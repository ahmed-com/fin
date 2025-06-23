package accounting

import (
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

// Storage buckets
var (
	BucketEvents                   = []byte("events")
	BucketAccounts                 = []byte("accounts")
	BucketTransactions             = []byte("transactions")
	BucketEntries                  = []byte("entries")
	BucketLedgers                  = []byte("ledgers")
	BucketPeriods                  = []byte("periods")
	BucketReconciliations          = []byte("reconciliations")
	BucketSchedules                = []byte("schedules")
	BucketReportingContexts        = []byte("reporting_contexts")
	BucketCompanies                = []byte("companies")
	BucketIntercompanyTransactions = []byte("intercompany_transactions")
	BucketConsolidationGroups      = []byte("consolidation_groups")
	// Zero-Based Budgeting buckets
	BucketBudgetPeriods     = []byte("budget_periods")
	BucketBudgetRequests    = []byte("budget_requests")
	BucketBudgetApprovals   = []byte("budget_approvals")
	BucketBudgetAllocations = []byte("budget_allocations")
	BucketBudgetTracking    = []byte("budget_tracking")
	// Compliance buckets
	BucketComplianceRules      = []byte("compliance_rules")
	BucketTaxRules             = []byte("tax_rules")
	BucketComplianceViolations = []byte("compliance_violations")
	BucketTaxReturns           = []byte("tax_returns")
	// AML buckets
	BucketAMLRules     = []byte("aml_rules")
	BucketAMLAlerts    = []byte("aml_alerts")
	BucketAMLCustomers = []byte("aml_customers")
)

// Storage provides persistent storage for the accounting system
type Storage struct {
	db *bbolt.DB
}

// NewStorage creates a new storage instance
func NewStorage(dbPath string) (*Storage, error) {
	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	storage := &Storage{db: db}
	if err := storage.initBuckets(); err != nil {
		return nil, fmt.Errorf("failed to initialize buckets: %w", err)
	}

	return storage, nil
}

// Close closes the database connection
func (s *Storage) Close() error {
	return s.db.Close()
}

// initBuckets creates all required buckets
func (s *Storage) initBuckets() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		buckets := [][]byte{
			BucketEvents, BucketAccounts, BucketTransactions,
			BucketEntries, BucketLedgers, BucketPeriods,
			BucketReconciliations, BucketSchedules, BucketReportingContexts,
			BucketCompanies, BucketIntercompanyTransactions, BucketConsolidationGroups,
			// Zero-Based Budgeting buckets
			BucketBudgetPeriods, BucketBudgetRequests, BucketBudgetApprovals,
			BucketBudgetAllocations, BucketBudgetTracking,
			// Compliance buckets
			BucketComplianceRules, BucketTaxRules, BucketComplianceViolations, BucketTaxReturns,
			// AML buckets
			BucketAMLRules, BucketAMLAlerts, BucketAMLCustomers,
		}

		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return fmt.Errorf("failed to create bucket %s: %w", bucket, err)
			}
		}
		return nil
	})
}

// AppendEvent appends a new event to the event log
func (s *Storage) AppendEvent(event *JournalEvent) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketEvents)
		data, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		// Use timestamp + ID as key for ordering
		key := fmt.Sprintf("%d_%s", event.TransactionTime.UnixNano(), event.ID)
		return b.Put([]byte(key), data)
	})
}

// GetEvents retrieves events within a time range
func (s *Storage) GetEvents(from, to time.Time) ([]*JournalEvent, error) {
	var events []*JournalEvent

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketEvents)
		c := b.Cursor()

		fromKey := []byte(fmt.Sprintf("%d", from.UnixNano()))
		toKey := []byte(fmt.Sprintf("%d", to.UnixNano()))

		for k, v := c.Seek(fromKey); k != nil && string(k) <= string(toKey); k, v = c.Next() {
			var event JournalEvent
			if err := json.Unmarshal(v, &event); err != nil {
				return fmt.Errorf("failed to unmarshal event: %w", err)
			}
			events = append(events, &event)
		}
		return nil
	})

	return events, err
}

// SaveAccount saves an account to storage
func (s *Storage) SaveAccount(account *Account) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketAccounts)
		data, err := json.Marshal(account)
		if err != nil {
			return fmt.Errorf("failed to marshal account: %w", err)
		}
		return b.Put([]byte(account.ID), data)
	})
}

// GetAccount retrieves an account by ID
func (s *Storage) GetAccount(id string) (*Account, error) {
	var account Account

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketAccounts)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("account not found: %s", id)
		}
		return json.Unmarshal(data, &account)
	})

	if err != nil {
		return nil, err
	}
	return &account, nil
}

// SaveTransaction saves a transaction to storage
func (s *Storage) SaveTransaction(txn *Transaction) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketTransactions)
		data, err := json.Marshal(txn)
		if err != nil {
			return fmt.Errorf("failed to marshal transaction: %w", err)
		}
		return b.Put([]byte(txn.ID), data)
	})
}

// GetTransaction retrieves a transaction by ID
func (s *Storage) GetTransaction(id string) (*Transaction, error) {
	var txn Transaction

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketTransactions)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("transaction not found: %s", id)
		}
		return json.Unmarshal(data, &txn)
	})

	if err != nil {
		return nil, err
	}
	return &txn, nil
}

// SaveEntry saves an entry to storage
func (s *Storage) SaveEntry(entry *Entry) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketEntries)
		data, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("failed to marshal entry: %w", err)
		}
		return b.Put([]byte(entry.ID), data)
	})
}

// GetEntriesByAccount retrieves all entries for a specific account
func (s *Storage) GetEntriesByAccount(accountID string) ([]*Entry, error) {
	var entries []*Entry

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketEntries)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var entry Entry
			if err := json.Unmarshal(v, &entry); err != nil {
				return fmt.Errorf("failed to unmarshal entry: %w", err)
			}
			if entry.AccountID == accountID {
				entries = append(entries, &entry)
			}
		}
		return nil
	})

	return entries, err
}

// SaveLedger saves a ledger to storage
func (s *Storage) SaveLedger(ledger *Ledger) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketLedgers)
		data, err := json.Marshal(ledger)
		if err != nil {
			return fmt.Errorf("failed to marshal ledger: %w", err)
		}
		return b.Put([]byte(ledger.ID), data)
	})
}

// SavePeriod saves a period to storage
func (s *Storage) SavePeriod(period *Period) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketPeriods)
		data, err := json.Marshal(period)
		if err != nil {
			return fmt.Errorf("failed to marshal period: %w", err)
		}
		return b.Put([]byte(period.ID), data)
	})
}

// GetPeriod retrieves a period by ID
func (s *Storage) GetPeriod(id string) (*Period, error) {
	var period Period

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketPeriods)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("period not found: %s", id)
		}
		return json.Unmarshal(data, &period)
	})

	if err != nil {
		return nil, err
	}
	return &period, nil
}

// SaveReconciliation saves a reconciliation to storage
func (s *Storage) SaveReconciliation(recon *Reconciliation) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketReconciliations)
		data, err := json.Marshal(recon)
		if err != nil {
			return fmt.Errorf("failed to marshal reconciliation: %w", err)
		}
		return b.Put([]byte(recon.ID), data)
	})
}

// SaveSchedule saves a recognition schedule to storage
func (s *Storage) SaveSchedule(schedule *RecognitionSchedule) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketSchedules)
		data, err := json.Marshal(schedule)
		if err != nil {
			return fmt.Errorf("failed to marshal schedule: %w", err)
		}
		return b.Put([]byte(schedule.ID), data)
	})
}

// GetAllSchedules retrieves all recognition schedules
func (s *Storage) GetAllSchedules() ([]*RecognitionSchedule, error) {
	var schedules []*RecognitionSchedule

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketSchedules)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var schedule RecognitionSchedule
			if err := json.Unmarshal(v, &schedule); err != nil {
				return fmt.Errorf("failed to unmarshal schedule: %w", err)
			}
			schedules = append(schedules, &schedule)
		}
		return nil
	})

	return schedules, err
}

// SaveCompany saves a company to storage
func (s *Storage) SaveCompany(company *Company) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketCompanies)
		data, err := json.Marshal(company)
		if err != nil {
			return fmt.Errorf("failed to marshal company: %w", err)
		}
		return b.Put([]byte(company.ID), data)
	})
}

// GetCompany retrieves a company by ID
func (s *Storage) GetCompany(id string) (*Company, error) {
	var company Company

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketCompanies)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("company not found: %s", id)
		}
		return json.Unmarshal(data, &company)
	})

	if err != nil {
		return nil, err
	}
	return &company, nil
}

// GetCompanies retrieves all companies
func (s *Storage) GetCompanies() ([]*Company, error) {
	var companies []*Company

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketCompanies)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var company Company
			if err := json.Unmarshal(v, &company); err != nil {
				return fmt.Errorf("failed to unmarshal company: %w", err)
			}
			companies = append(companies, &company)
		}
		return nil
	})

	return companies, err
}

// SaveIntercompanyTransaction saves an intercompany transaction to storage
func (s *Storage) SaveIntercompanyTransaction(txn *IntercompanyTransaction) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketIntercompanyTransactions)
		data, err := json.Marshal(txn)
		if err != nil {
			return fmt.Errorf("failed to marshal intercompany transaction: %w", err)
		}
		return b.Put([]byte(txn.ID), data)
	})
}

// GetIntercompanyTransaction retrieves an intercompany transaction by ID
func (s *Storage) GetIntercompanyTransaction(id string) (*IntercompanyTransaction, error) {
	var txn IntercompanyTransaction

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketIntercompanyTransactions)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("intercompany transaction not found: %s", id)
		}
		return json.Unmarshal(data, &txn)
	})

	if err != nil {
		return nil, err
	}
	return &txn, nil
}

// GetIntercompanyTransactionsByCompany retrieves all intercompany transactions for a specific company
func (s *Storage) GetIntercompanyTransactionsByCompany(companyID string) ([]*IntercompanyTransaction, error) {
	var txns []*IntercompanyTransaction

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketIntercompanyTransactions)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var txn IntercompanyTransaction
			if err := json.Unmarshal(v, &txn); err != nil {
				return fmt.Errorf("failed to unmarshal intercompany transaction: %w", err)
			}
			if txn.SourceCompanyID == companyID || txn.TargetCompanyID == companyID {
				txns = append(txns, &txn)
			}
		}
		return nil
	})

	return txns, err
}

// SaveConsolidationGroup saves a consolidation group to storage
func (s *Storage) SaveConsolidationGroup(group *ConsolidationGroup) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketConsolidationGroups)
		data, err := json.Marshal(group)
		if err != nil {
			return fmt.Errorf("failed to marshal consolidation group: %w", err)
		}
		return b.Put([]byte(group.ID), data)
	})
}

// GetConsolidationGroup retrieves a consolidation group by ID
func (s *Storage) GetConsolidationGroup(id string) (*ConsolidationGroup, error) {
	var group ConsolidationGroup

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketConsolidationGroups)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("consolidation group not found: %s", id)
		}
		return json.Unmarshal(data, &group)
	})

	if err != nil {
		return nil, err
	}
	return &group, nil
}

// GetConsolidationGroups retrieves all consolidation groups
func (s *Storage) GetConsolidationGroups() ([]*ConsolidationGroup, error) {
	var groups []*ConsolidationGroup

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketConsolidationGroups)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var group ConsolidationGroup
			if err := json.Unmarshal(v, &group); err != nil {
				return fmt.Errorf("failed to unmarshal consolidation group: %w", err)
			}
			groups = append(groups, &group)
		}
		return nil
	})

	return groups, err
}

// ----------------------------------------------------------------------------
// Zero-Based Budgeting Storage Methods
// ----------------------------------------------------------------------------

// SaveBudgetPeriod saves a budget period
func (s *Storage) SaveBudgetPeriod(period *BudgetPeriod) error {
	data, err := json.Marshal(period)
	if err != nil {
		return fmt.Errorf("failed to marshal budget period: %w", err)
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketBudgetPeriods)
		return b.Put([]byte(period.ID), data)
	})
}

// GetBudgetPeriod retrieves a budget period by ID
func (s *Storage) GetBudgetPeriod(id string) (*BudgetPeriod, error) {
	var period BudgetPeriod

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketBudgetPeriods)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("budget period not found")
		}
		return json.Unmarshal(data, &period)
	})

	return &period, err
}

// SaveBudgetRequest saves a budget request
func (s *Storage) SaveBudgetRequest(request *BudgetRequest) error {
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal budget request: %w", err)
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketBudgetRequests)
		return b.Put([]byte(request.ID), data)
	})
}

// GetBudgetRequest retrieves a budget request by ID
func (s *Storage) GetBudgetRequest(id string) (*BudgetRequest, error) {
	var request BudgetRequest

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketBudgetRequests)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("budget request not found")
		}
		return json.Unmarshal(data, &request)
	})

	return &request, err
}

// GetBudgetRequestsByPeriodAndDept retrieves budget requests by period and department
func (s *Storage) GetBudgetRequestsByPeriodAndDept(periodID, departmentID string) ([]*BudgetRequest, error) {
	var requests []*BudgetRequest

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketBudgetRequests)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var request BudgetRequest
			if err := json.Unmarshal(v, &request); err != nil {
				continue // Skip malformed requests
			}

			if request.PeriodID == periodID && request.DepartmentID == departmentID {
				requests = append(requests, &request)
			}
		}
		return nil
	})

	return requests, err
}

// SaveBudgetApproval saves a budget approval
func (s *Storage) SaveBudgetApproval(approval *BudgetApproval) error {
	data, err := json.Marshal(approval)
	if err != nil {
		return fmt.Errorf("failed to marshal budget approval: %w", err)
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketBudgetApprovals)
		return b.Put([]byte(approval.ID), data)
	})
}

// GetBudgetApprovalsByRequest retrieves budget approvals by request ID
func (s *Storage) GetBudgetApprovalsByRequest(requestID string) ([]*BudgetApproval, error) {
	var approvals []*BudgetApproval

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketBudgetApprovals)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var approval BudgetApproval
			if err := json.Unmarshal(v, &approval); err != nil {
				continue // Skip malformed approvals
			}

			if approval.RequestID == requestID {
				approvals = append(approvals, &approval)
			}
		}
		return nil
	})

	return approvals, err
}

// SaveBudgetAllocation saves a budget allocation
func (s *Storage) SaveBudgetAllocation(allocation *BudgetAllocation) error {
	data, err := json.Marshal(allocation)
	if err != nil {
		return fmt.Errorf("failed to marshal budget allocation: %w", err)
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketBudgetAllocations)
		return b.Put([]byte(allocation.ID), data)
	})
}

// GetBudgetAllocation retrieves a budget allocation by ID
func (s *Storage) GetBudgetAllocation(id string) (*BudgetAllocation, error) {
	var allocation BudgetAllocation

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketBudgetAllocations)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("budget allocation not found")
		}
		return json.Unmarshal(data, &allocation)
	})

	return &allocation, err
}

// GetBudgetAllocationsByPeriodAndDept retrieves budget allocations by period and department
func (s *Storage) GetBudgetAllocationsByPeriodAndDept(periodID, departmentID string) ([]*BudgetAllocation, error) {
	var allocations []*BudgetAllocation

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketBudgetAllocations)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var allocation BudgetAllocation
			if err := json.Unmarshal(v, &allocation); err != nil {
				continue // Skip malformed allocations
			}

			if allocation.PeriodID == periodID && allocation.DepartmentID == departmentID {
				allocations = append(allocations, &allocation)
			}
		}
		return nil
	})

	return allocations, err
}

// SaveBudgetTracking saves budget tracking record
func (s *Storage) SaveBudgetTracking(tracking *BudgetTracking) error {
	data, err := json.Marshal(tracking)
	if err != nil {
		return fmt.Errorf("failed to marshal budget tracking: %w", err)
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketBudgetTracking)
		// Use composite key: allocationID + transactionID
		key := fmt.Sprintf("%s_%s", tracking.AllocationID, tracking.TransactionID)
		return b.Put([]byte(key), data)
	})
}

// SaveComplianceRule saves a compliance rule
func (s *Storage) SaveComplianceRule(rule *ComplianceRule) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketComplianceRules)
		data, err := json.Marshal(rule)
		if err != nil {
			return fmt.Errorf("failed to marshal compliance rule: %w", err)
		}
		return b.Put([]byte(rule.ID), data)
	})
}

// GetComplianceRule retrieves a compliance rule by ID
func (s *Storage) GetComplianceRule(id string) (*ComplianceRule, error) {
	var rule ComplianceRule

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketComplianceRules)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("compliance rule not found: %s", id)
		}
		return json.Unmarshal(data, &rule)
	})

	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// GetAllComplianceRules retrieves all compliance rules
func (s *Storage) GetAllComplianceRules() ([]*ComplianceRule, error) {
	var rules []*ComplianceRule

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketComplianceRules)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var rule ComplianceRule
			if err := json.Unmarshal(v, &rule); err != nil {
				return fmt.Errorf("failed to unmarshal compliance rule: %w", err)
			}
			rules = append(rules, &rule)
		}
		return nil
	})

	return rules, err
}

// SaveTaxRule saves a tax rule
func (s *Storage) SaveTaxRule(rule *TaxRule) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketTaxRules)
		data, err := json.Marshal(rule)
		if err != nil {
			return fmt.Errorf("failed to marshal tax rule: %w", err)
		}
		return b.Put([]byte(rule.ID), data)
	})
}

// GetTaxRule retrieves a tax rule by ID
func (s *Storage) GetTaxRule(id string) (*TaxRule, error) {
	var rule TaxRule

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketTaxRules)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("tax rule not found: %s", id)
		}
		return json.Unmarshal(data, &rule)
	})

	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// GetAllTaxRules retrieves all tax rules
func (s *Storage) GetAllTaxRules() ([]*TaxRule, error) {
	var rules []*TaxRule

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketTaxRules)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var rule TaxRule
			if err := json.Unmarshal(v, &rule); err != nil {
				return fmt.Errorf("failed to unmarshal tax rule: %w", err)
			}
			rules = append(rules, &rule)
		}
		return nil
	})

	return rules, err
}

// SaveComplianceViolation saves a compliance violation
func (s *Storage) SaveComplianceViolation(violation *ComplianceViolation) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketComplianceViolations)
		data, err := json.Marshal(violation)
		if err != nil {
			return fmt.Errorf("failed to marshal compliance violation: %w", err)
		}
		return b.Put([]byte(violation.ID), data)
	})
}

// GetComplianceViolation retrieves a compliance violation by ID
func (s *Storage) GetComplianceViolation(id string) (*ComplianceViolation, error) {
	var violation ComplianceViolation

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketComplianceViolations)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("compliance violation not found: %s", id)
		}
		return json.Unmarshal(data, &violation)
	})

	if err != nil {
		return nil, err
	}
	return &violation, nil
}

// GetAllComplianceViolations retrieves all compliance violations
func (s *Storage) GetAllComplianceViolations() ([]*ComplianceViolation, error) {
	var violations []*ComplianceViolation

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketComplianceViolations)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var violation ComplianceViolation
			if err := json.Unmarshal(v, &violation); err != nil {
				return fmt.Errorf("failed to unmarshal compliance violation: %w", err)
			}
			violations = append(violations, &violation)
		}
		return nil
	})

	return violations, err
}

// SaveTaxReturn saves a tax return
func (s *Storage) SaveTaxReturn(taxReturn *TaxReturn) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketTaxReturns)
		data, err := json.Marshal(taxReturn)
		if err != nil {
			return fmt.Errorf("failed to marshal tax return: %w", err)
		}
		return b.Put([]byte(taxReturn.ID), data)
	})
}

// GetTaxReturn retrieves a tax return by ID
func (s *Storage) GetTaxReturn(id string) (*TaxReturn, error) {
	var taxReturn TaxReturn

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketTaxReturns)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("tax return not found: %s", id)
		}
		return json.Unmarshal(data, &taxReturn)
	})

	if err != nil {
		return nil, err
	}
	return &taxReturn, nil
}

// GetAllTaxReturns retrieves all tax returns
func (s *Storage) GetAllTaxReturns() ([]*TaxReturn, error) {
	var taxReturns []*TaxReturn

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketTaxReturns)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var taxReturn TaxReturn
			if err := json.Unmarshal(v, &taxReturn); err != nil {
				return fmt.Errorf("failed to unmarshal tax return: %w", err)
			}
			taxReturns = append(taxReturns, &taxReturn)
		}
		return nil
	})

	return taxReturns, err
}

// GetTaxRulesByJurisdiction retrieves tax rules by jurisdiction and tax type
func (s *Storage) GetTaxRulesByJurisdiction(jurisdiction TaxJurisdiction, taxType TaxType) ([]*TaxRule, error) {
	var rules []*TaxRule

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketTaxRules)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var rule TaxRule
			if err := json.Unmarshal(v, &rule); err != nil {
				return fmt.Errorf("failed to unmarshal tax rule: %w", err)
			}

			// Filter by jurisdiction and tax type
			if rule.Jurisdiction == jurisdiction && rule.TaxType == taxType {
				rules = append(rules, &rule)
			}
		}
		return nil
	})

	return rules, err
}

// GetTransactionsByDateRange retrieves transactions within a date range for a company
func (s *Storage) GetTransactionsByDateRange(companyID string, startDate, endDate time.Time) ([]*Transaction, error) {
	var transactions []*Transaction

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketTransactions)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var txn Transaction
			if err := json.Unmarshal(v, &txn); err != nil {
				continue // Skip malformed transactions
			}

			// Filter by company and date range
			// Note: In a real implementation, we'd need a company field on Transaction
			// For now, we'll just filter by date
			if txn.ValidTime.After(startDate) || txn.ValidTime.Equal(startDate) {
				if txn.ValidTime.Before(endDate) || txn.ValidTime.Equal(endDate) {
					transactions = append(transactions, &txn)
				}
			}
		}
		return nil
	})

	return transactions, err
}

// GetComplianceViolations retrieves compliance violations for a company
func (s *Storage) GetComplianceViolations(companyID string) ([]ComplianceViolation, error) {
	var violations []ComplianceViolation

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketComplianceViolations)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var violation ComplianceViolation
			if err := json.Unmarshal(v, &violation); err != nil {
				return fmt.Errorf("failed to unmarshal compliance violation: %w", err)
			}
			violations = append(violations, violation)
		}
		return nil
	})

	return violations, err
}

// QueryEntries queries entries based on provided options and filters
func (s *Storage) QueryEntries(options *QueryOptions) ([]*Entry, error) {
	var entries []*Entry

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketEntries)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var entry Entry
			if err := json.Unmarshal(v, &entry); err != nil {
				continue // Skip malformed entries
			}

			// Apply filters
			matches := true
			for _, filter := range options.Filters {
				switch filter.Field {
				case "company_id":
					// Note: In a real implementation, we'd need a company field on Entry
					// For now, we'll skip this filter or implement it differently
					continue
				case "account_id":
					if filter.Operator == "=" && entry.AccountID != fmt.Sprintf("%v", filter.Value) {
						matches = false
						break
					}
				case "transaction_id":
					if filter.Operator == "=" && entry.TransactionID != fmt.Sprintf("%v", filter.Value) {
						matches = false
						break
					}
				}
			}

			if matches {
				entries = append(entries, &entry)
			}
		}
		return nil
	})

	// Apply limit if specified
	if options.Limit > 0 && len(entries) > options.Limit {
		entries = entries[:options.Limit]
	}

	return entries, err
}

// ----------------------------------------------------------------------------
// AML Storage Methods
// ----------------------------------------------------------------------------

// SaveAMLRule saves an AML rule
func (s *Storage) SaveAMLRule(rule *AMLRule) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketAMLRules)
		data, err := json.Marshal(rule)
		if err != nil {
			return fmt.Errorf("failed to marshal AML rule: %w", err)
		}
		return b.Put([]byte(rule.ID), data)
	})
}

// GetAMLRule retrieves an AML rule by ID
func (s *Storage) GetAMLRule(id string) (*AMLRule, error) {
	var rule AMLRule

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketAMLRules)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("AML rule not found: %s", id)
		}
		return json.Unmarshal(data, &rule)
	})

	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// GetAllAMLRules retrieves all AML rules
func (s *Storage) GetAllAMLRules() ([]*AMLRule, error) {
	var rules []*AMLRule

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketAMLRules)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var rule AMLRule
			if err := json.Unmarshal(v, &rule); err != nil {
				return fmt.Errorf("failed to unmarshal AML rule: %w", err)
			}
			rules = append(rules, &rule)
		}
		return nil
	})

	return rules, err
}

// SaveAMLAlert saves an AML alert
func (s *Storage) SaveAMLAlert(alert *AMLAlert) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketAMLAlerts)
		data, err := json.Marshal(alert)
		if err != nil {
			return fmt.Errorf("failed to marshal AML alert: %w", err)
		}
		return b.Put([]byte(alert.ID), data)
	})
}

// GetAMLAlert retrieves an AML alert by ID
func (s *Storage) GetAMLAlert(id string) (*AMLAlert, error) {
	var alert AMLAlert

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketAMLAlerts)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("AML alert not found: %s", id)
		}
		return json.Unmarshal(data, &alert)
	})

	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// GetAMLAlerts retrieves all AML alerts
func (s *Storage) GetAMLAlerts() ([]*AMLAlert, error) {
	var alerts []*AMLAlert

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketAMLAlerts)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var alert AMLAlert
			if err := json.Unmarshal(v, &alert); err != nil {
				return fmt.Errorf("failed to unmarshal AML alert: %w", err)
			}
			alerts = append(alerts, &alert)
		}
		return nil
	})

	return alerts, err
}

// SaveAMLCustomer saves an AML customer
func (s *Storage) SaveAMLCustomer(customer *AMLCustomer) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketAMLCustomers)
		data, err := json.Marshal(customer)
		if err != nil {
			return fmt.Errorf("failed to marshal AML customer: %w", err)
		}
		return b.Put([]byte(customer.ID), data)
	})
}

// GetAMLCustomer retrieves an AML customer by ID
func (s *Storage) GetAMLCustomer(id string) (*AMLCustomer, error) {
	var customer AMLCustomer

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketAMLCustomers)
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("AML customer not found: %s", id)
		}
		return json.Unmarshal(data, &customer)
	})

	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// GetAllAMLCustomers retrieves all AML customers
func (s *Storage) GetAllAMLCustomers() ([]*AMLCustomer, error) {
	var customers []*AMLCustomer

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketAMLCustomers)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var customer AMLCustomer
			if err := json.Unmarshal(v, &customer); err != nil {
				return fmt.Errorf("failed to unmarshal AML customer: %w", err)
			}
			customers = append(customers, &customer)
		}
		return nil
	})

	return customers, err
}
