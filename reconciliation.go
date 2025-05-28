package accounting

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ReconciliationService handles bank statement and account reconciliation
type ReconciliationService struct {
	storage  *Storage
	queryAPI *QueryAPI
}

// NewReconciliationService creates a new reconciliation service
func NewReconciliationService(storage *Storage, queryAPI *QueryAPI) *ReconciliationService {
	return &ReconciliationService{
		storage:  storage,
		queryAPI: queryAPI,
	}
}

// ExternalStatement represents an external bank or statement line
type ExternalStatement struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      *Amount   `json:"amount"`
	Reference   string    `json:"reference"`
	BankAccount string    `json:"bank_account"`
}

// ReconciliationMatch represents a potential match between external statement and internal entries
type ReconciliationMatch struct {
	ExternalStatement *ExternalStatement `json:"external_statement"`
	InternalEntries   []*Entry           `json:"internal_entries"`
	MatchScore        float64            `json:"match_score"`
	MatchType         string             `json:"match_type"` // "EXACT", "PARTIAL", "SUGGESTED"
}

// ReconciliationSummary provides summary information about reconciliation status
type ReconciliationSummary struct {
	AccountID          string  `json:"account_id"`
	StatementBalance   *Amount `json:"statement_balance"`
	BookBalance        *Amount `json:"book_balance"`
	Difference         *Amount `json:"difference"`
	ReconciledCount    int     `json:"reconciled_count"`
	UnreconciledCount  int     `json:"unreconciled_count"`
	ReconciliationRate float64 `json:"reconciliation_rate"`
}

// AutoReconcile attempts to automatically match external statements with internal entries
func (rs *ReconciliationService) AutoReconcile(accountID string, statements []*ExternalStatement) ([]*ReconciliationMatch, error) {
	var matches []*ReconciliationMatch

	// Get unreconciled entries for the account
	entries, err := rs.getUnreconciledEntries(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unreconciled entries: %w", err)
	}

	for _, statement := range statements {
		match := rs.findBestMatch(statement, entries)
		if match != nil {
			matches = append(matches, match)
		}
	}

	return matches, nil
}

// findBestMatch finds the best matching internal entries for an external statement
func (rs *ReconciliationService) findBestMatch(statement *ExternalStatement, entries []*Entry) *ReconciliationMatch {
	var bestMatch *ReconciliationMatch
	bestScore := 0.0

	// Try exact amount matches first
	for _, entry := range entries {
		if rs.amountsMatch(statement.Amount, &entry.Amount) {
			// Check date proximity (within 3 days)
			txn, err := rs.storage.GetTransaction(entry.TransactionID)
			if err != nil {
				continue
			}

			daysDiff := rs.daysBetween(statement.Date, txn.ValidTime)
			if daysDiff <= 3 {
				score := 1.0 - (float64(daysDiff) * 0.1) // Reduce score by 0.1 per day difference

				if score > bestScore {
					bestScore = score
					bestMatch = &ReconciliationMatch{
						ExternalStatement: statement,
						InternalEntries:   []*Entry{entry},
						MatchScore:        score,
						MatchType:         "EXACT",
					}
				}
			}
		}
	}

	// If no exact match, try to find combination matches
	if bestMatch == nil {
		combinations := rs.findCombinationMatches(statement, entries)
		for _, combo := range combinations {
			if combo.MatchScore > bestScore {
				bestScore = combo.MatchScore
				bestMatch = combo
			}
		}
	}

	return bestMatch
}

// amountsMatch checks if two amounts are equal (considering currency)
func (rs *ReconciliationService) amountsMatch(amt1, amt2 *Amount) bool {
	return amt1.Value == amt2.Value && amt1.Currency == amt2.Currency
}

// daysBetween calculates the number of days between two dates
func (rs *ReconciliationService) daysBetween(date1, date2 time.Time) int {
	diff := date1.Sub(date2)
	if diff < 0 {
		diff = -diff
	}
	return int(diff.Hours() / 24)
}

// findCombinationMatches tries to find combinations of entries that sum to the statement amount
func (rs *ReconciliationService) findCombinationMatches(statement *ExternalStatement, entries []*Entry) []*ReconciliationMatch {
	var matches []*ReconciliationMatch

	// This is a simplified implementation
	// In practice, you'd use more sophisticated algorithms for subset sum problems

	// Try combinations of 2-3 entries
	for i, entry1 := range entries {
		for j := i + 1; j < len(entries); j++ {
			entry2 := entries[j]

			combinedAmount := &Amount{
				Value:    entry1.Amount.Value + entry2.Amount.Value,
				Currency: entry1.Amount.Currency,
			}

			if rs.amountsMatch(statement.Amount, combinedAmount) {
				matches = append(matches, &ReconciliationMatch{
					ExternalStatement: statement,
					InternalEntries:   []*Entry{entry1, entry2},
					MatchScore:        0.8, // Lower score for combination matches
					MatchType:         "PARTIAL",
				})
			}
		}
	}

	return matches
}

// getUnreconciledEntries gets all unreconciled entries for an account
func (rs *ReconciliationService) getUnreconciledEntries(accountID string) ([]*Entry, error) {
	allEntries, err := rs.storage.GetEntriesByAccount(accountID)
	if err != nil {
		return nil, err
	}

	var unreconciled []*Entry
	for _, entry := range allEntries {
		// Check if entry is already reconciled
		if !rs.isEntryReconciled(entry.ID) {
			unreconciled = append(unreconciled, entry)
		}
	}

	return unreconciled, nil
}

// isEntryReconciled checks if an entry is already reconciled
func (rs *ReconciliationService) isEntryReconciled(entryID string) bool {
	// In a real implementation, you'd query the reconciliation records
	// For now, return false (assume all entries are unreconciled)
	return false
}

// ConfirmReconciliation confirms a reconciliation match and creates a reconciliation record
func (rs *ReconciliationService) ConfirmReconciliation(match *ReconciliationMatch, userID string) (*Reconciliation, error) {
	entryIDs := make([]string, len(match.InternalEntries))
	for i, entry := range match.InternalEntries {
		entryIDs[i] = entry.ID
	}

	reconciliation := &Reconciliation{
		ID:          uuid.New().String(),
		ExternalRef: match.ExternalStatement.Reference,
		EntryIDs:    entryIDs,
		Status:      Reconciled,
		CreatedAt:   time.Now(),
		CompletedAt: &[]time.Time{time.Now()}[0],
	}

	if err := rs.storage.SaveReconciliation(reconciliation); err != nil {
		return nil, fmt.Errorf("failed to save reconciliation: %w", err)
	}

	return reconciliation, nil
}

// GetReconciliationSummary provides a summary of reconciliation status for an account
func (rs *ReconciliationService) GetReconciliationSummary(accountID string, asOfDate time.Time) (*ReconciliationSummary, error) {
	// Get book balance
	bookBalance, err := rs.queryAPI.GetAccountBalance(accountID, asOfDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get book balance: %w", err)
	}

	// Count reconciled and unreconciled entries
	entries, err := rs.storage.GetEntriesByAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get entries: %w", err)
	}

	reconciledCount := 0
	unreconciledCount := 0

	for _, entry := range entries {
		if rs.isEntryReconciled(entry.ID) {
			reconciledCount++
		} else {
			unreconciledCount++
		}
	}

	totalCount := reconciledCount + unreconciledCount
	reconciliationRate := 0.0
	if totalCount > 0 {
		reconciliationRate = float64(reconciledCount) / float64(totalCount)
	}

	summary := &ReconciliationSummary{
		AccountID:          accountID,
		BookBalance:        bookBalance.Balance,
		ReconciledCount:    reconciledCount,
		UnreconciledCount:  unreconciledCount,
		ReconciliationRate: reconciliationRate,
	}

	// For statement balance, you'd typically get this from bank feeds or manual input
	// For now, we'll set it to book balance (assuming perfect reconciliation scenario)
	summary.StatementBalance = &Amount{
		Value:    bookBalance.Balance.Value,
		Currency: bookBalance.Balance.Currency,
	}

	summary.Difference = &Amount{
		Value:    summary.StatementBalance.Value - summary.BookBalance.Value,
		Currency: summary.BookBalance.Currency,
	}

	return summary, nil
}

// CreateManualReconciliation creates a manual reconciliation entry
func (rs *ReconciliationService) CreateManualReconciliation(externalRef string, entryIDs []string, userID string) (*Reconciliation, error) {
	reconciliation := &Reconciliation{
		ID:          uuid.New().String(),
		ExternalRef: externalRef,
		EntryIDs:    entryIDs,
		Status:      Reconciled,
		CreatedAt:   time.Now(),
		CompletedAt: &[]time.Time{time.Now()}[0],
	}

	if err := rs.storage.SaveReconciliation(reconciliation); err != nil {
		return nil, fmt.Errorf("failed to save manual reconciliation: %w", err)
	}

	return reconciliation, nil
}
