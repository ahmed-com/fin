package accounting

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PostingEngine handles transaction posting with validation and balance checking
type PostingEngine struct {
	storage    *Storage
	eventStore *EventStore
	processor  *EventProcessor
}

// NewPostingEngine creates a new posting engine
func NewPostingEngine(storage *Storage, eventStore *EventStore, processor *EventProcessor) *PostingEngine {
	return &PostingEngine{
		storage:    storage,
		eventStore: eventStore,
		processor:  processor,
	}
}

// PostingError represents an error that occurred during posting
type PostingError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (pe PostingError) Error() string {
	return fmt.Sprintf("%s: %s", pe.Code, pe.Message)
}

// ValidationResult contains the result of transaction validation
type ValidationResult struct {
	Valid  bool           `json:"valid"`
	Errors []PostingError `json:"errors,omitempty"`
}

// ValidateTransaction validates a transaction before posting
func (pe *PostingEngine) ValidateTransaction(txn *Transaction) *ValidationResult {
	result := &ValidationResult{Valid: true}

	// Check if transaction balances (debits = credits)
	if err := pe.validateBalance(txn); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, PostingError{
			Code:    "UNBALANCED_TRANSACTION",
			Message: err.Error(),
		})
	}

	// Validate all accounts exist
	if err := pe.validateAccounts(txn); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, PostingError{
			Code:    "INVALID_ACCOUNT",
			Message: err.Error(),
		})
	}

	// Check period is open
	if err := pe.validatePeriod(txn.ValidTime); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, PostingError{
			Code:    "PERIOD_CLOSED",
			Message: err.Error(),
		})
	}

	return result
}

// validateBalance ensures debits equal credits
func (pe *PostingEngine) validateBalance(txn *Transaction) error {
	debitTotal := int64(0)
	creditTotal := int64(0)

	for _, entry := range txn.Entries {
		if entry.Type == Debit {
			debitTotal += entry.Amount.Value
		} else {
			creditTotal += entry.Amount.Value
		}
	}

	if debitTotal != creditTotal {
		return fmt.Errorf("transaction does not balance: debits=%d, credits=%d", debitTotal, creditTotal)
	}

	return nil
}

// validateAccounts ensures all referenced accounts exist
func (pe *PostingEngine) validateAccounts(txn *Transaction) error {
	for _, entry := range txn.Entries {
		_, err := pe.storage.GetAccount(entry.AccountID)
		if err != nil {
			return fmt.Errorf("account %s does not exist", entry.AccountID)
		}
	}
	return nil
}

// validatePeriod checks if the transaction date is in an open period
func (pe *PostingEngine) validatePeriod(validTime time.Time) error {
	// For now, assume all periods are open
	// In a real implementation, you'd check against period closing dates
	return nil
}

// PostTransaction posts a transaction to the ledger
func (pe *PostingEngine) PostTransaction(txn *Transaction, userID string) error {
	// Validate transaction
	validation := pe.ValidateTransaction(txn)
	if !validation.Valid {
		return fmt.Errorf("transaction validation failed: %v", validation.Errors)
	}

	// Set transaction status to posted
	txn.Status = Posted
	txn.UpdatedAt = time.Now()

	// Generate entries with IDs
	for i := range txn.Entries {
		if txn.Entries[i].ID == "" {
			txn.Entries[i].ID = uuid.New().String()
		}
		txn.Entries[i].TransactionID = txn.ID
	}

	// Create posting event
	event, err := pe.eventStore.CreateEvent(
		EventPostTransaction,
		TransactionPostedEvent{
			TransactionID: txn.ID,
			PostedAt:      time.Now(),
			Entries:       txn.Entries,
		},
		txn.ValidTime,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to create posting event: %w", err)
	}

	// Process the event
	if err := pe.processor.ProcessEvent(event); err != nil {
		return fmt.Errorf("failed to process posting event: %w", err)
	}

	return nil
}

// ReverseTransaction creates a reversing transaction
func (pe *PostingEngine) ReverseTransaction(originalTxnID string, description string, userID string) (*Transaction, error) {
	// Get original transaction
	originalTxn, err := pe.storage.GetTransaction(originalTxnID)
	if err != nil {
		return nil, fmt.Errorf("failed to get original transaction: %w", err)
	}

	if originalTxn.Status != Posted {
		return nil, fmt.Errorf("can only reverse posted transactions")
	}

	// Create reversing transaction
	reversingTxn := &Transaction{
		ID:              uuid.New().String(),
		Description:     description,
		ValidTime:       time.Now(),
		TransactionTime: time.Now(),
		Status:          Pending,
		SourceRef:       fmt.Sprintf("REVERSAL_%s", originalTxnID),
		UserID:          userID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Reverse all entries (flip debit/credit)
	for _, entry := range originalTxn.Entries {
		reversedType := Credit
		if entry.Type == Credit {
			reversedType = Debit
		}

		reversingEntry := Entry{
			ID:            uuid.New().String(),
			TransactionID: reversingTxn.ID,
			AccountID:     entry.AccountID,
			Type:          reversedType,
			Amount:        entry.Amount,
			Dimensions:    entry.Dimensions,
		}

		reversingTxn.Entries = append(reversingTxn.Entries, reversingEntry)
	}

	// Create transaction creation event
	_, err = pe.eventStore.CreateEvent(
		EventCreateTransaction,
		TransactionCreatedEvent{Transaction: reversingTxn},
		reversingTxn.ValidTime,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction event: %w", err)
	}

	// Save the transaction first
	if err := pe.storage.SaveTransaction(reversingTxn); err != nil {
		return nil, fmt.Errorf("failed to save reversing transaction: %w", err)
	}

	// Post the reversing transaction
	if err := pe.PostTransaction(reversingTxn, userID); err != nil {
		return nil, fmt.Errorf("failed to post reversing transaction: %w", err)
	}

	// Mark original transaction as reversed
	originalTxn.Status = Reversed
	originalTxn.UpdatedAt = time.Now()
	if err := pe.storage.SaveTransaction(originalTxn); err != nil {
		return nil, fmt.Errorf("failed to update original transaction status: %w", err)
	}

	return reversingTxn, nil
}

// CalculateAccountBalance calculates the current balance of an account
func (pe *PostingEngine) CalculateAccountBalance(accountID string, asOfDate time.Time) (*Amount, error) {
	account, err := pe.storage.GetAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	entries, err := pe.storage.GetEntriesByAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get entries: %w", err)
	}

	balance := &Amount{
		Value:    0,
		Currency: account.Currency,
	}

	for _, entry := range entries {
		// Get transaction to check valid time
		txn, err := pe.storage.GetTransaction(entry.TransactionID)
		if err != nil {
			continue // Skip if transaction not found
		}

		// Only include transactions valid up to the as-of date
		if txn.ValidTime.After(asOfDate) || txn.Status != Posted {
			continue
		}

		// Apply entry based on account type and entry type
		multiplier := pe.getBalanceMultiplier(account.Type, entry.Type)
		balance.Value += entry.Amount.Value * int64(multiplier)
	}

	return balance, nil
}

// getBalanceMultiplier returns the multiplier for calculating balance based on account and entry type
func (pe *PostingEngine) getBalanceMultiplier(accountType AccountType, entryType EntryType) int {
	// Normal balance sides:
	// Assets: Debit increases (+1), Credit decreases (-1)
	// Liabilities: Credit increases (+1), Debit decreases (-1)
	// Equity: Credit increases (+1), Debit decreases (-1)
	// Income: Credit increases (+1), Debit decreases (-1)
	// Expense: Debit increases (+1), Credit decreases (-1)

	switch accountType {
	case Asset, Expense:
		if entryType == Debit {
			return 1
		}
		return -1
	case Liability, Equity, Income:
		if entryType == Credit {
			return 1
		}
		return -1
	default:
		return 0
	}
}
