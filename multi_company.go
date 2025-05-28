package accounting

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Company represents a business entity in a multi-company environment
type Company struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	LegalName       string                 `json:"legal_name"`
	TaxID           string                 `json:"tax_id"`
	BaseCurrency    string                 `json:"base_currency"`
	FiscalYearEnd   time.Time              `json:"fiscal_year_end"`
	Address         *Address               `json:"address,omitempty"`
	Settings        *CompanySettings       `json:"settings"`
	CreatedAt       time.Time              `json:"created_at"`
	CreatedBy       string                 `json:"created_by"`
	Status          CompanyStatus          `json:"status"`
	ParentCompanyID string                 `json:"parent_company_id,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// Address represents a business address
type Address struct {
	Street1    string `json:"street1"`
	Street2    string `json:"street2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// CompanySettings represents company-specific settings
type CompanySettings struct {
	DefaultChartOfAccounts string             `json:"default_chart_of_accounts"`
	AllowIntercompanyTxn   bool               `json:"allow_intercompany_transactions"`
	RequireApprovalOver    *Amount            `json:"require_approval_over,omitempty"`
	AutoPostingRules       []*AutoPostingRule `json:"auto_posting_rules,omitempty"`
	PeriodLockingPolicy    string             `json:"period_locking_policy"`
	ReportingCurrency      string             `json:"reporting_currency"`
}

// AutoPostingRule represents automatic posting rules
type AutoPostingRule struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Condition string           `json:"condition"` // JSON-based condition
	Actions   []*PostingAction `json:"actions"`
	IsActive  bool             `json:"is_active"`
	CreatedAt time.Time        `json:"created_at"`
}

// PostingAction represents an action to take during auto-posting
type PostingAction struct {
	Type       string            `json:"type"` // "CREATE_ENTRY", "SET_DIMENSION", etc.
	Parameters map[string]string `json:"parameters"`
}

// CompanyStatus represents the status of a company
type CompanyStatus string

const (
	CompanyActive    CompanyStatus = "ACTIVE"
	CompanyInactive  CompanyStatus = "INACTIVE"
	CompanySuspended CompanyStatus = "SUSPENDED"
	CompanyMerged    CompanyStatus = "MERGED"
)

// IntercompanyTransaction represents a transaction between companies
type IntercompanyTransaction struct {
	ID                  string             `json:"id"`
	Description         string             `json:"description"`
	SourceCompanyID     string             `json:"source_company_id"`
	TargetCompanyID     string             `json:"target_company_id"`
	SourceTransactionID string             `json:"source_transaction_id"`
	TargetTransactionID string             `json:"target_transaction_id"`
	Amount              *Amount            `json:"amount"`
	ExchangeRate        float64            `json:"exchange_rate,omitempty"`
	MatchingStatus      IntercompanyStatus `json:"matching_status"`
	CreatedAt           time.Time          `json:"created_at"`
	CreatedBy           string             `json:"created_by"`
	ReconciledAt        *time.Time         `json:"reconciled_at,omitempty"`
	ReconciledBy        string             `json:"reconciled_by,omitempty"`
}

// IntercompanyStatus represents the status of intercompany transactions
type IntercompanyStatus string

const (
	IntercompanyPending    IntercompanyStatus = "PENDING"
	IntercompanyMatched    IntercompanyStatus = "MATCHED"
	IntercompanyReconciled IntercompanyStatus = "RECONCILED"
	IntercompanyDispute    IntercompanyStatus = "DISPUTE"
)

// ConsolidationGroup represents a group of companies for consolidation
type ConsolidationGroup struct {
	ID                  string             `json:"id"`
	Name                string             `json:"name"`
	ParentCompany       string             `json:"parent_company"`
	ChildCompanies      []string           `json:"child_companies"`
	ConsolidationMethod string             `json:"consolidation_method"` // "FULL", "EQUITY", "PROPORTIONAL"
	EliminationRules    []*EliminationRule `json:"elimination_rules"`
	CreatedAt           time.Time          `json:"created_at"`
	CreatedBy           string             `json:"created_by"`
}

// EliminationRule represents consolidation elimination rules
type EliminationRule struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	RuleType   string            `json:"rule_type"` // "INTERCOMPANY_SALES", "INVESTMENT", "DIVIDEND"
	Accounts   []string          `json:"accounts"`
	IsActive   bool              `json:"is_active"`
	Parameters map[string]string `json:"parameters"`
}

// MultiCompanyEngine manages multi-company operations
type MultiCompanyEngine struct {
	storage          Storage
	accountingEngine *AccountingEngine
	companies        map[string]*Company
	engines          map[string]*AccountingEngine // Cache for company accounting engines
}

// NewMultiCompanyEngine creates a new multi-company engine
func NewMultiCompanyEngine(storage Storage) *MultiCompanyEngine {
	return &MultiCompanyEngine{
		storage:   storage,
		companies: make(map[string]*Company),
		engines:   make(map[string]*AccountingEngine),
	}
}

// Close closes all cached accounting engines
func (mce *MultiCompanyEngine) Close() error {
	for _, engine := range mce.engines {
		if err := engine.Close(); err != nil {
			// Log error but continue closing others
			continue
		}
	}
	mce.engines = make(map[string]*AccountingEngine)
	return nil
}

// CreateCompany creates a new company
func (mce *MultiCompanyEngine) CreateCompany(company *Company, userID string) error {
	company.CreatedAt = time.Now()
	company.CreatedBy = userID
	company.Status = CompanyActive

	if err := mce.storage.SaveCompany(company); err != nil {
		return fmt.Errorf("failed to save company: %w", err)
	}

	// Create an accounting engine for this company
	engineKey := fmt.Sprintf("company_%s_%d", company.ID, time.Now().UnixNano())
	engine, err := NewAccountingEngine(engineKey + ".db")
	if err != nil {
		return fmt.Errorf("failed to create accounting engine for company: %w", err)
	}

	// Store in cache
	mce.companies[company.ID] = company
	mce.engines[company.ID] = engine // Cache the engine

	// Create standard chart of accounts if specified
	if company.Settings.DefaultChartOfAccounts != "" {
		if err := engine.CreateStandardAccounts(userID); err != nil {
			return fmt.Errorf("failed to create standard accounts: %w", err)
		}
	}

	return nil
}

// GetCompany retrieves a company by ID
func (mce *MultiCompanyEngine) GetCompany(companyID string) (*Company, error) {
	if company, exists := mce.companies[companyID]; exists {
		return company, nil
	}

	company, err := mce.storage.GetCompany(companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	mce.companies[companyID] = company
	return company, nil
}

// GetAccountingEngine gets the accounting engine for a specific company
func (mce *MultiCompanyEngine) GetAccountingEngine(companyID string) (*AccountingEngine, error) {
	// Verify company exists
	_, err := mce.GetCompany(companyID)
	if err != nil {
		return nil, err
	}

	// Check if engine is already cached
	if engine, exists := mce.engines[companyID]; exists {
		return engine, nil
	}

	// Create new accounting engine for this company
	engineKey := fmt.Sprintf("company_%s_%d", companyID, time.Now().UnixNano())
	engine, err := NewAccountingEngine(engineKey + ".db")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Cache the engine
	mce.engines[companyID] = engine
	return engine, nil
}

// CreateIntercompanyTransaction creates a transaction between two companies
func (mce *MultiCompanyEngine) CreateIntercompanyTransaction(
	sourceCompanyID, targetCompanyID string,
	amount *Amount,
	description string,
	userID string) (*IntercompanyTransaction, error) {

	// Verify both companies exist and allow intercompany transactions
	sourceCompany, err := mce.GetCompany(sourceCompanyID)
	if err != nil {
		return nil, fmt.Errorf("source company not found: %w", err)
	}

	targetCompany, err := mce.GetCompany(targetCompanyID)
	if err != nil {
		return nil, fmt.Errorf("target company not found: %w", err)
	}

	if !sourceCompany.Settings.AllowIntercompanyTxn || !targetCompany.Settings.AllowIntercompanyTxn {
		return nil, fmt.Errorf("intercompany transactions not allowed")
	}

	// Create intercompany transaction record
	intercompanyTxn := &IntercompanyTransaction{
		ID:              uuid.New().String(),
		Description:     description,
		SourceCompanyID: sourceCompanyID,
		TargetCompanyID: targetCompanyID,
		Amount:          amount,
		MatchingStatus:  IntercompanyPending,
		CreatedAt:       time.Now(),
		CreatedBy:       userID,
	}

	// Get accounting engines for both companies
	sourceEngine, err := mce.GetAccountingEngine(sourceCompanyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source accounting engine: %w", err)
	}

	targetEngine, err := mce.GetAccountingEngine(targetCompanyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get target accounting engine: %w", err)
	}

	// Create source transaction (outgoing)
	sourceTxn := &Transaction{
		Description: fmt.Sprintf("Intercompany transfer to %s: %s", targetCompany.Name, description),
		ValidTime:   time.Now(),
		Entries: []Entry{
			{
				AccountID: "intercompany_receivable",
				Type:      Debit,
				Amount:    *amount,
				Dimensions: []Dimension{
					{Key: "intercompany", Value: targetCompanyID},
					{Key: "transaction_type", Value: "intercompany_transfer"},
				},
			},
			{
				AccountID: "cash", // or appropriate source account
				Type:      Credit,
				Amount:    *amount,
			},
		},
	}

	if err := sourceEngine.CreateTransaction(sourceTxn, userID); err != nil {
		return nil, fmt.Errorf("failed to create source transaction: %w", err)
	}

	if err := sourceEngine.PostTransaction(sourceTxn.ID, userID); err != nil {
		return nil, fmt.Errorf("failed to post source transaction: %w", err)
	}

	// Create target transaction (incoming)
	targetTxn := &Transaction{
		Description: fmt.Sprintf("Intercompany transfer from %s: %s", sourceCompany.Name, description),
		ValidTime:   time.Now(),
		Entries: []Entry{
			{
				AccountID: "cash", // or appropriate target account
				Type:      Debit,
				Amount:    *amount,
			},
			{
				AccountID: "intercompany_payable",
				Type:      Credit,
				Amount:    *amount,
				Dimensions: []Dimension{
					{Key: "intercompany", Value: sourceCompanyID},
					{Key: "transaction_type", Value: "intercompany_transfer"},
				},
			},
		},
	}

	if err := targetEngine.CreateTransaction(targetTxn, userID); err != nil {
		return nil, fmt.Errorf("failed to create target transaction: %w", err)
	}

	if err := targetEngine.PostTransaction(targetTxn.ID, userID); err != nil {
		return nil, fmt.Errorf("failed to post target transaction: %w", err)
	}

	// Update intercompany transaction with transaction IDs
	intercompanyTxn.SourceTransactionID = sourceTxn.ID
	intercompanyTxn.TargetTransactionID = targetTxn.ID
	intercompanyTxn.MatchingStatus = IntercompanyMatched

	// Save intercompany transaction
	if err := mce.storage.SaveIntercompanyTransaction(intercompanyTxn); err != nil {
		return nil, fmt.Errorf("failed to save intercompany transaction: %w", err)
	}

	return intercompanyTxn, nil
}

// ReconcileIntercompanyTransactions reconciles pending intercompany transactions
func (mce *MultiCompanyEngine) ReconcileIntercompanyTransactions(companyID string, userID string) ([]*IntercompanyTransaction, error) {
	// Get all pending intercompany transactions for this company
	transactions, err := mce.storage.GetIntercompanyTransactionsByCompany(companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get intercompany transactions: %w", err)
	}

	reconciledTransactions := make([]*IntercompanyTransaction, 0)

	for _, txn := range transactions {
		if txn.MatchingStatus == IntercompanyMatched {
			// Mark as reconciled
			txn.MatchingStatus = IntercompanyReconciled
			now := time.Now()
			txn.ReconciledAt = &now
			txn.ReconciledBy = userID

			if err := mce.storage.SaveIntercompanyTransaction(txn); err != nil {
				continue // Skip on error
			}

			reconciledTransactions = append(reconciledTransactions, txn)
		}
	}

	return reconciledTransactions, nil
}

// GenerateConsolidatedTrialBalance generates a consolidated trial balance
func (mce *MultiCompanyEngine) GenerateConsolidatedTrialBalance(
	groupID string,
	asOfDate time.Time) (*ConsolidatedTrialBalance, error) {

	group, err := mce.storage.GetConsolidationGroup(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get consolidation group: %w", err)
	}

	consolidatedTB := &ConsolidatedTrialBalance{
		GroupName: group.Name,
		AsOfDate:  asOfDate,
		Companies: make(map[string]*TrialBalance),
	}

	// Get trial balance for each company
	for _, companyID := range group.ChildCompanies {
		engine, err := mce.GetAccountingEngine(companyID)
		if err != nil {
			continue // Skip on error
		}

		trialBalance, err := engine.GetTrialBalance(asOfDate, nil)
		if err != nil {
			continue // Skip on error
		}

		company, _ := mce.GetCompany(companyID)
		consolidatedTB.Companies[companyID] = &TrialBalance{
			CompanyName: company.Name,
			Balances:    trialBalance,
		}
	}

	// Apply elimination rules
	consolidatedTB.ConsolidatedBalances = mce.applyEliminationRules(consolidatedTB, group.EliminationRules)

	return consolidatedTB, nil
}

// ConsolidatedTrialBalance represents a consolidated trial balance
type ConsolidatedTrialBalance struct {
	GroupName            string                   `json:"group_name"`
	AsOfDate             time.Time                `json:"as_of_date"`
	Companies            map[string]*TrialBalance `json:"companies"`
	ConsolidatedBalances []*BalanceResult         `json:"consolidated_balances"`
	EliminationEntries   []*EliminationEntry      `json:"elimination_entries"`
}

// TrialBalance represents a company's trial balance
type TrialBalance struct {
	CompanyName string           `json:"company_name"`
	Balances    []*BalanceResult `json:"balances"`
}

// EliminationEntry represents an elimination entry for consolidation
type EliminationEntry struct {
	Description string  `json:"description"`
	AccountID   string  `json:"account_id"`
	Amount      *Amount `json:"amount"`
	RuleID      string  `json:"rule_id"`
}

// applyEliminationRules applies consolidation elimination rules
func (mce *MultiCompanyEngine) applyEliminationRules(
	consolidatedTB *ConsolidatedTrialBalance,
	rules []*EliminationRule) []*BalanceResult {

	// Combine all company balances
	combinedBalances := make(map[string]*BalanceResult)

	for _, companyTB := range consolidatedTB.Companies {
		for _, balance := range companyTB.Balances {
			if existing, exists := combinedBalances[balance.AccountID]; exists {
				existing.Balance.Value += balance.Balance.Value
			} else {
				combinedBalances[balance.AccountID] = &BalanceResult{
					AccountID:   balance.AccountID,
					AccountName: balance.AccountName,
					AccountType: balance.AccountType,
					Balance: &Amount{
						Value:    balance.Balance.Value,
						Currency: balance.Balance.Currency,
					},
				}
			}
		}
	}

	// Apply elimination rules
	for _, rule := range rules {
		if !rule.IsActive {
			continue
		}

		switch rule.RuleType {
		case "INTERCOMPANY_SALES":
			mce.eliminateIntercompanySales(combinedBalances, rule)
		case "INVESTMENT":
			mce.eliminateInvestments(combinedBalances, rule)
		}
	}

	// Convert map to slice
	result := make([]*BalanceResult, 0, len(combinedBalances))
	for _, balance := range combinedBalances {
		result = append(result, balance)
	}

	return result
}

// eliminateIntercompanySales eliminates intercompany sales
func (mce *MultiCompanyEngine) eliminateIntercompanySales(
	balances map[string]*BalanceResult,
	rule *EliminationRule) {

	// Look for intercompany receivables and payables
	for _, accountID := range rule.Accounts {
		if balance, exists := balances[accountID]; exists {
			if accountID == "intercompany_receivable" || accountID == "intercompany_payable" {
				// Zero out intercompany balances
				balance.Balance.Value = 0
			}
		}
	}
}

// eliminateInvestments eliminates investment accounts for subsidiaries
func (mce *MultiCompanyEngine) eliminateInvestments(
	balances map[string]*BalanceResult,
	rule *EliminationRule) {

	// Implementation depends on specific investment elimination rules
	for _, accountID := range rule.Accounts {
		if balance, exists := balances[accountID]; exists {
			if accountID == "investment_in_subsidiary" {
				// Eliminate against subsidiary equity
				balance.Balance.Value = 0
			}
		}
	}
}

// GetIntercompanyTransactions gets all intercompany transactions for a company
func (mce *MultiCompanyEngine) GetIntercompanyTransactions(companyID string) ([]*IntercompanyTransaction, error) {
	return mce.storage.GetIntercompanyTransactionsByCompany(companyID)
}

// CreateConsolidationGroup creates a new consolidation group
func (mce *MultiCompanyEngine) CreateConsolidationGroup(group *ConsolidationGroup, userID string) error {
	group.CreatedAt = time.Now()
	group.CreatedBy = userID

	return mce.storage.SaveConsolidationGroup(group)
}
