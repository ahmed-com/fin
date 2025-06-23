package accounting

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AccountingEngine is the main entry point for the accounting system
type AccountingEngine struct {
	storage               *Storage
	eventStore            *EventStore
	processor             *EventProcessor
	postingEngine         *PostingEngine
	queryAPI              *QueryAPI
	reconciliationService *ReconciliationService
	accrualService        *AccrualService
	reportingService      *ReportingService  // Add reporting service
	zbbService            *ZBBService        // Add ZBB service
	complianceService     *ComplianceService // Add compliance service
	amlService            *AMLService        // Add AML service
	forensicService       *ForensicService   // Add forensic service
}

// NewAccountingEngine creates a new accounting engine
func NewAccountingEngine(dbPath string) (*AccountingEngine, error) {
	// Initialize storage
	storage, err := NewStorage(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Initialize event store and processor
	eventStore := NewEventStore(storage)
	processor := NewEventProcessor(storage)

	// Initialize posting engine
	postingEngine := NewPostingEngine(storage, eventStore, processor)

	// Initialize query API
	queryAPI := NewQueryAPI(storage, postingEngine)

	// Initialize services
	reconciliationService := NewReconciliationService(storage, queryAPI)
	accrualService := NewAccrualService(storage, postingEngine, eventStore)
	reportingService := NewReportingService(storage, queryAPI)
	zbbService := NewZBBService(storage)                                     // Add ZBB service
	complianceService := NewComplianceService(*storage)                      // Add compliance service (dereference)
	forensicService := NewForensicService(storage, eventStore)               // Add forensic service
	amlService := NewAMLService(storage, complianceService, forensicService) // Add AML service

	return &AccountingEngine{
		storage:               storage,
		eventStore:            eventStore,
		processor:             processor,
		postingEngine:         postingEngine,
		queryAPI:              queryAPI,
		reconciliationService: reconciliationService,
		accrualService:        accrualService,
		reportingService:      reportingService,  // Add reporting service
		zbbService:            zbbService,        // Add ZBB service
		complianceService:     complianceService, // Add compliance service
		amlService:            amlService,        // Add AML service
		forensicService:       forensicService,   // Add forensic service
	}, nil
}

// Close closes the accounting engine and releases resources
func (ae *AccountingEngine) Close() error {
	return ae.storage.Close()
}

// CreateAccount creates a new account
func (ae *AccountingEngine) CreateAccount(account *Account, userID string) error {
	// Set timestamps
	account.CreatedAt = time.Now()
	if account.ID == "" {
		account.ID = uuid.New().String()
	}

	// Create account creation event
	_, err := ae.eventStore.CreateEvent(
		EventCreateAccount,
		AccountCreatedEvent{Account: account},
		time.Now(),
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to create account event: %w", err)
	}

	// Save account
	return ae.storage.SaveAccount(account)
}

// CreateTransaction creates a new transaction
func (ae *AccountingEngine) CreateTransaction(txn *Transaction, userID string) error {
	// Set timestamps and IDs
	if txn.ID == "" {
		txn.ID = uuid.New().String()
	}
	txn.CreatedAt = time.Now()
	txn.UpdatedAt = time.Now()
	txn.Status = Pending

	// Generate entry IDs
	for i := range txn.Entries {
		if txn.Entries[i].ID == "" {
			txn.Entries[i].ID = uuid.New().String()
		}
		txn.Entries[i].TransactionID = txn.ID
	}

	// Create transaction creation event
	_, err := ae.eventStore.CreateEvent(
		EventCreateTransaction,
		TransactionCreatedEvent{Transaction: txn},
		txn.ValidTime,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction event: %w", err)
	}

	// Process the event
	return ae.storage.SaveTransaction(txn)
}

// PostTransaction posts a transaction to the ledger
func (ae *AccountingEngine) PostTransaction(txnID string, userID string) error {
	txn, err := ae.storage.GetTransaction(txnID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	return ae.postingEngine.PostTransaction(txn, userID)
}

// GetAccountBalance gets the current balance of an account
func (ae *AccountingEngine) GetAccountBalance(accountID string, asOfDate time.Time) (*BalanceResult, error) {
	return ae.queryAPI.GetAccountBalance(accountID, asOfDate)
}

// GetTrialBalance generates a trial balance report
func (ae *AccountingEngine) GetTrialBalance(asOfDate time.Time, accountTypes []AccountType) ([]*BalanceResult, error) {
	return ae.queryAPI.GetTrialBalance(asOfDate, accountTypes)
}

// CreatePeriod creates a new accounting period
func (ae *AccountingEngine) CreatePeriod(period *Period, userID string) error {
	if period.ID == "" {
		period.ID = uuid.New().String()
	}

	// Create period creation event
	_, err := ae.eventStore.CreateEvent(
		EventCreatePeriod,
		period,
		time.Now(),
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to create period event: %w", err)
	}

	return ae.storage.SavePeriod(period)
}

// ClosePeriod closes an accounting period
func (ae *AccountingEngine) ClosePeriod(periodID string, softClose bool, userID string) error {
	period, err := ae.storage.GetPeriod(periodID)
	if err != nil {
		return fmt.Errorf("failed to get period: %w", err)
	}

	now := time.Now()
	if softClose {
		period.SoftClosedAt = &now
	} else {
		period.HardClosedAt = &now
	}

	// Create period close event
	_, err = ae.eventStore.CreateEvent(
		EventClosePeriod,
		map[string]interface{}{
			"period_id":  periodID,
			"soft_close": softClose,
			"closed_at":  now,
		},
		time.Now(),
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to create period close event: %w", err)
	}

	return ae.storage.SavePeriod(period)
}

// AutoReconcile performs automatic reconciliation
func (ae *AccountingEngine) AutoReconcile(accountID string, statements []*ExternalStatement) ([]*ReconciliationMatch, error) {
	return ae.reconciliationService.AutoReconcile(accountID, statements)
}

// ConfirmReconciliation confirms a reconciliation match
func (ae *AccountingEngine) ConfirmReconciliation(match *ReconciliationMatch, userID string) (*Reconciliation, error) {
	return ae.reconciliationService.ConfirmReconciliation(match, userID)
}

// CreateAccrualSchedule creates a new accrual/deferral schedule
func (ae *AccountingEngine) CreateAccrualSchedule(
	txnID string,
	totalAmount *Amount,
	frequency ScheduleFrequency,
	occurrences int,
	startDate time.Time,
	template *AccrualTemplate,
	userID string,
) (*RecognitionSchedule, error) {
	return ae.accrualService.CreateRecognitionSchedule(
		txnID, totalAmount, frequency, occurrences, startDate, template, userID,
	)
}

// ProcessAccruals processes pending accrual recognitions
func (ae *AccountingEngine) ProcessAccruals(upToDate time.Time, userID string) error {
	return ae.accrualService.ProcessPendingRecognitions(upToDate, userID)
}

// GetReconciliationSummary gets reconciliation summary for an account
func (ae *AccountingEngine) GetReconciliationSummary(accountID string, asOfDate time.Time) (*ReconciliationSummary, error) {
	return ae.reconciliationService.GetReconciliationSummary(accountID, asOfDate)
}

// ReverseTransaction creates a reversing transaction
func (ae *AccountingEngine) ReverseTransaction(originalTxnID string, description string, userID string) (*Transaction, error) {
	return ae.postingEngine.ReverseTransaction(originalTxnID, description, userID)
}

// GetEvents retrieves events within a time range
func (ae *AccountingEngine) GetEvents(from, to time.Time) ([]*JournalEvent, error) {
	return ae.eventStore.GetEvents(from, to)
}

// Financial Reporting Methods

// GenerateBalanceSheet generates a balance sheet as of a specific date
func (ae *AccountingEngine) GenerateBalanceSheet(asOfDate time.Time, currency string) (*FinancialStatement, error) {
	return ae.reportingService.GenerateBalanceSheet(asOfDate, currency)
}

// GenerateProfitAndLoss generates a P&L statement for a period
func (ae *AccountingEngine) GenerateProfitAndLoss(fromDate, toDate time.Time, currency string) (*FinancialStatement, error) {
	return ae.reportingService.GenerateProfitAndLoss(fromDate, toDate, currency)
}

// GenerateCashFlowStatement generates a cash flow statement for a period
func (ae *AccountingEngine) GenerateCashFlowStatement(fromDate, toDate time.Time, currency string) (*CashFlowStatement, error) {
	return ae.reportingService.GenerateCashFlowStatement(fromDate, toDate, currency)
}

// FormatFinancialStatement formats a financial statement for display
func (ae *AccountingEngine) FormatFinancialStatement(statement *FinancialStatement) string {
	return ae.reportingService.FormatFinancialStatement(statement)
}

// FormatCashFlowStatement formats a cash flow statement for display
func (ae *AccountingEngine) FormatCashFlowStatement(cf *CashFlowStatement) string {
	return ae.reportingService.FormatCashFlowStatement(cf)
}

// ----------------------------------------------------------------------------
// Zero-Based Budgeting Methods
// ----------------------------------------------------------------------------

// CreateBudgetPeriod creates a new budget period
func (ae *AccountingEngine) CreateBudgetPeriod(period *BudgetPeriod, userID string) error {
	return ae.zbbService.CreateBudgetPeriod(period, userID)
}

// CreateBudgetRequest creates a new zero-based budget request
func (ae *AccountingEngine) CreateBudgetRequest(request *BudgetRequest, userID string) error {
	return ae.zbbService.CreateBudgetRequest(request, userID)
}

// AddJustification adds business justification to a budget request
func (ae *AccountingEngine) AddJustification(requestID string, justification *Justification, userID string) error {
	return ae.zbbService.AddJustification(requestID, justification, userID)
}

// SubmitBudgetRequest submits request for approval
func (ae *AccountingEngine) SubmitBudgetRequest(requestID string, userID string) error {
	return ae.zbbService.SubmitBudgetRequest(requestID, userID)
}

// ApproveBudgetRequest approves a budget request
func (ae *AccountingEngine) ApproveBudgetRequest(requestID string, approverID string, approvedAmount *Amount, comments string) error {
	return ae.zbbService.ApproveBudgetRequest(requestID, approverID, approvedAmount, comments)
}

// CreateBudgetAllocation creates budget allocation from approved request
func (ae *AccountingEngine) CreateBudgetAllocation(requestID string, userID string) error {
	return ae.zbbService.CreateBudgetAllocation(requestID, userID)
}

// TrackBudgetSpending tracks actual spending against budget allocations
func (ae *AccountingEngine) TrackBudgetSpending(transactionID string, allocationID string) error {
	return ae.zbbService.TrackBudgetSpending(transactionID, allocationID)
}

// GetBudgetVariance calculates variance between budget and actual
func (ae *AccountingEngine) GetBudgetVariance(periodID string, departmentID string) (*BudgetVarianceReport, error) {
	return ae.zbbService.GetBudgetVariance(periodID, departmentID)
}

// GetDepartmentBudgetSummary provides summary of all budget requests for a department
func (ae *AccountingEngine) GetDepartmentBudgetSummary(periodID string, departmentID string) (*DepartmentBudgetSummary, error) {
	return ae.zbbService.GetDepartmentBudgetSummary(periodID, departmentID)
}

// Helper methods for common operations

// CreateStandardAccounts creates a basic chart of accounts
func (ae *AccountingEngine) CreateStandardAccounts(userID string) error {
	accounts := []*Account{
		{
			ID:   "cash",
			Code: "1001",
			Name: "Cash",
			Type: Asset,
		},
		{
			ID:   "accounts_receivable",
			Code: "1200",
			Name: "Accounts Receivable",
			Type: Asset,
		},
		{
			ID:   "intercompany_receivable",
			Code: "1300",
			Name: "Intercompany Receivable",
			Type: Asset,
		},
		{
			ID:   "accounts_payable",
			Code: "2001",
			Name: "Accounts Payable",
			Type: Liability,
		},
		{
			ID:   "intercompany_payable",
			Code: "2200",
			Name: "Intercompany Payable",
			Type: Liability,
		},
		{
			ID:   "revenue",
			Code: "4001",
			Name: "Revenue",
			Type: Income,
		},
		{
			ID:   "expenses",
			Code: "5001",
			Name: "Operating Expenses",
			Type: Expense,
		},
		{
			ID:   "unearned_revenue",
			Code: "2100",
			Name: "Unearned Revenue",
			Type: Liability,
		},
	}

	for _, account := range accounts {
		if err := ae.CreateAccount(account, userID); err != nil {
			return fmt.Errorf("failed to create account %s: %w", account.Name, err)
		}
	}

	return nil
}

// CreateLedger creates a new ledger
func (ae *AccountingEngine) CreateLedger(ledger *Ledger) error {
	if ledger.ID == "" {
		ledger.ID = uuid.New().String()
	}
	return ae.storage.SaveLedger(ledger)
}

// ----------------------------------------------------------------------------
// Service Getters
// ----------------------------------------------------------------------------

// GetAMLService returns the AML service
func (ae *AccountingEngine) GetAMLService() *AMLService {
	return ae.amlService
}

// GetComplianceService returns the compliance service
func (ae *AccountingEngine) GetComplianceService() *ComplianceService {
	return ae.complianceService
}

// GetForensicService returns the forensic service
func (ae *AccountingEngine) GetForensicService() *ForensicService {
	return ae.forensicService
}
