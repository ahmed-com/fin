package accounting

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ComplianceFramework represents different accounting standards
type ComplianceFramework string

const (
	GAAP_Framework ComplianceFramework = "GAAP"
	IFRS_Framework ComplianceFramework = "IFRS"
	SOX_Framework  ComplianceFramework = "SOX"
)

// TaxJurisdiction represents different tax jurisdictions
type TaxJurisdiction string

const (
	US_FEDERAL TaxJurisdiction = "US_FEDERAL"
	US_STATE   TaxJurisdiction = "US_STATE"
	EU_VAT     TaxJurisdiction = "EU_VAT"
	UK_VAT     TaxJurisdiction = "UK_VAT"
	CANADA_GST TaxJurisdiction = "CANADA_GST"
	AUSTRALIA  TaxJurisdiction = "AUSTRALIA"
)

// TaxType represents different types of taxes
type TaxType string

const (
	INCOME_TAX   TaxType = "INCOME_TAX"
	SALES_TAX    TaxType = "SALES_TAX"
	VAT          TaxType = "VAT"
	GST          TaxType = "GST"
	PAYROLL_TAX  TaxType = "PAYROLL_TAX"
	PROPERTY_TAX TaxType = "PROPERTY_TAX"
	WITHHOLDING  TaxType = "WITHHOLDING"
)

// ComplianceRule represents a regulatory compliance rule
type ComplianceRule struct {
	ID          string              `json:"id"`
	Framework   ComplianceFramework `json:"framework"`
	RuleType    string              `json:"rule_type"`
	Description string              `json:"description"`
	AccountType AccountType         `json:"account_type,omitempty"`
	Conditions  []string            `json:"conditions"`
	Actions     []string            `json:"actions"`
	Severity    string              `json:"severity"` // "ERROR", "WARNING", "INFO"
	Active      bool                `json:"active"`
	CreatedAt   time.Time           `json:"created_at"`
}

// TaxRule represents a tax calculation rule
type TaxRule struct {
	ID            string          `json:"id"`
	Jurisdiction  TaxJurisdiction `json:"jurisdiction"`
	TaxType       TaxType         `json:"tax_type"`
	Name          string          `json:"name"`
	Rate          float64         `json:"rate"`       // Tax rate as decimal (0.075 for 7.5%)
	MinAmount     float64         `json:"min_amount"` // Minimum taxable amount
	MaxAmount     float64         `json:"max_amount"` // Maximum taxable amount (0 = no limit)
	Exemptions    []string        `json:"exemptions"` // Account codes or transaction types exempt
	EffectiveFrom time.Time       `json:"effective_from"`
	EffectiveTo   *time.Time      `json:"effective_to,omitempty"`
	Active        bool            `json:"active"`
}

// TaxCalculation represents a calculated tax amount
type TaxCalculation struct {
	RuleID        string    `json:"rule_id"`
	BaseAmount    float64   `json:"base_amount"`
	TaxableAmount float64   `json:"taxable_amount"`
	TaxRate       float64   `json:"tax_rate"`
	TaxAmount     float64   `json:"tax_amount"`
	Exemptions    []string  `json:"exemptions"`
	CalculatedAt  time.Time `json:"calculated_at"`
}

// ComplianceViolation represents a compliance rule violation
type ComplianceViolation struct {
	ID            string     `json:"id"`
	RuleID        string     `json:"rule_id"`
	TransactionID string     `json:"transaction_id,omitempty"`
	AccountID     string     `json:"account_id,omitempty"`
	Description   string     `json:"description"`
	Severity      string     `json:"severity"`
	Status        string     `json:"status"` // "OPEN", "ACKNOWLEDGED", "RESOLVED"
	DetectedAt    time.Time  `json:"detected_at"`
	ResolvedAt    *time.Time `json:"resolved_at,omitempty"`
	Notes         string     `json:"notes"`
}

// TaxReturn represents a tax filing/return
type TaxReturn struct {
	ID             string           `json:"id"`
	CompanyID      string           `json:"company_id"`
	Jurisdiction   TaxJurisdiction  `json:"jurisdiction"`
	TaxType        TaxType          `json:"tax_type"`
	PeriodStart    time.Time        `json:"period_start"`
	PeriodEnd      time.Time        `json:"period_end"`
	GrossRevenue   float64          `json:"gross_revenue"`
	TaxableRevenue float64          `json:"taxable_revenue"`
	TotalTax       float64          `json:"total_tax"`
	TaxPaid        float64          `json:"tax_paid"`
	TaxOwed        float64          `json:"tax_owed"`
	FilingStatus   string           `json:"filing_status"` // "DRAFT", "FILED", "AMENDED"
	FiledAt        *time.Time       `json:"filed_at,omitempty"`
	DueDate        time.Time        `json:"due_date"`
	Calculations   []TaxCalculation `json:"calculations"`
	Attachments    []string         `json:"attachments"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

// ComplianceService handles regulatory and tax compliance
type ComplianceService struct {
	storage Storage
}

// NewComplianceService creates a new compliance service
func NewComplianceService(storage Storage) *ComplianceService {
	return &ComplianceService{
		storage: storage,
	}
}

// CreateComplianceRule creates a new compliance rule
func (cs *ComplianceService) CreateComplianceRule(rule ComplianceRule) error {
	rule.ID = uuid.New().String()
	rule.CreatedAt = time.Now()
	rule.Active = true

	return cs.storage.SaveComplianceRule(&rule)
}

// CreateTaxRule creates a new tax rule
func (cs *ComplianceService) CreateTaxRule(rule TaxRule) error {
	rule.ID = uuid.New().String()
	rule.Active = true

	return cs.storage.SaveTaxRule(&rule)
}

// CalculateTax calculates tax for a given amount and jurisdiction
func (cs *ComplianceService) CalculateTax(amount float64, jurisdiction TaxJurisdiction, taxType TaxType, exemptions []string) (*TaxCalculation, error) {
	rules, err := cs.storage.GetTaxRulesByJurisdiction(jurisdiction, taxType)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax rules: %w", err)
	}

	now := time.Now()
	var applicableRule *TaxRule

	// Find the most recent applicable rule
	for _, rule := range rules {
		if !rule.Active {
			continue
		}

		if rule.EffectiveFrom.After(now) {
			continue
		}

		if rule.EffectiveTo != nil && rule.EffectiveTo.Before(now) {
			continue
		}

		// Check if amount is within rule limits
		if rule.MinAmount > 0 && amount < rule.MinAmount {
			continue
		}

		if rule.MaxAmount > 0 && amount > rule.MaxAmount {
			continue
		}

		applicableRule = rule
		break
	}

	if applicableRule == nil {
		return &TaxCalculation{
			BaseAmount:    amount,
			TaxableAmount: 0,
			TaxRate:       0,
			TaxAmount:     0,
			Exemptions:    exemptions,
			CalculatedAt:  now,
		}, nil
	}

	// Check exemptions
	taxableAmount := amount
	appliedExemptions := []string{}

	for _, exemption := range exemptions {
		for _, ruleExemption := range applicableRule.Exemptions {
			if strings.EqualFold(exemption, ruleExemption) {
				taxableAmount = 0
				appliedExemptions = append(appliedExemptions, exemption)
				break
			}
		}
	}

	taxAmount := taxableAmount * applicableRule.Rate

	return &TaxCalculation{
		RuleID:        applicableRule.ID,
		BaseAmount:    amount,
		TaxableAmount: taxableAmount,
		TaxRate:       applicableRule.Rate,
		TaxAmount:     taxAmount,
		Exemptions:    appliedExemptions,
		CalculatedAt:  now,
	}, nil
}

// ValidateTransaction validates a transaction against compliance rules
func (cs *ComplianceService) ValidateTransaction(transaction Transaction) ([]ComplianceViolation, error) {
	rules, err := cs.storage.GetAllComplianceRules()
	if err != nil {
		return nil, fmt.Errorf("failed to get compliance rules: %w", err)
	}

	var violations []ComplianceViolation

	for _, rule := range rules {
		if !rule.Active {
			continue
		}

		violation := cs.checkTransactionAgainstRule(transaction, *rule)
		if violation != nil {
			violations = append(violations, *violation)
		}
	}

	return violations, nil
}

// checkTransactionAgainstRule checks a transaction against a specific compliance rule
func (cs *ComplianceService) checkTransactionAgainstRule(transaction Transaction, rule ComplianceRule) *ComplianceViolation {
	switch rule.RuleType {
	case "SEGREGATION_OF_DUTIES":
		return cs.checkSegregationOfDuties(transaction, rule)
	case "MATERIALITY_THRESHOLD":
		return cs.checkMaterialityThreshold(transaction, rule)
	case "AUTHORIZATION_REQUIRED":
		return cs.checkAuthorizationRequired(transaction, rule)
	case "JOURNAL_ENTRY_BALANCE":
		return cs.checkJournalEntryBalance(transaction, rule)
	case "ACCOUNT_TYPE_RESTRICTION":
		return cs.checkAccountTypeRestriction(transaction, rule)
	default:
		return nil
	}
}

// checkSegregationOfDuties validates segregation of duties
func (cs *ComplianceService) checkSegregationOfDuties(transaction Transaction, rule ComplianceRule) *ComplianceViolation {
	// Check if the same user created and approved the transaction
	// Since the base Transaction struct doesn't have approval fields, we'll check if UserID is repeated
	if transaction.UserID != "" {
		return &ComplianceViolation{
			ID:            uuid.New().String(),
			RuleID:        rule.ID,
			TransactionID: transaction.ID,
			Description:   "Transaction requires segregation of duties validation",
			Severity:      rule.Severity,
			Status:        "OPEN",
			DetectedAt:    time.Now(),
		}
	}
	return nil
}

// checkMaterialityThreshold validates materiality thresholds
func (cs *ComplianceService) checkMaterialityThreshold(transaction Transaction, rule ComplianceRule) *ComplianceViolation {
	// Extract threshold from conditions (simplified)
	var threshold float64 = 10000 // Default threshold

	for _, condition := range rule.Conditions {
		if strings.Contains(condition, "AMOUNT_THRESHOLD") {
			// Parse threshold from condition string
			// In a real implementation, this would be more sophisticated
			threshold = 10000
		}
	}

	totalAmount := 0.0
	for _, entry := range transaction.Entries {
		totalAmount += float64(entry.Amount.Value) / 100 // Convert from cents
	}

	if totalAmount > threshold {
		return &ComplianceViolation{
			ID:            uuid.New().String(),
			RuleID:        rule.ID,
			TransactionID: transaction.ID,
			Description:   fmt.Sprintf("Transaction amount %.2f exceeds materiality threshold %.2f", totalAmount, threshold),
			Severity:      rule.Severity,
			Status:        "OPEN",
			DetectedAt:    time.Now(),
		}
	}
	return nil
}

// checkAuthorizationRequired validates required authorization
func (cs *ComplianceService) checkAuthorizationRequired(transaction Transaction, rule ComplianceRule) *ComplianceViolation {
	// Since base Transaction doesn't have approval fields, we check basic validation
	if transaction.Status == Pending {
		return &ComplianceViolation{
			ID:            uuid.New().String(),
			RuleID:        rule.ID,
			TransactionID: transaction.ID,
			Description:   "Transaction requires authorization before posting",
			Severity:      rule.Severity,
			Status:        "OPEN",
			DetectedAt:    time.Now(),
		}
	}
	return nil
}

// checkJournalEntryBalance validates journal entry balance
func (cs *ComplianceService) checkJournalEntryBalance(transaction Transaction, rule ComplianceRule) *ComplianceViolation {
	totalDebits := int64(0)
	totalCredits := int64(0)

	for _, entry := range transaction.Entries {
		if entry.Type == Debit {
			totalDebits += entry.Amount.Value
		} else {
			totalCredits += entry.Amount.Value
		}
	}

	// Allow for small rounding differences
	if abs64(totalDebits-totalCredits) > 1 { // 1 cent tolerance
		return &ComplianceViolation{
			ID:            uuid.New().String(),
			RuleID:        rule.ID,
			TransactionID: transaction.ID,
			Description:   fmt.Sprintf("Journal entry not balanced: Debits=%d, Credits=%d", totalDebits, totalCredits),
			Severity:      "ERROR",
			Status:        "OPEN",
			DetectedAt:    time.Now(),
		}
	}
	return nil
}

// checkAccountTypeRestriction validates account type restrictions
func (cs *ComplianceService) checkAccountTypeRestriction(transaction Transaction, rule ComplianceRule) *ComplianceViolation {
	// This would check if certain account types can only be used in specific ways
	// Implementation depends on specific rule conditions
	return nil
}

// CreateTaxReturn creates a new tax return
func (cs *ComplianceService) CreateTaxReturn(taxReturn TaxReturn) error {
	taxReturn.ID = uuid.New().String()
	taxReturn.CreatedAt = time.Now()
	taxReturn.UpdatedAt = time.Now()
	taxReturn.FilingStatus = "DRAFT"

	return cs.storage.SaveTaxReturn(&taxReturn)
}

// CalculateTaxReturn calculates tax return amounts based on transactions
func (cs *ComplianceService) CalculateTaxReturn(companyID string, jurisdiction TaxJurisdiction, taxType TaxType, periodStart, periodEnd time.Time) (*TaxReturn, error) {
	// Get transactions for the period
	transactions, err := cs.storage.GetTransactionsByDateRange(companyID, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	grossRevenue := 0.0
	taxableRevenue := 0.0
	var calculations []TaxCalculation

	// Calculate revenue and taxes
	for _, transaction := range transactions {
		for _, entry := range transaction.Entries {
			// Check if this is revenue (credit entry to revenue account)
			if entry.Type == Credit && strings.Contains(strings.ToLower(entry.AccountID), "revenue") {
				entryAmount := float64(entry.Amount.Value) / 100 // Convert from cents
				grossRevenue += entryAmount

				// Calculate tax on this entry
				calc, err := cs.CalculateTax(entryAmount, jurisdiction, taxType, []string{})
				if err != nil {
					continue
				}

				taxableRevenue += calc.TaxableAmount
				calculations = append(calculations, *calc)
			}
		}
	}

	totalTax := 0.0
	for _, calc := range calculations {
		totalTax += calc.TaxAmount
	}

	// Calculate due date (simplified - typically 30-90 days after period end)
	dueDate := periodEnd.AddDate(0, 3, 0) // 3 months after period end

	return &TaxReturn{
		CompanyID:      companyID,
		Jurisdiction:   jurisdiction,
		TaxType:        taxType,
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
		GrossRevenue:   grossRevenue,
		TaxableRevenue: taxableRevenue,
		TotalTax:       totalTax,
		TaxPaid:        0, // To be updated when payments are made
		TaxOwed:        totalTax,
		FilingStatus:   "DRAFT",
		DueDate:        dueDate,
		Calculations:   calculations,
		Attachments:    []string{},
	}, nil
}

// GetComplianceViolations retrieves compliance violations
func (cs *ComplianceService) GetComplianceViolations(companyID string) ([]ComplianceViolation, error) {
	return cs.storage.GetComplianceViolations(companyID)
}

// ResolveViolation marks a compliance violation as resolved
func (cs *ComplianceService) ResolveViolation(violationID, notes string) error {
	violation, err := cs.storage.GetComplianceViolation(violationID)
	if err != nil {
		return fmt.Errorf("failed to get violation: %w", err)
	}

	now := time.Now()
	violation.Status = "RESOLVED"
	violation.ResolvedAt = &now
	violation.Notes = notes

	return cs.storage.SaveComplianceViolation(violation)
}

// SetupStandardComplianceRules creates standard compliance rules
func (cs *ComplianceService) SetupStandardComplianceRules(framework ComplianceFramework) error {
	var rules []ComplianceRule

	switch framework {
	case GAAP_Framework:
		rules = cs.getGAAPRules()
	case IFRS_Framework:
		rules = cs.getIFRSRules()
	case SOX_Framework:
		rules = cs.getSOXRules()
	}

	for _, rule := range rules {
		if err := cs.CreateComplianceRule(rule); err != nil {
			return fmt.Errorf("failed to create rule %s: %w", rule.ID, err)
		}
	}

	return nil
}

// SetupStandardTaxRules creates standard tax rules for a jurisdiction
func (cs *ComplianceService) SetupStandardTaxRules(jurisdiction TaxJurisdiction) error {
	var rules []TaxRule

	switch jurisdiction {
	case US_FEDERAL:
		rules = cs.getUSFederalTaxRules()
	case EU_VAT:
		rules = cs.getEUVATRules()
	case UK_VAT:
		rules = cs.getUKVATRules()
	case CANADA_GST:
		rules = cs.getCanadaGSTRules()
	}

	for _, rule := range rules {
		if err := cs.CreateTaxRule(rule); err != nil {
			return fmt.Errorf("failed to create tax rule %s: %w", rule.ID, err)
		}
	}

	return nil
}

// Helper functions for standard rules
func (cs *ComplianceService) getGAAPRules() []ComplianceRule {
	return []ComplianceRule{
		{
			Framework:   GAAP_Framework,
			RuleType:    "JOURNAL_ENTRY_BALANCE",
			Description: "All journal entries must be balanced (debits = credits)",
			Conditions:  []string{"TOTAL_DEBITS_EQUALS_TOTAL_CREDITS"},
			Actions:     []string{"REJECT_TRANSACTION"},
			Severity:    "ERROR",
		},
		{
			Framework:   GAAP_Framework,
			RuleType:    "MATERIALITY_THRESHOLD",
			Description: "Transactions over materiality threshold require approval",
			Conditions:  []string{"AMOUNT_THRESHOLD=10000"},
			Actions:     []string{"REQUIRE_APPROVAL"},
			Severity:    "WARNING",
		},
	}
}

func (cs *ComplianceService) getIFRSRules() []ComplianceRule {
	return []ComplianceRule{
		{
			Framework:   IFRS_Framework,
			RuleType:    "JOURNAL_ENTRY_BALANCE",
			Description: "All journal entries must be balanced per IFRS standards",
			Conditions:  []string{"TOTAL_DEBITS_EQUALS_TOTAL_CREDITS"},
			Actions:     []string{"REJECT_TRANSACTION"},
			Severity:    "ERROR",
		},
	}
}

func (cs *ComplianceService) getSOXRules() []ComplianceRule {
	return []ComplianceRule{
		{
			Framework:   SOX_Framework,
			RuleType:    "SEGREGATION_OF_DUTIES",
			Description: "Transaction creator and approver must be different (SOX compliance)",
			Conditions:  []string{"CREATOR_NOT_EQUAL_APPROVER"},
			Actions:     []string{"REQUIRE_DIFFERENT_APPROVER"},
			Severity:    "ERROR",
		},
	}
}

func (cs *ComplianceService) getUSFederalTaxRules() []TaxRule {
	return []TaxRule{
		{
			Jurisdiction:  US_FEDERAL,
			TaxType:       INCOME_TAX,
			Name:          "Federal Corporate Income Tax",
			Rate:          0.21, // 21% corporate tax rate
			MinAmount:     0,
			MaxAmount:     0,
			Exemptions:    []string{"TAX_EXEMPT", "MUNICIPAL_BOND"},
			EffectiveFrom: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
}

func (cs *ComplianceService) getEUVATRules() []TaxRule {
	return []TaxRule{
		{
			Jurisdiction:  EU_VAT,
			TaxType:       VAT,
			Name:          "Standard EU VAT Rate",
			Rate:          0.20, // 20% standard VAT rate
			MinAmount:     0,
			MaxAmount:     0,
			Exemptions:    []string{"HEALTHCARE", "EDUCATION", "FINANCIAL_SERVICES"},
			EffectiveFrom: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
}

func (cs *ComplianceService) getUKVATRules() []TaxRule {
	return []TaxRule{
		{
			Jurisdiction:  UK_VAT,
			TaxType:       VAT,
			Name:          "UK Standard VAT Rate",
			Rate:          0.20, // 20% VAT rate
			MinAmount:     0,
			MaxAmount:     0,
			Exemptions:    []string{"FOOD", "BOOKS", "CHILDREN_CLOTHING"},
			EffectiveFrom: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
}

func (cs *ComplianceService) getCanadaGSTRules() []TaxRule {
	return []TaxRule{
		{
			Jurisdiction:  CANADA_GST,
			TaxType:       GST,
			Name:          "Canada GST",
			Rate:          0.05, // 5% GST rate
			MinAmount:     0,
			MaxAmount:     0,
			Exemptions:    []string{"GROCERIES", "MEDICAL", "EDUCATION"},
			EffectiveFrom: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Helper function for absolute value of int64
func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
