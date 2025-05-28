package accounting

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// QueryOptions represents query parameters for entry searches
type QueryOptions struct {
	Filters   []Filter `json:"filters"`
	SortBy    string   `json:"sort_by,omitempty"`
	SortOrder string   `json:"sort_order,omitempty"`
	Limit     int      `json:"limit,omitempty"`
}

// Filter represents a query filter
type Filter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// generateUUID generates a new UUID string
func generateUUID() string {
	return uuid.New().String()
}

// ForensicService provides forensic accounting capabilities
type ForensicService struct {
	storage    *Storage
	eventStore *EventStore
}

// NewForensicService creates a new forensic service
func NewForensicService(storage *Storage, eventStore *EventStore) *ForensicService {
	return &ForensicService{
		storage:    storage,
		eventStore: eventStore,
	}
}

// MoneyTrail represents a complete path of money movement
type MoneyTrail struct {
	ID          string           `json:"id"`
	StartAmount Amount           `json:"start_amount"`
	EndAmount   Amount           `json:"end_amount"`
	StartDate   time.Time        `json:"start_date"`
	EndDate     time.Time        `json:"end_date"`
	Path        []MoneyTrailStep `json:"path"`
	TotalSteps  int              `json:"total_steps"`
	Suspicious  bool             `json:"suspicious"`
	Flags       []ForensicFlag   `json:"flags"`
}

// MoneyTrailStep represents one step in money movement
type MoneyTrailStep struct {
	TransactionID string    `json:"transaction_id"`
	FromAccount   string    `json:"from_account"`
	ToAccount     string    `json:"to_account"`
	Amount        Amount    `json:"amount"`
	Date          time.Time `json:"date"`
	Description   string    `json:"description"`
	UserID        string    `json:"user_id"`
}

// ForensicFlag represents suspicious activity indicators
type ForensicFlag struct {
	Type        FlagType  `json:"type"`
	Severity    Severity  `json:"severity"`
	Description string    `json:"description"`
	Evidence    []string  `json:"evidence"`
	Triggered   time.Time `json:"triggered"`
}

type FlagType string

const (
	FlagRoundAmounts        FlagType = "round_amounts"
	FlagHighFrequency       FlagType = "high_frequency"
	FlagUnusualTiming       FlagType = "unusual_timing"
	FlagComplexRouting      FlagType = "complex_routing"
	FlagLayering            FlagType = "layering"
	FlagRapidMovement       FlagType = "rapid_movement"
	FlagStructuring         FlagType = "structuring"
	FlagCircularTransfers   FlagType = "circular_transfers"
	FlagDormantReactivation FlagType = "dormant_reactivation"
)

type Severity string

const (
	SeverityLow    Severity = "low"
	SeverityMedium Severity = "medium"
	SeverityHigh   Severity = "high"
)

// TransactionGraph represents relationships between accounts
type TransactionGraph struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

// GraphNode represents an account in the transaction graph
type GraphNode struct {
	AccountID    string                 `json:"account_id"`
	AccountName  string                 `json:"account_name"`
	TotalInflow  Amount                 `json:"total_inflow"`
	TotalOutflow Amount                 `json:"total_outflow"`
	Centrality   float64                `json:"centrality"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// GraphEdge represents a relationship between accounts
type GraphEdge struct {
	FromAccount      string    `json:"from_account"`
	ToAccount        string    `json:"to_account"`
	TotalAmount      Amount    `json:"total_amount"`
	TransactionCount int       `json:"transaction_count"`
	FirstTransaction time.Time `json:"first_transaction"`
	LastTransaction  time.Time `json:"last_transaction"`
	Weight           float64   `json:"weight"`
}

// SuspiciousPattern represents detected patterns
type SuspiciousPattern struct {
	ID           string      `json:"id"`
	Type         FlagType    `json:"type"`
	Severity     Severity    `json:"severity"`
	Description  string      `json:"description"`
	Accounts     []string    `json:"accounts"`
	Transactions []string    `json:"transactions"`
	Timeline     []time.Time `json:"timeline"`
	Evidence     []string    `json:"evidence"`
	Confidence   float64     `json:"confidence"`
	DetectedAt   time.Time   `json:"detected_at"`
}

// TrackMoneyTrail follows money movement from source to destination
func (fs *ForensicService) TrackMoneyTrail(sourceAccountID string, startDate, endDate time.Time, minAmount int64) (*MoneyTrail, error) {
	// Get all entries for the source account in the date range
	entries, err := fs.storage.GetEntriesByAccount(sourceAccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get entries: %w", err)
	}

	// Get transactions for date filtering
	txnMap := make(map[string]*Transaction)
	var filteredEntries []*Entry

	for _, entry := range entries {
		// Get transaction if not already cached
		if _, exists := txnMap[entry.TransactionID]; !exists {
			txn, err := fs.storage.GetTransaction(entry.TransactionID)
			if err != nil {
				continue
			}
			txnMap[entry.TransactionID] = txn
		}

		txn := txnMap[entry.TransactionID]
		// Filter by date range and amount using transaction's ValidTime
		if (txn.ValidTime.After(startDate) || txn.ValidTime.Equal(startDate)) &&
			(txn.ValidTime.Before(endDate) || txn.ValidTime.Equal(endDate)) &&
			entry.Amount.Value >= minAmount {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	// Build money trail
	trail := &MoneyTrail{
		ID:        generateUUID(),
		StartDate: startDate,
		EndDate:   endDate,
		Path:      []MoneyTrailStep{},
	}

	// Group filtered entries by transaction
	txnEntryMap := make(map[string][]*Entry)
	for _, entry := range filteredEntries {
		txnEntryMap[entry.TransactionID] = append(txnEntryMap[entry.TransactionID], entry)
	}

	// Build trail steps
	for txnID, txnEntries := range txnEntryMap {
		// Get transaction details
		txn := txnMap[txnID]
		if txn == nil {
			continue
		}

		// Filter by minimum amount
		if len(txnEntries) > 0 && txnEntries[0].Amount.Value < minAmount {
			continue
		}

		// Find counterparty entries
		for _, entry := range txnEntries {
			if entry.AccountID == sourceAccountID {
				// Find corresponding entries in the same transaction
				for _, otherEntry := range txn.Entries {
					if otherEntry.ID != entry.ID {
						step := MoneyTrailStep{
							TransactionID: txnID,
							Amount:        entry.Amount,
							Date:          txn.ValidTime,
							Description:   txn.Description,
						}

						if entry.Type == Debit {
							step.FromAccount = otherEntry.AccountID
							step.ToAccount = entry.AccountID
						} else {
							step.FromAccount = entry.AccountID
							step.ToAccount = otherEntry.AccountID
						}

						trail.Path = append(trail.Path, step)
					}
				}
			}
		}
	}

	// Sort by date
	sort.Slice(trail.Path, func(i, j int) bool {
		return trail.Path[i].Date.Before(trail.Path[j].Date)
	})

	trail.TotalSteps = len(trail.Path)

	// Calculate total amounts
	if len(trail.Path) > 0 {
		trail.StartAmount = trail.Path[0].Amount
		trail.EndAmount = trail.Path[len(trail.Path)-1].Amount
	}

	// Analyze for suspicious patterns
	trail.Flags = fs.analyzeTrailForFlags(trail)
	trail.Suspicious = len(trail.Flags) > 0

	return trail, nil
}

// BuildTransactionGraph creates a graph representation of account relationships
func (fs *ForensicService) BuildTransactionGraph(startDate, endDate time.Time, companyID string) (*TransactionGraph, error) {
	// Get all transactions in the period
	query := &QueryOptions{
		Filters: []Filter{
			{Field: "valid_time", Operator: ">=", Value: startDate},
			{Field: "valid_time", Operator: "<=", Value: endDate},
		},
	}

	if companyID != "" {
		query.Filters = append(query.Filters, Filter{Field: "company_id", Operator: "=", Value: companyID})
	}

	entries, err := fs.storage.QueryEntries(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query entries: %w", err)
	}

	// Build nodes and edges
	nodeMap := make(map[string]*GraphNode)
	edgeMap := make(map[string]*GraphEdge)

	// Group entries by transaction
	txnMap := make(map[string][]*Entry)
	for _, entry := range entries {
		txnMap[entry.TransactionID] = append(txnMap[entry.TransactionID], entry)
	}

	// Process each transaction
	for txnID, txnEntries := range txnMap {
		txn, err := fs.storage.GetTransaction(txnID)
		if err != nil {
			continue
		}

		// Process entries in pairs (debit/credit)
		for i, entry1 := range txnEntries {
			// Initialize node if not exists
			if _, exists := nodeMap[entry1.AccountID]; !exists {
				account, _ := fs.storage.GetAccount(entry1.AccountID)
				nodeMap[entry1.AccountID] = &GraphNode{
					AccountID:    entry1.AccountID,
					AccountName:  account.Name,
					TotalInflow:  Amount{Value: 0, Currency: entry1.Amount.Currency},
					TotalOutflow: Amount{Value: 0, Currency: entry1.Amount.Currency},
					Metadata:     make(map[string]interface{}),
				}
			}

			// Update node totals
			if entry1.Type == Debit {
				nodeMap[entry1.AccountID].TotalInflow.Value += entry1.Amount.Value
			} else {
				nodeMap[entry1.AccountID].TotalOutflow.Value += entry1.Amount.Value
			}

			// Process relationships with other entries
			for j, entry2 := range txnEntries {
				if i != j {
					var fromAccount, toAccount string
					if entry1.Type == Debit && entry2.Type == Credit {
						fromAccount = entry2.AccountID
						toAccount = entry1.AccountID
					} else if entry1.Type == Credit && entry2.Type == Debit {
						fromAccount = entry1.AccountID
						toAccount = entry2.AccountID
					} else {
						continue
					}

					edgeKey := fmt.Sprintf("%s->%s", fromAccount, toAccount)
					if _, exists := edgeMap[edgeKey]; !exists {
						edgeMap[edgeKey] = &GraphEdge{
							FromAccount:      fromAccount,
							ToAccount:        toAccount,
							TotalAmount:      Amount{Value: 0, Currency: entry1.Amount.Currency},
							TransactionCount: 0,
							FirstTransaction: txn.ValidTime,
							LastTransaction:  txn.ValidTime,
						}
					}

					edge := edgeMap[edgeKey]
					edge.TotalAmount.Value += entry1.Amount.Value
					edge.TransactionCount++
					if txn.ValidTime.Before(edge.FirstTransaction) {
						edge.FirstTransaction = txn.ValidTime
					}
					if txn.ValidTime.After(edge.LastTransaction) {
						edge.LastTransaction = txn.ValidTime
					}
				}
			}
		}
	}

	// Calculate centrality measures
	fs.calculateCentrality(nodeMap, edgeMap)

	// Convert maps to slices
	graph := &TransactionGraph{
		Nodes: make([]GraphNode, 0, len(nodeMap)),
		Edges: make([]GraphEdge, 0, len(edgeMap)),
	}

	for _, node := range nodeMap {
		graph.Nodes = append(graph.Nodes, *node)
	}

	for _, edge := range edgeMap {
		graph.Edges = append(graph.Edges, *edge)
	}

	return graph, nil
}

// DetectSuspiciousPatterns analyzes transactions for suspicious patterns
func (fs *ForensicService) DetectSuspiciousPatterns(startDate, endDate time.Time, companyID string) ([]SuspiciousPattern, error) {
	var patterns []SuspiciousPattern

	// Get transactions for analysis
	query := &QueryOptions{
		Filters: []Filter{
			{Field: "valid_time", Operator: ">=", Value: startDate},
			{Field: "valid_time", Operator: "<=", Value: endDate},
		},
	}

	if companyID != "" {
		query.Filters = append(query.Filters, Filter{Field: "company_id", Operator: "=", Value: companyID})
	}

	entries, err := fs.storage.QueryEntries(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query entries: %w", err)
	}

	// Detect various suspicious patterns
	patterns = append(patterns, fs.detectRoundAmountPattern(entries)...)
	patterns = append(patterns, fs.detectHighFrequencyPattern(entries)...)
	patterns = append(patterns, fs.detectStructuringPattern(entries)...)
	patterns = append(patterns, fs.detectCircularTransferPattern(entries)...)
	patterns = append(patterns, fs.detectUnusualTimingPattern(entries)...)

	return patterns, nil
}

// Helper methods for pattern detection

func (fs *ForensicService) analyzeTrailForFlags(trail *MoneyTrail) []ForensicFlag {
	var flags []ForensicFlag

	// Check for rapid movement (multiple transactions in short time)
	if len(trail.Path) > 5 {
		rapidMovements := 0
		for i := 1; i < len(trail.Path); i++ {
			if trail.Path[i].Date.Sub(trail.Path[i-1].Date) < time.Hour {
				rapidMovements++
			}
		}

		if rapidMovements > 3 {
			flags = append(flags, ForensicFlag{
				Type:        FlagRapidMovement,
				Severity:    SeverityMedium,
				Description: "Multiple transactions within short time periods",
				Evidence:    []string{fmt.Sprintf("%d rapid movements detected", rapidMovements)},
				Triggered:   time.Now(),
			})
		}
	}

	// Check for complex routing (many intermediate accounts)
	if len(trail.Path) > 10 {
		flags = append(flags, ForensicFlag{
			Type:        FlagComplexRouting,
			Severity:    SeverityHigh,
			Description: "Complex routing through multiple accounts",
			Evidence:    []string{fmt.Sprintf("%d transaction steps", len(trail.Path))},
			Triggered:   time.Now(),
		})
	}

	return flags
}

func (fs *ForensicService) detectRoundAmountPattern(entries []*Entry) []SuspiciousPattern {
	roundAmountCount := 0
	var roundTransactions []string

	for _, entry := range entries {
		// Check if amount is suspiciously round (ends in multiple zeros)
		if entry.Amount.Value%100000 == 0 && entry.Amount.Value > 0 { // Ends in 00000 (like $1000.00)
			roundAmountCount++
			roundTransactions = append(roundTransactions, entry.TransactionID)
		}
	}

	var patterns []SuspiciousPattern
	if roundAmountCount > 10 { // More than 10 round amounts
		patterns = append(patterns, SuspiciousPattern{
			ID:           generateUUID(),
			Type:         FlagRoundAmounts,
			Severity:     SeverityMedium,
			Description:  "High frequency of round number transactions",
			Transactions: roundTransactions,
			Evidence:     []string{fmt.Sprintf("%d round amount transactions", roundAmountCount)},
			Confidence:   float64(roundAmountCount) / float64(len(entries)),
			DetectedAt:   time.Now(),
		})
	}

	return patterns
}

func (fs *ForensicService) detectHighFrequencyPattern(entries []*Entry) []SuspiciousPattern {
	// Group by account and day
	dailyActivity := make(map[string]map[string]int) // account -> date -> count

	// Get transaction details for date information
	txnCache := make(map[string]*Transaction)

	for _, entry := range entries {
		// Get transaction if not cached
		if _, exists := txnCache[entry.TransactionID]; !exists {
			txn, err := fs.storage.GetTransaction(entry.TransactionID)
			if err != nil {
				continue
			}
			txnCache[entry.TransactionID] = txn
		}

		txn := txnCache[entry.TransactionID]
		dateKey := txn.TransactionTime.Format("2006-01-02")
		if dailyActivity[entry.AccountID] == nil {
			dailyActivity[entry.AccountID] = make(map[string]int)
		}
		dailyActivity[entry.AccountID][dateKey]++
	}

	var patterns []SuspiciousPattern
	for accountID, activity := range dailyActivity {
		for date, count := range activity {
			if count > 50 { // More than 50 transactions per day
				patterns = append(patterns, SuspiciousPattern{
					ID:          generateUUID(),
					Type:        FlagHighFrequency,
					Severity:    SeverityHigh,
					Description: "Unusually high transaction frequency",
					Accounts:    []string{accountID},
					Evidence:    []string{fmt.Sprintf("%d transactions on %s", count, date)},
					Confidence:  0.8,
					DetectedAt:  time.Now(),
				})
			}
		}
	}

	return patterns
}

func (fs *ForensicService) detectStructuringPattern(entries []*Entry) []SuspiciousPattern {
	// Look for amounts just under reporting thresholds
	const threshold = 1000000 // $10,000 threshold
	var suspiciousAmounts []string
	structuringCount := 0

	for _, entry := range entries {
		// Check if amount is just under threshold (within 5%)
		if entry.Amount.Value > threshold*95/100 && entry.Amount.Value < threshold {
			structuringCount++
			suspiciousAmounts = append(suspiciousAmounts, fmt.Sprintf("$%.2f", float64(entry.Amount.Value)/100))
		}
	}

	var patterns []SuspiciousPattern
	if structuringCount > 5 {
		patterns = append(patterns, SuspiciousPattern{
			ID:          generateUUID(),
			Type:        FlagStructuring,
			Severity:    SeverityHigh,
			Description: "Potential structuring - amounts just under reporting thresholds",
			Evidence:    []string{fmt.Sprintf("%d transactions near threshold: %s", structuringCount, strings.Join(suspiciousAmounts, ", "))},
			Confidence:  0.9,
			DetectedAt:  time.Now(),
		})
	}

	return patterns
}

func (fs *ForensicService) detectCircularTransferPattern(entries []*Entry) []SuspiciousPattern {
	// This is a simplified version - a full implementation would use graph algorithms
	// to detect actual cycles in the transaction flow
	accountPairs := make(map[string]int) // "account1->account2" -> count

	// Group entries by transaction
	txnMap := make(map[string][]*Entry)
	for _, entry := range entries {
		txnMap[entry.TransactionID] = append(txnMap[entry.TransactionID], entry)
	}

	// Look for back-and-forth transactions
	for _, txnEntries := range txnMap {
		if len(txnEntries) == 2 {
			var account1, account2 string
			for _, entry := range txnEntries {
				if entry.Type == Debit {
					account2 = entry.AccountID
				} else {
					account1 = entry.AccountID
				}
			}

			if account1 != "" && account2 != "" {
				pair := fmt.Sprintf("%s->%s", account1, account2)
				reversePair := fmt.Sprintf("%s->%s", account2, account1)
				accountPairs[pair]++

				// Check if reverse transaction exists
				if accountPairs[reversePair] > 0 {
					// Potential circular transfer detected
				}
			}
		}
	}

	// For now, return empty - full cycle detection would require more complex graph analysis
	return []SuspiciousPattern{}
}

func (fs *ForensicService) detectUnusualTimingPattern(entries []*Entry) []SuspiciousPattern {
	weekendCount := 0
	afterHoursCount := 0

	// Get transaction details for timing information
	txnCache := make(map[string]*Transaction)

	for _, entry := range entries {
		// Get transaction if not cached
		if _, exists := txnCache[entry.TransactionID]; !exists {
			txn, err := fs.storage.GetTransaction(entry.TransactionID)
			if err != nil {
				continue
			}
			txnCache[entry.TransactionID] = txn
		}

		txn := txnCache[entry.TransactionID]
		weekday := txn.TransactionTime.Weekday()
		hour := txn.TransactionTime.Hour()

		// Weekend transactions
		if weekday == time.Saturday || weekday == time.Sunday {
			weekendCount++
		}

		// After hours (before 8 AM or after 6 PM)
		if hour < 8 || hour > 18 {
			afterHoursCount++
		}
	}

	var patterns []SuspiciousPattern
	totalEntries := len(entries)

	if weekendCount > totalEntries/10 { // More than 10% on weekends
		patterns = append(patterns, SuspiciousPattern{
			ID:          generateUUID(),
			Type:        FlagUnusualTiming,
			Severity:    SeverityMedium,
			Description: "High percentage of weekend transactions",
			Evidence:    []string{fmt.Sprintf("%d weekend transactions (%.1f%%)", weekendCount, float64(weekendCount)*100/float64(totalEntries))},
			Confidence:  0.7,
			DetectedAt:  time.Now(),
		})
	}

	if afterHoursCount > totalEntries/5 { // More than 20% after hours
		patterns = append(patterns, SuspiciousPattern{
			ID:          generateUUID(),
			Type:        FlagUnusualTiming,
			Severity:    SeverityMedium,
			Description: "High percentage of after-hours transactions",
			Evidence:    []string{fmt.Sprintf("%d after-hours transactions (%.1f%%)", afterHoursCount, float64(afterHoursCount)*100/float64(totalEntries))},
			Confidence:  0.6,
			DetectedAt:  time.Now(),
		})
	}

	return patterns
}

func (fs *ForensicService) calculateCentrality(nodeMap map[string]*GraphNode, edgeMap map[string]*GraphEdge) {
	// Simple degree centrality calculation
	for _, node := range nodeMap {
		inDegree := 0
		outDegree := 0

		for _, edge := range edgeMap {
			if edge.ToAccount == node.AccountID {
				inDegree++
			}
			if edge.FromAccount == node.AccountID {
				outDegree++
			}
		}

		// Simple centrality measure
		node.Centrality = float64(inDegree + outDegree)
	}
}
