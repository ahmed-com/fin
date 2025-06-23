package accounting

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ----------------------------------------------------------------------------
// AML (Anti-Money Laundering) Structures and Constants
// ----------------------------------------------------------------------------

// AMLFramework represents different AML regulatory frameworks
type AMLFramework string

const (
	BSA_Framework    AMLFramework = "BSA"    // Bank Secrecy Act (US)
	AMLD_Framework   AMLFramework = "AMLD"   // Anti-Money Laundering Directive (EU)
	FATF_Framework   AMLFramework = "FATF"   // Financial Action Task Force
	FINCEN_Framework AMLFramework = "FINCEN" // Financial Crimes Enforcement Network
	OFAC_Framework   AMLFramework = "OFAC"   // Office of Foreign Assets Control
)

// AMLRiskLevel represents the risk level of an entity or transaction
type AMLRiskLevel string

const (
	RiskLow      AMLRiskLevel = "LOW"
	RiskMedium   AMLRiskLevel = "MEDIUM"
	RiskHigh     AMLRiskLevel = "HIGH"
	RiskCritical AMLRiskLevel = "CRITICAL"
)

// AMLRuleType represents different types of AML rules
type AMLRuleType string

const (
	// Transaction-based rules
	RuleCTR           AMLRuleType = "CTR"            // Currency Transaction Report (>$10K)
	RuleSAR           AMLRuleType = "SAR"            // Suspicious Activity Report
	RuleStructuring   AMLRuleType = "STRUCTURING"    // Avoiding reporting thresholds
	RuleSmurfing      AMLRuleType = "SMURFING"       // Multiple small transactions
	RuleLayering      AMLRuleType = "LAYERING"       // Complex transaction chains
	RuleRapidMovement AMLRuleType = "RAPID_MOVEMENT" // Rapid in/out movement
	RuleRoundAmounts  AMLRuleType = "ROUND_AMOUNTS"  // Unusual round amounts

	// Customer-based rules
	RuleKYC       AMLRuleType = "KYC"       // Know Your Customer
	RuleCDD       AMLRuleType = "CDD"       // Customer Due Diligence
	RuleEDD       AMLRuleType = "EDD"       // Enhanced Due Diligence
	RulePEP       AMLRuleType = "PEP"       // Politically Exposed Person
	RuleSanctions AMLRuleType = "SANCTIONS" // Sanctions screening

	// Geographic rules
	RuleHighRiskJuris AMLRuleType = "HIGH_RISK_JURISDICTION"
	RuleNCCT          AMLRuleType = "NCCT" // Non-Cooperative Countries

	// Pattern-based rules
	RuleVelocity      AMLRuleType = "VELOCITY"      // Transaction velocity
	RuleFrequency     AMLRuleType = "FREQUENCY"     // Transaction frequency
	RuleConcentration AMLRuleType = "CONCENTRATION" // Account concentration

	// Additional common AML rules
	RuleUnusualTiming       AMLRuleType = "UNUSUAL_TIMING"        // Transactions at unusual times
	RuleCircularTransfers   AMLRuleType = "CIRCULAR_TRANSFERS"    // Circular money movements
	RuleCashIntensive       AMLRuleType = "CASH_INTENSIVE"        // High cash transaction volumes
	RuleHighRiskProducts    AMLRuleType = "HIGH_RISK_PRODUCTS"    // High-risk financial products
	RuleAccountDormancy     AMLRuleType = "ACCOUNT_DORMANCY"      // Dormant account activity
	RuleIdentityVerif       AMLRuleType = "IDENTITY_VERIFICATION" // Identity verification issues
	RuleSourceOfFunds       AMLRuleType = "SOURCE_OF_FUNDS"       // Unexplained source of funds
	RuleWireStripping       AMLRuleType = "WIRE_STRIPPING"        // Removal of wire transfer info
	RuleThirdPartyCheck     AMLRuleType = "THIRD_PARTY_CHECK"     // Third-party check deposits
	RuleNegativeMedia       AMLRuleType = "NEGATIVE_MEDIA"        // Adverse media screening
	RuleAccountTakeover     AMLRuleType = "ACCOUNT_TAKEOVER"      // Potential account takeover
	RuleTradeBasedML        AMLRuleType = "TRADE_BASED_ML"        // Trade-based money laundering
	RuleShellCompany        AMLRuleType = "SHELL_COMPANY"         // Shell company indicators
	RulePrepaidCards        AMLRuleType = "PREPAID_CARDS"         // Prepaid card abuse
	RuleCryptocurrency      AMLRuleType = "CRYPTOCURRENCY"        // Cryptocurrency transactions
	RuleJustUnderThreshold  AMLRuleType = "JUST_UNDER_THRESHOLD"  // Amounts just under thresholds
	RuleUnexpectedGeography AMLRuleType = "UNEXPECTED_GEOGRAPHY"  // Unexpected geographical activity
)

// AMLAlert represents an AML compliance alert
type AMLAlert struct {
	ID             string            `json:"id"`
	RuleType       AMLRuleType       `json:"rule_type"`
	Framework      AMLFramework      `json:"framework"`
	RiskLevel      AMLRiskLevel      `json:"risk_level"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	EntityID       string            `json:"entity_id"`   // Customer, account, or transaction ID
	EntityType     string            `json:"entity_type"` // "CUSTOMER", "ACCOUNT", "TRANSACTION"
	TransactionIDs []string          `json:"transaction_ids"`
	AccountIDs     []string          `json:"account_ids"`
	Amount         *Amount           `json:"amount,omitempty"`
	Currency       string            `json:"currency"`
	DetectedAt     time.Time         `json:"detected_at"`
	Status         string            `json:"status"` // "OPEN", "INVESTIGATING", "CLOSED", "ESCALATED"
	AssignedTo     string            `json:"assigned_to"`
	Investigation  *AMLInvestigation `json:"investigation,omitempty"`
	Evidence       []AMLEvidence     `json:"evidence"`
	Dispositions   []AMLDisposition  `json:"dispositions"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// AMLInvestigation represents an investigation into an AML alert
type AMLInvestigation struct {
	ID           string                `json:"id"`
	AlertID      string                `json:"alert_id"`
	Investigator string                `json:"investigator"`
	StartedAt    time.Time             `json:"started_at"`
	CompletedAt  *time.Time            `json:"completed_at,omitempty"`
	Status       string                `json:"status"`   // "ACTIVE", "COMPLETED", "ESCALATED"
	Priority     string                `json:"priority"` // "LOW", "MEDIUM", "HIGH", "URGENT"
	Findings     []string              `json:"findings"`
	Actions      []InvestigationAction `json:"actions"`
	Notes        []InvestigationNote   `json:"notes"`
}

// InvestigationAction represents an action taken during investigation
type InvestigationAction struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // "REVIEW", "REQUEST_INFO", "ESCALATE", "CLOSE"
	Description string    `json:"description"`
	TakenBy     string    `json:"taken_by"`
	TakenAt     time.Time `json:"taken_at"`
	Result      string    `json:"result"`
}

// InvestigationNote represents a note in an investigation
type InvestigationNote struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// AMLEvidence represents evidence supporting an AML alert
type AMLEvidence struct {
	Type        string      `json:"type"` // "TRANSACTION", "PATTERN", "EXTERNAL", "DOCUMENT"
	Description string      `json:"description"`
	Value       interface{} `json:"value"`
	Source      string      `json:"source"`
	Confidence  float64     `json:"confidence"` // 0.0 to 1.0
	CollectedAt time.Time   `json:"collected_at"`
}

// AMLDisposition represents the final disposition of an AML alert
type AMLDisposition struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // "NO_ACTION", "SAR_FILED", "ACCOUNT_CLOSED", "ESCALATED"
	Description string    `json:"description"`
	DecidedBy   string    `json:"decided_by"`
	DecidedAt   time.Time `json:"decided_at"`
	Rationale   string    `json:"rationale"`
	SARNumber   string    `json:"sar_number,omitempty"`
	ReportedTo  []string  `json:"reported_to,omitempty"`
}

// AMLRule represents an AML compliance rule
type AMLRule struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Type        AMLRuleType  `json:"type"`
	Framework   AMLFramework `json:"framework"`
	Description string       `json:"description"`
	Enabled     bool         `json:"enabled"`

	// Rule parameters
	Thresholds  map[string]interface{} `json:"thresholds"`
	TimeWindows map[string]int         `json:"time_windows"` // in hours
	Currencies  []string               `json:"currencies"`
	Countries   []string               `json:"countries"`

	// Risk scoring
	BaseScore    int     `json:"base_score"`    // 1-100
	RiskMultiple float64 `json:"risk_multiple"` // Multiplier for risk calculation

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AMLCustomer represents customer information for AML purposes
type AMLCustomer struct {
	ID             string       `json:"id"`
	CustomerID     string       `json:"customer_id"`
	Name           string       `json:"name"`
	Type           string       `json:"type"` // "INDIVIDUAL", "BUSINESS", "GOVERNMENT"
	RiskLevel      AMLRiskLevel `json:"risk_level"`
	Country        string       `json:"country"`
	IsPEP          bool         `json:"is_pep"` // Politically Exposed Person
	IsHighRisk     bool         `json:"is_high_risk"`
	SanctionsMatch bool         `json:"sanctions_match"`

	// Due diligence dates
	LastKYCDate    *time.Time `json:"last_kyc_date"`
	LastCDDDate    *time.Time `json:"last_cdd_date"`
	NextReviewDate *time.Time `json:"next_review_date"`

	// Business relationship
	OnboardingDate   time.Time `json:"onboarding_date"`
	ExpectedActivity string    `json:"expected_activity"`
	BusinessPurpose  string    `json:"business_purpose"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AMLTransaction represents transaction data enriched for AML analysis
type AMLTransaction struct {
	TransactionID  string    `json:"transaction_id"`
	Amount         *Amount   `json:"amount"`
	Currency       string    `json:"currency"`
	Date           time.Time `json:"date"`
	FromCustomerID string    `json:"from_customer_id"`
	ToCustomerID   string    `json:"to_customer_id"`
	FromCountry    string    `json:"from_country"`
	ToCountry      string    `json:"to_country"`
	Purpose        string    `json:"purpose"`
	Channel        string    `json:"channel"` // "WIRE", "ACH", "CASH", "CHECK"
	RiskScore      int       `json:"risk_score"`
	IsStructured   bool      `json:"is_structured"`
	IsSuspicious   bool      `json:"is_suspicious"`
	Flags          []string  `json:"flags"`
}

// ----------------------------------------------------------------------------
// AML Service
// ----------------------------------------------------------------------------

// AMLService handles anti-money laundering compliance
type AMLService struct {
	storage     *Storage
	compliance  *ComplianceService
	forensic    *ForensicService
	rules       map[string]*AMLRule
	customers   map[string]*AMLCustomer
	alertsCache map[string]*AMLAlert
}

// NewAMLService creates a new AML service
func NewAMLService(storage *Storage, compliance *ComplianceService, forensic *ForensicService) *AMLService {
	return &AMLService{
		storage:     storage,
		compliance:  compliance,
		forensic:    forensic,
		rules:       make(map[string]*AMLRule),
		customers:   make(map[string]*AMLCustomer),
		alertsCache: make(map[string]*AMLAlert),
	}
}

// ----------------------------------------------------------------------------
// Core AML Rule Implementation
// ----------------------------------------------------------------------------

// SetupStandardAMLRules creates standard AML rules based on framework
func (aml *AMLService) SetupStandardAMLRules(framework AMLFramework) error {
	switch framework {
	case BSA_Framework:
		return aml.setupBSARules()
	case AMLD_Framework:
		return aml.setupAMLDRules()
	case FATF_Framework:
		return aml.setupFATFRules()
	case FINCEN_Framework:
		return aml.setupFinCENRules()
	case OFAC_Framework:
		return aml.setupOFACRules()
	default:
		return fmt.Errorf("unsupported AML framework: %s", framework)
	}
}

// setupBSARules creates Bank Secrecy Act rules (US)
func (aml *AMLService) setupBSARules() error {
	rules := []*AMLRule{
		{
			ID:          generateUUID(),
			Name:        "CTR - Currency Transaction Report",
			Type:        RuleCTR,
			Framework:   BSA_Framework,
			Description: "Report cash transactions over $10,000",
			Enabled:     true,
			Thresholds: map[string]interface{}{
				"single_transaction": 1000000, // $10,000 in cents
				"daily_aggregate":    1000000,
			},
			TimeWindows: map[string]int{
				"aggregation_window": 24, // 24 hours
			},
			Currencies:   []string{"USD"},
			BaseScore:    75,
			RiskMultiple: 1.5,
		},
		{
			ID:          generateUUID(),
			Name:        "SAR - Suspicious Activity Threshold",
			Type:        RuleSAR,
			Framework:   BSA_Framework,
			Description: "Flag transactions for potential SAR filing",
			Enabled:     true,
			Thresholds: map[string]interface{}{
				"minimum_amount": 500000, // $5,000 in cents
				"risk_score":     70,
			},
			BaseScore:    80,
			RiskMultiple: 2.0,
		},
		{
			ID:          generateUUID(),
			Name:        "Structuring Detection",
			Type:        RuleStructuring,
			Framework:   BSA_Framework,
			Description: "Detect transactions designed to avoid CTR reporting",
			Enabled:     true,
			Thresholds: map[string]interface{}{
				"threshold_percentage": 0.9, // 90% of CTR threshold
				"transaction_count":    3,
				"amount_variation":     0.1, // 10% variation
			},
			TimeWindows: map[string]int{
				"detection_window": 72, // 3 days
			},
			BaseScore:    85,
			RiskMultiple: 2.5,
		},
	}

	for _, rule := range rules {
		rule.CreatedAt = time.Now()
		rule.UpdatedAt = time.Now()
		aml.rules[rule.ID] = rule

		// Save to storage
		if err := aml.storage.SaveAMLRule(rule); err != nil {
			return fmt.Errorf("failed to save AML rule %s: %w", rule.Name, err)
		}
	}

	return nil
}

// setupAMLDRules creates Anti-Money Laundering Directive rules (EU)
func (aml *AMLService) setupAMLDRules() error {
	rules := []*AMLRule{
		{
			ID:          generateUUID(),
			Name:        "EU Suspicious Transaction Threshold",
			Type:        RuleSAR,
			Framework:   AMLD_Framework,
			Description: "Flag suspicious transactions under AMLD",
			Enabled:     true,
			Thresholds: map[string]interface{}{
				"minimum_amount": 1500000, // €15,000 in cents
			},
			Currencies:   []string{"EUR", "USD", "GBP"},
			BaseScore:    70,
			RiskMultiple: 1.8,
		},
		{
			ID:           generateUUID(),
			Name:         "High-Risk Third Countries",
			Type:         RuleHighRiskJuris,
			Framework:    AMLD_Framework,
			Description:  "Enhanced monitoring for high-risk jurisdictions",
			Enabled:      true,
			Countries:    []string{"AF", "IR", "KP", "PK"}, // Afghanistan, Iran, North Korea, Pakistan
			BaseScore:    90,
			RiskMultiple: 3.0,
		},
	}

	for _, rule := range rules {
		rule.CreatedAt = time.Now()
		rule.UpdatedAt = time.Now()
		aml.rules[rule.ID] = rule

		if err := aml.storage.SaveAMLRule(rule); err != nil {
			return fmt.Errorf("failed to save AML rule %s: %w", rule.Name, err)
		}
	}

	return nil
}

// setupFATFRules creates FATF-based rules
func (aml *AMLService) setupFATFRules() error {
	rules := []*AMLRule{
		{
			ID:          generateUUID(),
			Name:        "FATF Velocity Monitoring",
			Type:        RuleVelocity,
			Framework:   FATF_Framework,
			Description: "Monitor transaction velocity patterns",
			Enabled:     true,
			Thresholds: map[string]interface{}{
				"transactions_per_hour": 10,
				"amount_per_hour":       10000000, // $100,000 in cents
			},
			TimeWindows: map[string]int{
				"velocity_window": 1, // 1 hour
			},
			BaseScore:    65,
			RiskMultiple: 1.5,
		},
		{
			ID:          generateUUID(),
			Name:        "Rapid Movement Pattern",
			Type:        RuleRapidMovement,
			Framework:   FATF_Framework,
			Description: "Detect rapid in-and-out movement of funds",
			Enabled:     true,
			Thresholds: map[string]interface{}{
				"rapid_movement_hours": 24,
				"minimum_amount":       1000000, // $10,000 in cents
			},
			BaseScore:    75,
			RiskMultiple: 2.0,
		},
	}

	for _, rule := range rules {
		rule.CreatedAt = time.Now()
		rule.UpdatedAt = time.Now()
		aml.rules[rule.ID] = rule

		if err := aml.storage.SaveAMLRule(rule); err != nil {
			return fmt.Errorf("failed to save AML rule %s: %w", rule.Name, err)
		}
	}

	return nil
}

// setupFinCENRules creates FinCEN-specific rules
func (aml *AMLService) setupFinCENRules() error {
	rules := []*AMLRule{
		{
			ID:          generateUUID(),
			Name:        "FinCEN Beneficial Ownership",
			Type:        RuleCDD,
			Framework:   FINCEN_Framework,
			Description: "Enhanced due diligence for beneficial ownership",
			Enabled:     true,
			Thresholds: map[string]interface{}{
				"ownership_percentage": 25.0, // 25% ownership threshold
			},
			BaseScore:    60,
			RiskMultiple: 1.3,
		},
	}

	for _, rule := range rules {
		rule.CreatedAt = time.Now()
		rule.UpdatedAt = time.Now()
		aml.rules[rule.ID] = rule

		if err := aml.storage.SaveAMLRule(rule); err != nil {
			return fmt.Errorf("failed to save AML rule %s: %w", rule.Name, err)
		}
	}

	return nil
}

// setupOFACRules creates OFAC sanctions rules
func (aml *AMLService) setupOFACRules() error {
	rules := []*AMLRule{
		{
			ID:           generateUUID(),
			Name:         "OFAC Sanctions Screening",
			Type:         RuleSanctions,
			Framework:    OFAC_Framework,
			Description:  "Screen against OFAC sanctions lists",
			Enabled:      true,
			BaseScore:    100,
			RiskMultiple: 5.0,
		},
	}

	for _, rule := range rules {
		rule.CreatedAt = time.Now()
		rule.UpdatedAt = time.Now()
		aml.rules[rule.ID] = rule

		if err := aml.storage.SaveAMLRule(rule); err != nil {
			return fmt.Errorf("failed to save AML rule %s: %w", rule.Name, err)
		}
	}

	return nil
}

// ----------------------------------------------------------------------------
// Additional Common AML Rule Implementations
// ----------------------------------------------------------------------------

// setupCommonAMLRules adds the most common AML detection rules
func (aml *AMLService) setupCommonAMLRules() error {
	rules := []*AMLRule{
		// 1. Cash Intensive Activity Detection
		{
			ID:          generateUUID(),
			Name:        "Cash Intensive Activity",
			Type:        RuleCashIntensive,
			Framework:   BSA_Framework,
			Description: "Detect customers with unusually high cash activity",
			Enabled:     true,
			Thresholds: map[string]interface{}{
				"cash_percentage": 80.0,    // 80% of transactions are cash
				"minimum_volume":  5000000, // $50,000 minimum volume in cents
				"time_window":     30,      // 30 days
			},
			BaseScore:    70,
			RiskMultiple: 1.8,
		},

		// 2. Just Under Threshold Detection
		{
			ID:          generateUUID(),
			Name:        "Just Under Threshold",
			Type:        RuleJustUnderThreshold,
			Framework:   BSA_Framework,
			Description: "Detect transactions just under reporting thresholds",
			Enabled:     true,
			Thresholds: map[string]interface{}{
				"threshold_amount": 1000000, // $10,000 in cents
				"tolerance_pct":    5.0,     // Within 5% of threshold
				"frequency_limit":  3,       // 3 or more occurrences
				"time_window":      7,       // Within 7 days
			},
			BaseScore:    85,
			RiskMultiple: 2.2,
		},

		// 3. Unusual Timing Detection
		{
			ID:          generateUUID(),
			Name:        "Unusual Timing",
			Type:        RuleUnusualTiming,
			Framework:   BSA_Framework,
			Description: "Detect transactions at unusual times (nights, weekends, holidays)",
			Enabled:     true,
			Thresholds: map[string]interface{}{
				"night_start_hour": 22,     // 10 PM
				"night_end_hour":   6,      // 6 AM
				"minimum_amount":   100000, // $1,000 minimum
			},
			BaseScore:    40,
			RiskMultiple: 1.2,
		},

		// 4. Account Dormancy Reactivation
		{
			ID:          generateUUID(),
			Name:        "Dormant Account Reactivation",
			Type:        RuleAccountDormancy,
			Framework:   BSA_Framework,
			Description: "Detect sudden activity in previously dormant accounts",
			Enabled:     true,
			Thresholds: map[string]interface{}{
				"dormancy_period":  90,     // 90 days of inactivity
				"reactivation_amt": 500000, // $5,000 or more
			},
			BaseScore:    60,
			RiskMultiple: 1.5,
		},

		// 5. Wire Stripping Detection
		{
			ID:          generateUUID(),
			Name:        "Wire Stripping",
			Type:        RuleWireStripping,
			Framework:   BSA_Framework,
			Description: "Detect removal of wire transfer information",
			Enabled:     true,
			Thresholds: map[string]interface{}{
				"minimum_amount": 300000, // $3,000
			},
			BaseScore:    80,
			RiskMultiple: 2.0,
		},

		// 6. High-Risk Geography
		{
			ID:           generateUUID(),
			Name:         "Unexpected Geography",
			Type:         RuleUnexpectedGeography,
			Framework:    FATF_Framework,
			Description:  "Detect transactions from unexpected geographical locations",
			Enabled:      true,
			Countries:    []string{"AF", "MM", "KP", "IR", "SY"}, // High-risk countries
			BaseScore:    90,
			RiskMultiple: 3.0,
		},

		// 7. Cryptocurrency Transactions
		{
			ID:          generateUUID(),
			Name:        "Cryptocurrency Activity",
			Type:        RuleCryptocurrency,
			Framework:   BSA_Framework,
			Description: "Monitor cryptocurrency-related transactions",
			Enabled:     true,
			Thresholds: map[string]interface{}{
				"minimum_amount": 100000, // $1,000
			},
			BaseScore:    50,
			RiskMultiple: 1.4,
		},

		// 8. Shell Company Indicators
		{
			ID:          generateUUID(),
			Name:        "Shell Company Indicators",
			Type:        RuleShellCompany,
			Framework:   FATF_Framework,
			Description: "Detect potential shell company characteristics",
			Enabled:     true,
			Thresholds: map[string]interface{}{
				"transaction_complexity": 5, // Complex transaction patterns
				"minimal_operations":     1, // Minimal legitimate operations
			},
			BaseScore:    80,
			RiskMultiple: 2.5,
		},

		// 9. Trade-Based Money Laundering
		{
			ID:          generateUUID(),
			Name:        "Trade-Based Money Laundering",
			Type:        RuleTradeBasedML,
			Framework:   FATF_Framework,
			Description: "Detect trade-based money laundering patterns",
			Enabled:     true,
			Thresholds: map[string]interface{}{
				"price_variance_pct": 20.0,    // 20% price variance from market
				"volume_threshold":   1000000, // $10,000
			},
			BaseScore:    75,
			RiskMultiple: 2.0,
		},

		// 10. Third-Party Check Deposits
		{
			ID:          generateUUID(),
			Name:        "Third-Party Check Deposits",
			Type:        RuleThirdPartyCheck,
			Framework:   BSA_Framework,
			Description: "Monitor third-party check deposits",
			Enabled:     true,
			Thresholds: map[string]interface{}{
				"frequency_limit": 5,     // 5 or more per month
				"minimum_amount":  50000, // $500
			},
			BaseScore:    45,
			RiskMultiple: 1.3,
		},
	}

	for _, rule := range rules {
		rule.CreatedAt = time.Now()
		rule.UpdatedAt = time.Now()
		aml.rules[rule.ID] = rule

		if err := aml.storage.SaveAMLRule(rule); err != nil {
			return fmt.Errorf("failed to save AML rule %s: %w", rule.Name, err)
		}
	}

	return nil
}

// SetupAllStandardAMLRules sets up all common AML rules across frameworks
func (aml *AMLService) SetupAllStandardAMLRules() error {
	frameworks := []AMLFramework{
		BSA_Framework,
		FATF_Framework,
		AMLD_Framework,
		FINCEN_Framework,
		OFAC_Framework,
	}

	for _, framework := range frameworks {
		if err := aml.SetupStandardAMLRules(framework); err != nil {
			return fmt.Errorf("failed to setup %s rules: %w", framework, err)
		}
	}

	// Add common rules not specific to any framework
	if err := aml.setupCommonAMLRules(); err != nil {
		return fmt.Errorf("failed to setup common AML rules: %w", err)
	}

	return nil
}

// ----------------------------------------------------------------------------
// Transaction Monitoring
// ----------------------------------------------------------------------------

// MonitorTransaction analyzes a transaction against AML rules
func (aml *AMLService) MonitorTransaction(txn *Transaction, customerInfo map[string]*AMLCustomer) ([]*AMLAlert, error) {
	var alerts []*AMLAlert

	// Convert transaction to AML format
	amlTxn := aml.convertToAMLTransaction(txn, customerInfo)

	// Run traditional rule-based checks
	for _, rule := range aml.rules {
		if !rule.Enabled {
			continue
		}

		alert := aml.evaluateRule(rule, amlTxn, customerInfo)
		if alert != nil {
			alerts = append(alerts, alert)

			// Save alert
			if err := aml.storage.SaveAMLAlert(alert); err != nil {
				return nil, fmt.Errorf("failed to save AML alert: %w", err)
			}

			// Cache for quick access
			aml.alertsCache[alert.ID] = alert
		}
	}

	// Run advanced AML checks
	advancedChecks := []func(*Transaction) (*AMLAlert, error){
		aml.CheckJustUnderThreshold,
		aml.CheckUnusualTiming,
		aml.CheckDormantAccountReactivation,
	}

	for _, check := range advancedChecks {
		alert, err := check(txn)
		if err != nil {
			// Log error but continue with other checks
			continue
		}
		if alert != nil {
			alerts = append(alerts, alert)

			// Save alert
			if err := aml.storage.SaveAMLAlert(alert); err != nil {
				return nil, fmt.Errorf("failed to save AML alert: %w", err)
			}

			// Cache for quick access
			aml.alertsCache[alert.ID] = alert
		}
	}

	// Run customer-specific checks if customer info is available
	for _, customer := range customerInfo {
		if alert, err := aml.CheckHighRiskGeography(txn, customer); err == nil && alert != nil {
			alerts = append(alerts, alert)

			if err := aml.storage.SaveAMLAlert(alert); err != nil {
				return nil, fmt.Errorf("failed to save AML alert: %w", err)
			}

			aml.alertsCache[alert.ID] = alert
		}

		// Check cash intensive activity (periodic check)
		if alert, err := aml.CheckCashIntensiveActivity(customer.CustomerID, 30); err == nil && alert != nil {
			alerts = append(alerts, alert)

			if err := aml.storage.SaveAMLAlert(alert); err != nil {
				return nil, fmt.Errorf("failed to save AML alert: %w", err)
			}

			aml.alertsCache[alert.ID] = alert
		}
	}

	return alerts, nil
}

// convertToAMLTransaction converts a regular transaction to AML format
func (aml *AMLService) convertToAMLTransaction(txn *Transaction, customerInfo map[string]*AMLCustomer) *AMLTransaction {
	amlTxn := &AMLTransaction{
		TransactionID: txn.ID,
		Date:          txn.ValidTime,
		Currency:      "USD", // Default currency
		Flags:         []string{},
	}

	// Calculate total amount and determine flow
	var totalAmount int64
	var currency Currency = "USD"
	for _, entry := range txn.Entries {
		totalAmount += entry.Amount.Value
		if entry.Amount.Currency != "" {
			currency = entry.Amount.Currency
		}
	}

	amlTxn.Amount = &Amount{
		Value:    totalAmount / 2, // Divide by 2 since double-entry means each amount is counted twice
		Currency: currency,
	}

	amlTxn.Currency = string(currency)

	// For demo purposes, derive customer information from transaction description and reference
	// In a real system, this would come from proper metadata or foreign keys
	if txn.SourceRef != "" {
		// Simple mapping based on reference patterns
		switch {
		case txn.SourceRef == "SAL-001":
			amlTxn.FromCustomerID = "customer_001"
			amlTxn.Purpose = "Monthly salary payment"
			amlTxn.Channel = "ACH"
		case txn.SourceRef == "CASH-001":
			amlTxn.FromCustomerID = "customer_002"
			amlTxn.Purpose = "Business cash deposit"
			amlTxn.Channel = "CASH"
		case txn.SourceRef == "WIRE-001":
			amlTxn.FromCustomerID = "customer_002"
			amlTxn.ToCustomerID = "customer_001"
			amlTxn.Purpose = "" // Vague purpose
			amlTxn.Channel = "WIRE"
		case txn.SourceRef == "INTL-001":
			amlTxn.FromCustomerID = "customer_003" // Sanctioned customer
			amlTxn.ToCustomerID = "customer_001"
			amlTxn.Purpose = "Business payment"
			amlTxn.Channel = "WIRE"
		}
	}

	// Set countries based on customer info
	if customer, exists := customerInfo[amlTxn.FromCustomerID]; exists {
		amlTxn.FromCountry = customer.Country
	}
	if customer, exists := customerInfo[amlTxn.ToCustomerID]; exists {
		amlTxn.ToCountry = customer.Country
	}

	return amlTxn
}

// evaluateRule evaluates a transaction against a specific AML rule
func (aml *AMLService) evaluateRule(rule *AMLRule, txn *AMLTransaction, customerInfo map[string]*AMLCustomer) *AMLAlert {
	switch rule.Type {
	case RuleCTR:
		return aml.evaluateCTRRule(rule, txn)
	case RuleSAR:
		return aml.evaluateSARRule(rule, txn)
	case RuleStructuring:
		return aml.evaluateStructuringRule(rule, txn)
	case RuleVelocity:
		return aml.evaluateVelocityRule(rule, txn)
	case RuleRapidMovement:
		return aml.evaluateRapidMovementRule(rule, txn)
	case RuleHighRiskJuris:
		return aml.evaluateHighRiskJurisdictionRule(rule, txn)
	case RuleSanctions:
		return aml.evaluateSanctionsRule(rule, txn, customerInfo)
	default:
		return nil
	}
}

// evaluateCTRRule evaluates Currency Transaction Report rule
func (aml *AMLService) evaluateCTRRule(rule *AMLRule, txn *AMLTransaction) *AMLAlert {
	threshold, ok := rule.Thresholds["single_transaction"].(int)
	if !ok {
		return nil
	}

	if txn.Amount.Value >= int64(threshold) && txn.Channel == "CASH" {
		return &AMLAlert{
			ID:             generateUUID(),
			RuleType:       rule.Type,
			Framework:      rule.Framework,
			RiskLevel:      RiskHigh,
			Title:          "Currency Transaction Report Required",
			Description:    fmt.Sprintf("Cash transaction of %s %d exceeds CTR threshold", txn.Currency, txn.Amount.Value),
			EntityID:       txn.TransactionID,
			EntityType:     "TRANSACTION",
			TransactionIDs: []string{txn.TransactionID},
			Amount:         txn.Amount,
			Currency:       txn.Currency,
			DetectedAt:     time.Now(),
			Status:         "OPEN",
			Evidence: []AMLEvidence{
				{
					Type:        "TRANSACTION",
					Description: "High-value cash transaction",
					Value:       txn.Amount.Value,
					Source:      "TRANSACTION_MONITOR",
					Confidence:  0.95,
					CollectedAt: time.Now(),
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	return nil
}

// evaluateSARRule evaluates Suspicious Activity Report rule
func (aml *AMLService) evaluateSARRule(rule *AMLRule, txn *AMLTransaction) *AMLAlert {
	minAmount, ok := rule.Thresholds["minimum_amount"].(int)
	if !ok {
		return nil
	}

	if txn.Amount.Value >= int64(minAmount) {
		// Calculate suspicion score based on various factors
		suspicionScore := aml.calculateSuspicionScore(txn)

		if suspicionScore >= 70 { // Threshold for SAR consideration
			riskLevel := RiskMedium
			if suspicionScore >= 90 {
				riskLevel = RiskCritical
			} else if suspicionScore >= 80 {
				riskLevel = RiskHigh
			}

			return &AMLAlert{
				ID:             generateUUID(),
				RuleType:       rule.Type,
				Framework:      rule.Framework,
				RiskLevel:      riskLevel,
				Title:          "Potential Suspicious Activity",
				Description:    fmt.Sprintf("Transaction shows suspicious patterns (score: %d)", suspicionScore),
				EntityID:       txn.TransactionID,
				EntityType:     "TRANSACTION",
				TransactionIDs: []string{txn.TransactionID},
				Amount:         txn.Amount,
				Currency:       txn.Currency,
				DetectedAt:     time.Now(),
				Status:         "OPEN",
				Evidence: []AMLEvidence{
					{
						Type:        "PATTERN",
						Description: "Suspicious activity pattern detected",
						Value:       suspicionScore,
						Source:      "PATTERN_ANALYZER",
						Confidence:  float64(suspicionScore) / 100.0,
						CollectedAt: time.Now(),
					},
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}
	}

	return nil
}

// evaluateStructuringRule evaluates structuring detection rule
func (aml *AMLService) evaluateStructuringRule(rule *AMLRule, txn *AMLTransaction) *AMLAlert {
	// This would require historical transaction analysis
	// For now, implement basic round amount detection
	if aml.isRoundAmount(txn.Amount.Value) {
		return &AMLAlert{
			ID:             generateUUID(),
			RuleType:       rule.Type,
			Framework:      rule.Framework,
			RiskLevel:      RiskMedium,
			Title:          "Potential Structuring Activity",
			Description:    "Transaction uses suspiciously round amounts",
			EntityID:       txn.TransactionID,
			EntityType:     "TRANSACTION",
			TransactionIDs: []string{txn.TransactionID},
			Amount:         txn.Amount,
			Currency:       txn.Currency,
			DetectedAt:     time.Now(),
			Status:         "OPEN",
			Evidence: []AMLEvidence{
				{
					Type:        "TRANSACTION",
					Description: "Round amount transaction",
					Value:       txn.Amount.Value,
					Source:      "AMOUNT_ANALYZER",
					Confidence:  0.7,
					CollectedAt: time.Now(),
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	return nil
}

// evaluateVelocityRule evaluates transaction velocity rule
func (aml *AMLService) evaluateVelocityRule(rule *AMLRule, txn *AMLTransaction) *AMLAlert {
	// This would require querying recent transactions
	// Implementation would check transaction frequency and amounts
	// For demonstration, return nil (would need historical data)
	return nil
}

// evaluateRapidMovementRule evaluates rapid movement rule
func (aml *AMLService) evaluateRapidMovementRule(rule *AMLRule, txn *AMLTransaction) *AMLAlert {
	// This would require tracking fund flows
	// Implementation would detect rapid in-and-out patterns
	// For demonstration, return nil (would need flow analysis)
	return nil
}

// evaluateHighRiskJurisdictionRule evaluates high-risk jurisdiction rule
func (aml *AMLService) evaluateHighRiskJurisdictionRule(rule *AMLRule, txn *AMLTransaction) *AMLAlert {
	for _, country := range rule.Countries {
		if txn.FromCountry == country || txn.ToCountry == country {
			return &AMLAlert{
				ID:             generateUUID(),
				RuleType:       rule.Type,
				Framework:      rule.Framework,
				RiskLevel:      RiskHigh,
				Title:          "High-Risk Jurisdiction Transaction",
				Description:    fmt.Sprintf("Transaction involves high-risk country: %s", country),
				EntityID:       txn.TransactionID,
				EntityType:     "TRANSACTION",
				TransactionIDs: []string{txn.TransactionID},
				Amount:         txn.Amount,
				Currency:       txn.Currency,
				DetectedAt:     time.Now(),
				Status:         "OPEN",
				Evidence: []AMLEvidence{
					{
						Type:        "GEOGRAPHIC",
						Description: "High-risk jurisdiction involvement",
						Value:       country,
						Source:      "JURISDICTION_CHECKER",
						Confidence:  0.9,
						CollectedAt: time.Now(),
					},
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}
	}

	return nil
}

// evaluateSanctionsRule evaluates sanctions screening rule
func (aml *AMLService) evaluateSanctionsRule(rule *AMLRule, txn *AMLTransaction, customerInfo map[string]*AMLCustomer) *AMLAlert {
	// Check if any involved customer has sanctions match
	for customerID, customer := range customerInfo {
		if (customerID == txn.FromCustomerID || customerID == txn.ToCustomerID) && customer.SanctionsMatch {
			return &AMLAlert{
				ID:             generateUUID(),
				RuleType:       rule.Type,
				Framework:      rule.Framework,
				RiskLevel:      RiskCritical,
				Title:          "Sanctions Match Detected",
				Description:    fmt.Sprintf("Customer %s matches sanctions list", customer.Name),
				EntityID:       customerID,
				EntityType:     "CUSTOMER",
				TransactionIDs: []string{txn.TransactionID},
				Amount:         txn.Amount,
				Currency:       txn.Currency,
				DetectedAt:     time.Now(),
				Status:         "OPEN",
				Evidence: []AMLEvidence{
					{
						Type:        "SANCTIONS",
						Description: "Customer sanctions list match",
						Value:       customer.Name,
						Source:      "SANCTIONS_SCREENER",
						Confidence:  1.0,
						CollectedAt: time.Now(),
					},
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}
	}

	return nil
}

// ----------------------------------------------------------------------------
// Risk Assessment and Scoring
// ----------------------------------------------------------------------------

// calculateSuspicionScore calculates a suspicion score for a transaction
func (aml *AMLService) calculateSuspicionScore(txn *AMLTransaction) int {
	score := 0

	// Round amounts (higher suspicion)
	if aml.isRoundAmount(txn.Amount.Value) {
		score += 20
	}

	// High-value transactions
	if txn.Amount.Value >= 1000000 { // $10,000+
		score += 15
	} else if txn.Amount.Value >= 500000 { // $5,000+
		score += 10
	}

	// Cash transactions
	if txn.Channel == "CASH" {
		score += 25
	}

	// Cross-border transactions
	if txn.FromCountry != "" && txn.ToCountry != "" && txn.FromCountry != txn.ToCountry {
		score += 10
	}

	// Unusual timing (weekends, holidays, after hours)
	if aml.isUnusualTiming(txn.Date) {
		score += 15
	}

	// Vague purpose
	if txn.Purpose == "" || len(strings.TrimSpace(txn.Purpose)) < 5 {
		score += 10
	}

	return score
}

// isRoundAmount checks if an amount is suspiciously round
func (aml *AMLService) isRoundAmount(amount int64) bool {
	// Check for round thousands, ten-thousands, etc.
	return amount%100000 == 0 || // Multiples of $1,000
		amount%1000000 == 0 || // Multiples of $10,000
		amount%10000000 == 0 // Multiples of $100,000
}

// isUnusualTiming checks if transaction timing is unusual
func (aml *AMLService) isUnusualTiming(timestamp time.Time) bool {
	// Weekend transactions
	weekday := timestamp.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return true
	}

	// After hours (before 9 AM or after 5 PM)
	hour := timestamp.Hour()
	if hour < 9 || hour > 17 {
		return true
	}

	return false
}

// ----------------------------------------------------------------------------
// Alert Management
// ----------------------------------------------------------------------------

// GetAMLAlerts retrieves AML alerts with filtering
func (aml *AMLService) GetAMLAlerts(status string, riskLevel AMLRiskLevel, limit int) ([]*AMLAlert, error) {
	alerts, err := aml.storage.GetAMLAlerts()
	if err != nil {
		return nil, err
	}

	// Filter alerts
	var filtered []*AMLAlert
	for _, alert := range alerts {
		if status != "" && alert.Status != status {
			continue
		}
		if riskLevel != "" && alert.RiskLevel != riskLevel {
			continue
		}
		filtered = append(filtered, alert)
	}

	// Sort by detection date (newest first)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].DetectedAt.After(filtered[j].DetectedAt)
	})

	// Apply limit
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}

	return filtered, nil
}

// UpdateAlertStatus updates the status of an AML alert
func (aml *AMLService) UpdateAlertStatus(alertID, status, userID string) error {
	alert, err := aml.storage.GetAMLAlert(alertID)
	if err != nil {
		return err
	}

	alert.Status = status
	alert.UpdatedAt = time.Now()

	// Add disposition if closing
	if status == "CLOSED" {
		disposition := AMLDisposition{
			ID:          generateUUID(),
			Type:        "NO_ACTION",
			Description: "Alert reviewed and closed",
			DecidedBy:   userID,
			DecidedAt:   time.Now(),
			Rationale:   "No suspicious activity found upon review",
		}
		alert.Dispositions = append(alert.Dispositions, disposition)
	}

	return aml.storage.SaveAMLAlert(alert)
}

// CreateInvestigation creates a new investigation for an alert
func (aml *AMLService) CreateInvestigation(alertID, investigatorID string) (*AMLInvestigation, error) {
	investigation := &AMLInvestigation{
		ID:           generateUUID(),
		AlertID:      alertID,
		Investigator: investigatorID,
		StartedAt:    time.Now(),
		Status:       "ACTIVE",
		Priority:     "MEDIUM",
		Findings:     []string{},
		Actions:      []InvestigationAction{},
		Notes:        []InvestigationNote{},
	}

	// Update alert
	alert, err := aml.storage.GetAMLAlert(alertID)
	if err != nil {
		return nil, err
	}

	alert.Status = "INVESTIGATING"
	alert.AssignedTo = investigatorID
	alert.Investigation = investigation
	alert.UpdatedAt = time.Now()

	if err := aml.storage.SaveAMLAlert(alert); err != nil {
		return nil, err
	}

	return investigation, nil
}

// AddInvestigationNote adds a note to an investigation
func (aml *AMLService) AddInvestigationNote(alertID, content, userID string) error {
	alert, err := aml.storage.GetAMLAlert(alertID)
	if err != nil {
		return err
	}

	if alert.Investigation == nil {
		return fmt.Errorf("no active investigation for alert %s", alertID)
	}

	note := InvestigationNote{
		ID:        generateUUID(),
		Content:   content,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}

	alert.Investigation.Notes = append(alert.Investigation.Notes, note)
	alert.UpdatedAt = time.Now()

	return aml.storage.SaveAMLAlert(alert)
}

// ----------------------------------------------------------------------------
// Customer Management
// ----------------------------------------------------------------------------

// RegisterCustomer registers a customer for AML monitoring
func (aml *AMLService) RegisterCustomer(customer *AMLCustomer) error {
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	aml.customers[customer.ID] = customer
	return aml.storage.SaveAMLCustomer(customer)
}

// UpdateCustomerRisk updates a customer's risk level
func (aml *AMLService) UpdateCustomerRisk(customerID string, riskLevel AMLRiskLevel, reason string) error {
	customer, err := aml.storage.GetAMLCustomer(customerID)
	if err != nil {
		return err
	}

	customer.RiskLevel = riskLevel
	customer.UpdatedAt = time.Now()

	return aml.storage.SaveAMLCustomer(customer)
}

// PerformKYC performs Know Your Customer check
func (aml *AMLService) PerformKYC(customerID string, userID string) error {
	customer, err := aml.storage.GetAMLCustomer(customerID)
	if err != nil {
		return err
	}

	now := time.Now()
	customer.LastKYCDate = &now

	// Set next review date (annually for low risk, semi-annually for high risk)
	nextReview := now.AddDate(1, 0, 0) // 1 year
	if customer.RiskLevel == RiskHigh || customer.RiskLevel == RiskCritical {
		nextReview = now.AddDate(0, 6, 0) // 6 months
	}
	customer.NextReviewDate = &nextReview

	customer.UpdatedAt = time.Now()

	return aml.storage.SaveAMLCustomer(customer)
}

// GetCustomersForReview gets customers that need KYC review
func (aml *AMLService) GetCustomersForReview() ([]*AMLCustomer, error) {
	customers, err := aml.storage.GetAllAMLCustomers()
	if err != nil {
		return nil, err
	}

	var needReview []*AMLCustomer
	now := time.Now()

	for _, customer := range customers {
		if customer.NextReviewDate != nil && customer.NextReviewDate.Before(now) {
			needReview = append(needReview, customer)
		}
	}

	return needReview, nil
}

// ----------------------------------------------------------------------------
// Reporting and Analytics
// ----------------------------------------------------------------------------

// GenerateAMLReport generates various AML reports
func (aml *AMLService) GenerateAMLReport(reportType string, startDate, endDate time.Time) (interface{}, error) {
	switch reportType {
	case "ALERTS_SUMMARY":
		return aml.generateAlertsSummaryReport(startDate, endDate)
	case "RISK_ASSESSMENT":
		return aml.generateRiskAssessmentReport()
	case "CTR_REPORT":
		return aml.generateCTRReport(startDate, endDate)
	case "SAR_CANDIDATES":
		return aml.generateSARCandidatesReport(startDate, endDate)
	default:
		return nil, fmt.Errorf("unsupported report type: %s", reportType)
	}
}

// generateAlertsSummaryReport generates summary of AML alerts
func (aml *AMLService) generateAlertsSummaryReport(startDate, endDate time.Time) (map[string]interface{}, error) {
	alerts, err := aml.storage.GetAMLAlerts()
	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"period":        fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		"total_alerts":  0,
		"by_risk_level": make(map[string]int),
		"by_rule_type":  make(map[string]int),
		"by_status":     make(map[string]int),
	}

	riskCounts := make(map[string]int)
	ruleCounts := make(map[string]int)
	statusCounts := make(map[string]int)

	for _, alert := range alerts {
		if alert.DetectedAt.Before(startDate) || alert.DetectedAt.After(endDate) {
			continue
		}

		summary["total_alerts"] = summary["total_alerts"].(int) + 1
		riskCounts[string(alert.RiskLevel)]++
		ruleCounts[string(alert.RuleType)]++
		statusCounts[alert.Status]++
	}

	summary["by_risk_level"] = riskCounts
	summary["by_rule_type"] = ruleCounts
	summary["by_status"] = statusCounts

	return summary, nil
}

// generateRiskAssessmentReport generates risk assessment report
func (aml *AMLService) generateRiskAssessmentReport() (map[string]interface{}, error) {
	customers, err := aml.storage.GetAllAMLCustomers()
	if err != nil {
		return nil, err
	}

	riskDistribution := make(map[string]int)
	pepCount := 0
	sanctionsCount := 0
	totalCustomers := len(customers)

	for _, customer := range customers {
		riskDistribution[string(customer.RiskLevel)]++
		if customer.IsPEP {
			pepCount++
		}
		if customer.SanctionsMatch {
			sanctionsCount++
		}
	}

	return map[string]interface{}{
		"total_customers":   totalCustomers,
		"risk_distribution": riskDistribution,
		"pep_count":         pepCount,
		"sanctions_matches": sanctionsCount,
		"generated_at":      time.Now(),
	}, nil
}

// generateCTRReport generates CTR report
func (aml *AMLService) generateCTRReport(startDate, endDate time.Time) ([]map[string]interface{}, error) {
	alerts, err := aml.storage.GetAMLAlerts()
	if err != nil {
		return nil, err
	}

	var ctrAlerts []map[string]interface{}

	for _, alert := range alerts {
		if alert.RuleType == RuleCTR &&
			alert.DetectedAt.After(startDate) &&
			alert.DetectedAt.Before(endDate) {

			ctrAlerts = append(ctrAlerts, map[string]interface{}{
				"alert_id":       alert.ID,
				"transaction_id": alert.EntityID,
				"amount":         alert.Amount,
				"currency":       alert.Currency,
				"detected_at":    alert.DetectedAt,
				"status":         alert.Status,
			})
		}
	}

	return ctrAlerts, nil
}

// generateSARCandidatesReport generates SAR candidates report
func (aml *AMLService) generateSARCandidatesReport(startDate, endDate time.Time) ([]map[string]interface{}, error) {
	alerts, err := aml.storage.GetAMLAlerts()
	if err != nil {
		return nil, err
	}

	var sarCandidates []map[string]interface{}

	for _, alert := range alerts {
		if (alert.RuleType == RuleSAR || alert.RiskLevel == RiskHigh || alert.RiskLevel == RiskCritical) &&
			alert.DetectedAt.After(startDate) &&
			alert.DetectedAt.Before(endDate) &&
			alert.Status == "OPEN" {

			sarCandidates = append(sarCandidates, map[string]interface{}{
				"alert_id":    alert.ID,
				"rule_type":   alert.RuleType,
				"risk_level":  alert.RiskLevel,
				"description": alert.Description,
				"amount":      alert.Amount,
				"detected_at": alert.DetectedAt,
				"entity_id":   alert.EntityID,
				"entity_type": alert.EntityType,
			})
		}
	}

	return sarCandidates, nil
}

// ----------------------------------------------------------------------------
// AML Monitoring Dashboard and Reporting
// ----------------------------------------------------------------------------

// AMLDashboard represents comprehensive AML monitoring data
type AMLDashboard struct {
	PeriodStart        time.Time             `json:"period_start"`
	PeriodEnd          time.Time             `json:"period_end"`
	TotalAlerts        int                   `json:"total_alerts"`
	AlertsByRiskLevel  map[AMLRiskLevel]int  `json:"alerts_by_risk_level"`
	AlertsByType       map[AMLRuleType]int   `json:"alerts_by_type"`
	TopRiskyCustomers  []CustomerRiskSummary `json:"top_risky_customers"`
	ComplianceMetrics  AMLComplianceMetrics  `json:"compliance_metrics"`
	TrendAnalysis      AMLTrendAnalysis      `json:"trend_analysis"`
	RecommendedActions []AMLRecommendation   `json:"recommended_actions"`
}

type CustomerRiskSummary struct {
	CustomerID   string    `json:"customer_id"`
	CustomerName string    `json:"customer_name"`
	RiskScore    int       `json:"risk_score"`
	AlertCount   int       `json:"alert_count"`
	TotalVolume  int64     `json:"total_volume"`
	LastActivity time.Time `json:"last_activity"`
	RiskFactors  []string  `json:"risk_factors"`
}

type AMLComplianceMetrics struct {
	CTRFilingRate         float64 `json:"ctr_filing_rate"`
	SARFilingRate         float64 `json:"sar_filing_rate"`
	FalsePositiveRate     float64 `json:"false_positive_rate"`
	AverageResolutionTime int     `json:"avg_resolution_time_hours"`
	ComplianceScore       int     `json:"compliance_score"`
}

type AMLTrendAnalysis struct {
	AlertTrend30Days  []int     `json:"alert_trend_30_days"`
	VolumeTrend30Days []int64   `json:"volume_trend_30_days"`
	RiskScoreTrend    []float64 `json:"risk_score_trend"`
	EmergingPatterns  []string  `json:"emerging_patterns"`
}

type AMLRecommendation struct {
	Priority    string    `json:"priority"` // "HIGH", "MEDIUM", "LOW"
	Category    string    `json:"category"` // "RULE_TUNING", "INVESTIGATION", "TRAINING"
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ActionItems []string  `json:"action_items"`
	DueDate     time.Time `json:"due_date"`
}

// GenerateAMLDashboard creates a comprehensive AML monitoring dashboard
func (aml *AMLService) GenerateAMLDashboard(startDate, endDate time.Time) (*AMLDashboard, error) {
	dashboard := &AMLDashboard{
		PeriodStart:       startDate,
		PeriodEnd:         endDate,
		AlertsByRiskLevel: make(map[AMLRiskLevel]int),
		AlertsByType:      make(map[AMLRuleType]int),
	}

	// Get all alerts for the period
	alerts, err := aml.getAlertsForPeriod(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %w", err)
	}

	dashboard.TotalAlerts = len(alerts)

	// Analyze alerts by risk level and type
	for _, alert := range alerts {
		dashboard.AlertsByRiskLevel[alert.RiskLevel]++
		dashboard.AlertsByType[alert.RuleType]++
	}

	// Generate customer risk summaries
	dashboard.TopRiskyCustomers = aml.generateCustomerRiskSummaries(alerts, 10)

	// Calculate compliance metrics
	dashboard.ComplianceMetrics = aml.calculateComplianceMetrics(alerts)

	// Perform trend analysis
	dashboard.TrendAnalysis = aml.performTrendAnalysis(startDate, endDate)

	// Generate recommendations
	dashboard.RecommendedActions = aml.generateRecommendations(alerts, dashboard.ComplianceMetrics)

	return dashboard, nil
}

// getAlertsForPeriod retrieves alerts for a specific time period
func (aml *AMLService) getAlertsForPeriod(startDate, endDate time.Time) ([]*AMLAlert, error) {
	var alerts []*AMLAlert

	// In a real implementation, this would query the database
	// For now, return cached alerts that fall within the period
	for _, alert := range aml.alertsCache {
		if alert.DetectedAt.After(startDate) && alert.DetectedAt.Before(endDate) {
			alerts = append(alerts, alert)
		}
	}

	return alerts, nil
}

// generateCustomerRiskSummaries creates risk summaries for top risky customers
func (aml *AMLService) generateCustomerRiskSummaries(alerts []*AMLAlert, limit int) []CustomerRiskSummary {
	customerRisks := make(map[string]*CustomerRiskSummary)

	for _, alert := range alerts {
		if alert.EntityType != "CUSTOMER" {
			continue
		}

		if summary, exists := customerRisks[alert.EntityID]; exists {
			summary.AlertCount++
			if alert.Amount != nil {
				summary.TotalVolume += alert.Amount.Value
			}
			summary.RiskFactors = append(summary.RiskFactors, string(alert.RuleType))
		} else {
			volume := int64(0)
			if alert.Amount != nil {
				volume = alert.Amount.Value
			}

			customerRisks[alert.EntityID] = &CustomerRiskSummary{
				CustomerID:   alert.EntityID,
				CustomerName: fmt.Sprintf("Customer %s", alert.EntityID), // In real system, lookup actual name
				RiskScore:    70,                                         // Calculate based on alert types and frequencies
				AlertCount:   1,
				TotalVolume:  volume,
				LastActivity: alert.DetectedAt,
				RiskFactors:  []string{string(alert.RuleType)},
			}
		}
	}

	// Convert to slice and sort by risk score
	var summaries []CustomerRiskSummary
	for _, summary := range customerRisks {
		summaries = append(summaries, *summary)
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].RiskScore > summaries[j].RiskScore
	})

	if len(summaries) > limit {
		summaries = summaries[:limit]
	}

	return summaries
}

// calculateComplianceMetrics computes various compliance metrics
func (aml *AMLService) calculateComplianceMetrics(alerts []*AMLAlert) AMLComplianceMetrics {
	var ctrCount, sarCount, totalTransactions int
	var resolvedAlerts, falsePositives int
	var totalResolutionTime int

	for _, alert := range alerts {
		switch alert.RuleType {
		case RuleCTR:
			ctrCount++
		case RuleSAR:
			sarCount++
		}

		if alert.Status == "CLOSED" {
			resolvedAlerts++
			// Simulate resolution time (in a real system, calculate from timestamps)
			totalResolutionTime += 24 // Average 24 hours
		}

		// Simulate false positive identification
		if strings.Contains(alert.Description, "normal") {
			falsePositives++
		}
	}

	// Simulate total transactions for rate calculations
	totalTransactions = len(alerts) * 10 // Rough estimate

	avgResolutionTime := 0
	if resolvedAlerts > 0 {
		avgResolutionTime = totalResolutionTime / resolvedAlerts
	}

	falsePositiveRate := 0.0
	if len(alerts) > 0 {
		falsePositiveRate = float64(falsePositives) / float64(len(alerts)) * 100
	}

	// Calculate compliance score (0-100)
	complianceScore := 100
	if falsePositiveRate > 10 {
		complianceScore -= 20
	}
	if avgResolutionTime > 48 {
		complianceScore -= 15
	}

	return AMLComplianceMetrics{
		CTRFilingRate:         float64(ctrCount) / float64(totalTransactions) * 100,
		SARFilingRate:         float64(sarCount) / float64(totalTransactions) * 100,
		FalsePositiveRate:     falsePositiveRate,
		AverageResolutionTime: avgResolutionTime,
		ComplianceScore:       complianceScore,
	}
}

// performTrendAnalysis analyzes trends over the specified period
func (aml *AMLService) performTrendAnalysis(startDate, endDate time.Time) AMLTrendAnalysis {
	// Simulate trend data (in a real system, query historical data)
	return AMLTrendAnalysis{
		AlertTrend30Days:  []int{5, 7, 3, 8, 12, 6, 9, 4, 11, 8},
		VolumeTrend30Days: []int64{1000000, 1200000, 800000, 1500000, 2000000, 1100000, 1300000, 900000, 1600000, 1400000},
		RiskScoreTrend:    []float64{65.5, 67.2, 63.8, 69.4, 72.1, 66.9, 68.5, 64.2, 70.8, 69.3},
		EmergingPatterns: []string{
			"Increased cash transactions in hospitality sector",
			"Rise in round-amount wire transfers",
			"Growing cryptocurrency-related activity",
		},
	}
}

// generateRecommendations creates actionable recommendations based on analysis
func (aml *AMLService) generateRecommendations(alerts []*AMLAlert, metrics AMLComplianceMetrics) []AMLRecommendation {
	var recommendations []AMLRecommendation

	// High false positive rate recommendation
	if metrics.FalsePositiveRate > 15 {
		recommendations = append(recommendations, AMLRecommendation{
			Priority:    "HIGH",
			Category:    "RULE_TUNING",
			Title:       "Reduce False Positive Rate",
			Description: fmt.Sprintf("False positive rate is %.1f%%, exceeding acceptable threshold", metrics.FalsePositiveRate),
			ActionItems: []string{
				"Review and tune transaction amount thresholds",
				"Implement customer risk scoring refinements",
				"Train model with recent false positive data",
			},
			DueDate: time.Now().AddDate(0, 0, 14),
		})
	}

	// Long resolution time recommendation
	if metrics.AverageResolutionTime > 48 {
		recommendations = append(recommendations, AMLRecommendation{
			Priority:    "MEDIUM",
			Category:    "INVESTIGATION",
			Title:       "Improve Alert Resolution Time",
			Description: fmt.Sprintf("Average resolution time is %d hours, target is 48 hours", metrics.AverageResolutionTime),
			ActionItems: []string{
				"Automate initial alert triage",
				"Provide additional investigator training",
				"Implement priority-based alert routing",
			},
			DueDate: time.Now().AddDate(0, 0, 30),
		})
	}

	// High-risk patterns recommendation
	highRiskCount := 0
	for _, alert := range alerts {
		if alert.RiskLevel == RiskHigh || alert.RiskLevel == RiskCritical {
			highRiskCount++
		}
	}

	if highRiskCount > len(alerts)/4 {
		recommendations = append(recommendations, AMLRecommendation{
			Priority:    "HIGH",
			Category:    "INVESTIGATION",
			Title:       "Increased High-Risk Activity",
			Description: fmt.Sprintf("%d high/critical risk alerts require immediate attention", highRiskCount),
			ActionItems: []string{
				"Conduct enhanced due diligence on flagged customers",
				"Review transaction patterns for coordinated activity",
				"Consider filing SARs for suspicious patterns",
			},
			DueDate: time.Now().AddDate(0, 0, 7),
		})
	}

	return recommendations
}

// ExportAMLReport generates a detailed AML compliance report
func (aml *AMLService) ExportAMLReport(dashboard *AMLDashboard, format string) ([]byte, error) {
	switch format {
	case "JSON":
		return aml.exportJSONReport(dashboard)
	case "CSV":
		return aml.exportCSVReport(dashboard)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

func (aml *AMLService) exportJSONReport(dashboard *AMLDashboard) ([]byte, error) {
	return fmt.Appendf(nil, `{
  "report_generated": "%s",
  "period": "%s to %s",
  "total_alerts": %d,
  "high_risk_alerts": %d,
  "compliance_score": %d,
  "recommendations": %d
}`,
		time.Now().Format("2006-01-02 15:04:05"),
		dashboard.PeriodStart.Format("2006-01-02"),
		dashboard.PeriodEnd.Format("2006-01-02"),
		dashboard.TotalAlerts,
		dashboard.AlertsByRiskLevel[RiskHigh]+dashboard.AlertsByRiskLevel[RiskCritical],
		dashboard.ComplianceMetrics.ComplianceScore,
		len(dashboard.RecommendedActions),
	), nil
}

func (aml *AMLService) exportCSVReport(dashboard *AMLDashboard) ([]byte, error) {
	csv := "Alert_Type,Count,Risk_Level\n"
	for alertType, count := range dashboard.AlertsByType {
		csv += fmt.Sprintf("%s,%d,MEDIUM\n", alertType, count)
	}
	return []byte(csv), nil
}

// ----------------------------------------------------------------------------
// Comprehensive AML Check Methods
// ----------------------------------------------------------------------------

// CheckCashIntensiveActivity analyzes if a customer has unusually high cash activity
func (aml *AMLService) CheckCashIntensiveActivity(customerID string, timeWindow int) (*AMLAlert, error) {
	// Get customer transactions for the specified time window
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -timeWindow)

	query := &QueryOptions{
		Filters: []Filter{
			{Field: "valid_time", Operator: ">=", Value: startDate},
			{Field: "valid_time", Operator: "<=", Value: endDate},
		},
	}

	entries, err := aml.storage.QueryEntries(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query entries: %w", err)
	}

	var totalVolume, cashVolume int64
	var cashTransactions []string

	for _, entry := range entries {
		totalVolume += entry.Amount.Value

		// Check if this is a cash transaction (simplified check based on description/account)
		txn, err := aml.storage.GetTransaction(entry.TransactionID)
		if err != nil {
			continue
		}

		if strings.Contains(strings.ToLower(txn.Description), "cash") {
			cashVolume += entry.Amount.Value
			cashTransactions = append(cashTransactions, entry.TransactionID)
		}
	}

	if totalVolume == 0 {
		return nil, nil // No transactions
	}

	cashPercentage := float64(cashVolume) / float64(totalVolume) * 100

	// Check thresholds
	rule := aml.findRuleByType(RuleCashIntensive)
	if rule == nil {
		return nil, nil
	}

	minPercentage := rule.Thresholds["cash_percentage"].(float64)
	minVolume := int64(rule.Thresholds["minimum_volume"].(int))

	if cashPercentage >= minPercentage && totalVolume >= minVolume {
		return &AMLAlert{
			ID:             generateUUID(),
			RuleType:       RuleCashIntensive,
			Framework:      rule.Framework,
			RiskLevel:      RiskHigh,
			Title:          "Cash Intensive Activity Detected",
			Description:    fmt.Sprintf("Customer has %.1f%% cash transactions with volume $%.2f", cashPercentage, float64(totalVolume)/100),
			EntityID:       customerID,
			EntityType:     "CUSTOMER",
			TransactionIDs: cashTransactions,
			DetectedAt:     time.Now(),
			Status:         "OPEN",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}, nil
	}

	return nil, nil
}

// CheckJustUnderThreshold detects transactions just under reporting thresholds
func (aml *AMLService) CheckJustUnderThreshold(txn *Transaction) (*AMLAlert, error) {
	rule := aml.findRuleByType(RuleJustUnderThreshold)
	if rule == nil {
		return nil, nil
	}

	threshold := int64(rule.Thresholds["threshold_amount"].(int))
	tolerancePct := rule.Thresholds["tolerance_pct"].(float64)

	for _, entry := range txn.Entries {
		lowerBound := int64(float64(threshold) * (100 - tolerancePct) / 100)

		if entry.Amount.Value >= lowerBound && entry.Amount.Value < threshold {
			return &AMLAlert{
				ID:             generateUUID(),
				RuleType:       RuleJustUnderThreshold,
				Framework:      rule.Framework,
				RiskLevel:      RiskHigh,
				Title:          "Just Under Threshold Transaction",
				Description:    fmt.Sprintf("Transaction amount $%.2f is just under $%.2f threshold", float64(entry.Amount.Value)/100, float64(threshold)/100),
				EntityID:       txn.ID,
				EntityType:     "TRANSACTION",
				TransactionIDs: []string{txn.ID},
				Amount:         &entry.Amount,
				Currency:       string(entry.Amount.Currency),
				DetectedAt:     time.Now(),
				Status:         "OPEN",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}, nil
		}
	}

	return nil, nil
}

// CheckUnusualTiming detects transactions at unusual times
func (aml *AMLService) CheckUnusualTiming(txn *Transaction) (*AMLAlert, error) {
	rule := aml.findRuleByType(RuleUnusualTiming)
	if rule == nil {
		return nil, nil
	}

	nightStart := rule.Thresholds["night_start_hour"].(int)
	nightEnd := rule.Thresholds["night_end_hour"].(int)
	minAmount := int64(rule.Thresholds["minimum_amount"].(int))

	hour := txn.TransactionTime.Hour()
	isWeekend := txn.TransactionTime.Weekday() == time.Saturday || txn.TransactionTime.Weekday() == time.Sunday
	isNightTime := hour >= nightStart || hour <= nightEnd

	var totalAmount int64
	for _, entry := range txn.Entries {
		totalAmount += entry.Amount.Value
	}
	totalAmount /= 2 // Adjust for double-entry

	if (isNightTime || isWeekend) && totalAmount >= minAmount {
		timeDescription := "night time"
		if isWeekend {
			timeDescription = "weekend"
		}

		return &AMLAlert{
			ID:             generateUUID(),
			RuleType:       RuleUnusualTiming,
			Framework:      rule.Framework,
			RiskLevel:      RiskMedium,
			Title:          "Unusual Timing Transaction",
			Description:    fmt.Sprintf("$%.2f transaction during %s", float64(totalAmount)/100, timeDescription),
			EntityID:       txn.ID,
			EntityType:     "TRANSACTION",
			TransactionIDs: []string{txn.ID},
			Amount:         &Amount{Value: totalAmount, Currency: txn.Entries[0].Amount.Currency},
			Currency:       string(txn.Entries[0].Amount.Currency),
			DetectedAt:     time.Now(),
			Status:         "OPEN",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}, nil
	}

	return nil, nil
}

// CheckDormantAccountReactivation detects activity in previously dormant accounts
func (aml *AMLService) CheckDormantAccountReactivation(txn *Transaction) (*AMLAlert, error) {
	rule := aml.findRuleByType(RuleAccountDormancy)
	if rule == nil {
		return nil, nil
	}

	dormancyPeriod := rule.Thresholds["dormancy_period"].(int)
	minReactivationAmount := int64(rule.Thresholds["reactivation_amt"].(int))

	for _, entry := range txn.Entries {
		if entry.Amount.Value < minReactivationAmount {
			continue
		}

		// Check if account was dormant
		checkDate := time.Now().AddDate(0, 0, -dormancyPeriod)
		recentEntries, err := aml.storage.GetEntriesByAccount(entry.AccountID)
		if err != nil {
			continue
		}

		// Count transactions in dormancy period (excluding current transaction)
		var recentActivity int
		for _, recentEntry := range recentEntries {
			recentTxn, err := aml.storage.GetTransaction(recentEntry.TransactionID)
			if err != nil {
				continue
			}
			if recentTxn.TransactionTime.After(checkDate) && recentTxn.ID != txn.ID {
				recentActivity++
			}
		}

		if recentActivity == 0 { // Account was dormant
			return &AMLAlert{
				ID:             generateUUID(),
				RuleType:       RuleAccountDormancy,
				Framework:      rule.Framework,
				RiskLevel:      RiskMedium,
				Title:          "Dormant Account Reactivation",
				Description:    fmt.Sprintf("Account %s reactivated with $%.2f after %d days dormancy", entry.AccountID, float64(entry.Amount.Value)/100, dormancyPeriod),
				EntityID:       entry.AccountID,
				EntityType:     "ACCOUNT",
				AccountIDs:     []string{entry.AccountID},
				TransactionIDs: []string{txn.ID},
				Amount:         &entry.Amount,
				Currency:       string(entry.Amount.Currency),
				DetectedAt:     time.Now(),
				Status:         "OPEN",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}, nil
		}
	}

	return nil, nil
}

// CheckHighRiskGeography detects transactions from high-risk jurisdictions
func (aml *AMLService) CheckHighRiskGeography(txn *Transaction, customerInfo *AMLCustomer) (*AMLAlert, error) {
	rule := aml.findRuleByType(RuleUnexpectedGeography)
	if rule == nil || customerInfo == nil {
		return nil, nil
	}

	// Check if customer's country or transaction countries are high-risk
	highRiskCountries := rule.Countries
	isHighRisk := false
	riskCountry := ""

	for _, country := range highRiskCountries {
		if customerInfo.Country == country {
			isHighRisk = true
			riskCountry = country
			break
		}
	}

	if isHighRisk {
		var totalAmount int64
		for _, entry := range txn.Entries {
			totalAmount += entry.Amount.Value
		}
		totalAmount /= 2

		return &AMLAlert{
			ID:             generateUUID(),
			RuleType:       RuleUnexpectedGeography,
			Framework:      rule.Framework,
			RiskLevel:      RiskHigh,
			Title:          "High-Risk Geography Transaction",
			Description:    fmt.Sprintf("$%.2f transaction from high-risk country: %s", float64(totalAmount)/100, riskCountry),
			EntityID:       txn.ID,
			EntityType:     "TRANSACTION",
			TransactionIDs: []string{txn.ID},
			Amount:         &Amount{Value: totalAmount, Currency: txn.Entries[0].Amount.Currency},
			Currency:       string(txn.Entries[0].Amount.Currency),
			DetectedAt:     time.Now(),
			Status:         "OPEN",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}, nil
	}

	return nil, nil
}

// Helper function to find rule by type
func (aml *AMLService) findRuleByType(ruleType AMLRuleType) *AMLRule {
	for _, rule := range aml.rules {
		if rule.Type == ruleType && rule.Enabled {
			return rule
		}
	}
	return nil
}

// ----------------------------------------------------------------------------
