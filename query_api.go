package accounting

import (
	"fmt"
	"time"
)

// QueryFilter represents filters for querying accounting data
type QueryFilter struct {
	ValidTimeFrom       *time.Time          `json:"valid_time_from,omitempty"`
	ValidTimeTo         *time.Time          `json:"valid_time_to,omitempty"`
	TransactionTimeFrom *time.Time          `json:"transaction_time_from,omitempty"`
	TransactionTimeTo   *time.Time          `json:"transaction_time_to,omitempty"`
	AccountIDs          []string            `json:"account_ids,omitempty"`
	AccountTypes        []AccountType       `json:"account_types,omitempty"`
	Dimensions          []Dimension         `json:"dimensions,omitempty"`
	Status              []TransactionStatus `json:"status,omitempty"`
	Currencies          []Currency          `json:"currencies,omitempty"`
}

// QueryResult contains paginated query results
type QueryResult struct {
	Transactions []*Transaction `json:"transactions"`
	Entries      []*Entry       `json:"entries"`
	TotalCount   int            `json:"total_count"`
	Page         int            `json:"page"`
	PageSize     int            `json:"page_size"`
}

// BalanceResult represents account balance information
type BalanceResult struct {
	AccountID   string      `json:"account_id"`
	AccountName string      `json:"account_name"`
	AccountType AccountType `json:"account_type"`
	Balance     *Amount     `json:"balance"`
	AsOfDate    time.Time   `json:"as_of_date"`
}

// DimensionRollup represents aggregated data by dimensions
type DimensionRollup struct {
	Dimensions []Dimension `json:"dimensions"`
	Amount     *Amount     `json:"amount"`
	Count      int         `json:"count"`
}

// QueryAPI provides query functionality for the accounting system
type QueryAPI struct {
	storage       *Storage
	postingEngine *PostingEngine
}

// NewQueryAPI creates a new query API
func NewQueryAPI(storage *Storage, postingEngine *PostingEngine) *QueryAPI {
	return &QueryAPI{
		storage:       storage,
		postingEngine: postingEngine,
	}
}

// GetTransactions retrieves transactions based on filters
func (qa *QueryAPI) GetTransactions(filter *QueryFilter, page, pageSize int) (*QueryResult, error) {
	// This is a simplified implementation
	// In a real system, you'd implement proper indexing and pagination

	// For now, we'll get all transactions and filter in memory
	// This is not efficient for large datasets

	result := &QueryResult{
		Transactions: []*Transaction{},
		Page:         page,
		PageSize:     pageSize,
	}

	// In a real implementation, you'd use indexed queries
	// For this MVP, we'll use a simple approach

	return result, nil
}

// GetAccountBalance gets the balance of an account as of a specific date
func (qa *QueryAPI) GetAccountBalance(accountID string, asOfDate time.Time) (*BalanceResult, error) {
	account, err := qa.storage.GetAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	balance, err := qa.postingEngine.CalculateAccountBalance(accountID, asOfDate)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate balance: %w", err)
	}

	return &BalanceResult{
		AccountID:   account.ID,
		AccountName: account.Name,
		AccountType: account.Type,
		Balance:     balance,
		AsOfDate:    asOfDate,
	}, nil
}

// GetTrialBalance generates a trial balance report
func (qa *QueryAPI) GetTrialBalance(asOfDate time.Time, accountTypes []AccountType) ([]*BalanceResult, error) {
	var results []*BalanceResult

	// Get all accounts - for this implementation, we'll use a simple approach
	// In a real system, you'd have proper account indexing

	// Get accounts from storage by iterating through known account IDs
	// This is a simplified approach for the demo
	standardAccountIDs := []string{"cash", "accounts_receivable", "accounts_payable", "revenue", "expenses", "unearned_revenue"}

	for _, accountID := range standardAccountIDs {
		account, err := qa.storage.GetAccount(accountID)
		if err != nil {
			continue // Skip if account doesn't exist
		}

		// Filter by account type if specified
		if len(accountTypes) > 0 {
			found := false
			for _, accountType := range accountTypes {
				if account.Type == accountType {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		balance, err := qa.postingEngine.CalculateAccountBalance(account.ID, asOfDate)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate balance for account %s: %w", account.ID, err)
		}

		// Include all accounts in trial balance, even those with zero balance
		results = append(results, &BalanceResult{
			AccountID:   account.ID,
			AccountName: account.Name,
			AccountType: account.Type,
			Balance:     balance,
			AsOfDate:    asOfDate,
		})
	}

	return results, nil
}

// GetDimensionRollup aggregates amounts by dimensions
func (qa *QueryAPI) GetDimensionRollup(filter *QueryFilter, rollupDimensions []DimensionKey) ([]*DimensionRollup, error) {
	rollupMap := make(map[string]*DimensionRollup)

	// Get all entries matching the filter
	// This is simplified - in a real system you'd use proper indexing
	entries := []*Entry{} // Placeholder

	for _, entry := range entries {
		// Create a key from the rollup dimensions
		key := qa.createDimensionKey(entry.Dimensions, rollupDimensions)

		if rollup, exists := rollupMap[key]; exists {
			// Add to existing rollup
			rollup.Amount.Value += entry.Amount.Value
			rollup.Count++
		} else {
			// Create new rollup
			dimensions := qa.extractDimensions(entry.Dimensions, rollupDimensions)
			rollupMap[key] = &DimensionRollup{
				Dimensions: dimensions,
				Amount: &Amount{
					Value:    entry.Amount.Value,
					Currency: entry.Amount.Currency,
				},
				Count: 1,
			}
		}
	}

	// Convert map to slice
	var results []*DimensionRollup
	for _, rollup := range rollupMap {
		results = append(results, rollup)
	}

	return results, nil
}

// createDimensionKey creates a unique key from dimension values
func (qa *QueryAPI) createDimensionKey(dimensions []Dimension, rollupKeys []DimensionKey) string {
	key := ""
	for _, rollupKey := range rollupKeys {
		value := qa.getDimensionValue(dimensions, rollupKey)
		key += fmt.Sprintf("%s:%s|", rollupKey, value)
	}
	return key
}

// getDimensionValue gets the value of a specific dimension
func (qa *QueryAPI) getDimensionValue(dimensions []Dimension, key DimensionKey) string {
	for _, dim := range dimensions {
		if dim.Key == key {
			return dim.Value
		}
	}
	return "N/A"
}

// extractDimensions extracts specific dimensions for rollup
func (qa *QueryAPI) extractDimensions(dimensions []Dimension, rollupKeys []DimensionKey) []Dimension {
	var result []Dimension
	for _, rollupKey := range rollupKeys {
		value := qa.getDimensionValue(dimensions, rollupKey)
		result = append(result, Dimension{
			Key:   rollupKey,
			Value: value,
		})
	}
	return result
}

// GetAccountHierarchy builds the account hierarchy tree
func (qa *QueryAPI) GetAccountHierarchy() ([]*AccountNode, error) {
	// Get all accounts
	// In a real system, this would be properly indexed
	accounts := []*Account{} // Placeholder

	// Build hierarchy
	accountMap := make(map[string]*AccountNode)
	var roots []*AccountNode

	// Create all nodes
	for _, account := range accounts {
		node := &AccountNode{
			Account:  account,
			Children: []*AccountNode{},
		}
		accountMap[account.ID] = node
	}

	// Build parent-child relationships
	for _, account := range accounts {
		node := accountMap[account.ID]
		if account.ParentID == "" {
			roots = append(roots, node)
		} else {
			if parent, exists := accountMap[account.ParentID]; exists {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return roots, nil
}

// AccountNode represents a node in the account hierarchy
type AccountNode struct {
	Account  *Account       `json:"account"`
	Children []*AccountNode `json:"children"`
}
