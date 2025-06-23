package main

import (
	"accounting"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// Clean up any existing test database
	dbPath := "aml_demo.db"
	os.Remove(dbPath)

	// Initialize accounting engine
	engine, err := accounting.NewAccountingEngine(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize engine: %v", err)
	}
	defer engine.Close()

	// Set up standard chart of accounts
	if err := engine.CreateStandardAccounts("admin"); err != nil {
		log.Fatalf("Failed to create chart of accounts: %v", err)
	}

	// Get AML service
	amlService := engine.GetAMLService()

	fmt.Println("🛡️ AML (Anti-Money Laundering) Demo")
	fmt.Println("===================================")

	// 1. Setup AML Rules
	fmt.Println("\n1. ✅ Setting up AML Rules...")

	// Setup BSA (Bank Secrecy Act) rules
	if err := amlService.SetupStandardAMLRules(accounting.BSA_Framework); err != nil {
		log.Fatalf("Failed to setup BSA rules: %v", err)
	}

	// Setup FATF rules
	if err := amlService.SetupStandardAMLRules(accounting.FATF_Framework); err != nil {
		log.Fatalf("Failed to setup FATF rules: %v", err)
	}

	// Setup EU AMLD rules
	if err := amlService.SetupStandardAMLRules(accounting.AMLD_Framework); err != nil {
		log.Fatalf("Failed to setup AMLD rules: %v", err)
	}

	fmt.Printf("   ✓ BSA rules (US) configured\n")
	fmt.Printf("   ✓ FATF rules (International) configured\n")
	fmt.Printf("   ✓ AMLD rules (EU) configured\n")

	// 2. Register AML Customers
	fmt.Println("\n2. ✅ Registering AML Customers...")

	// Normal customer
	normalCustomer := &accounting.AMLCustomer{
		ID:               "cust_001",
		CustomerID:       "customer_001",
		Name:             "John Smith",
		Type:             "INDIVIDUAL",
		RiskLevel:        accounting.RiskLow,
		Country:          "US",
		IsPEP:            false,
		IsHighRisk:       false,
		SanctionsMatch:   false,
		OnboardingDate:   time.Now().AddDate(-1, 0, 0), // 1 year ago
		ExpectedActivity: "Personal banking, salary deposits, utility payments",
		BusinessPurpose:  "Personal finance management",
	}

	// High-risk customer (PEP)
	pepCustomer := &accounting.AMLCustomer{
		ID:               "cust_002",
		CustomerID:       "customer_002",
		Name:             "Maria Rodriguez",
		Type:             "INDIVIDUAL",
		RiskLevel:        accounting.RiskHigh,
		Country:          "US",
		IsPEP:            true, // Politically Exposed Person
		IsHighRisk:       true,
		SanctionsMatch:   false,
		OnboardingDate:   time.Now().AddDate(0, -3, 0), // 3 months ago
		ExpectedActivity: "Large international transfers, investment activities",
		BusinessPurpose:  "Investment and wealth management",
	}

	// Sanctioned customer (for demo)
	sanctionedCustomer := &accounting.AMLCustomer{
		ID:               "cust_003",
		CustomerID:       "customer_003",
		Name:             "Suspected Entity Inc",
		Type:             "BUSINESS",
		RiskLevel:        accounting.RiskCritical,
		Country:          "IR", // Iran - high-risk jurisdiction
		IsPEP:            false,
		IsHighRisk:       true,
		SanctionsMatch:   true,                         // Matches sanctions list
		OnboardingDate:   time.Now().AddDate(0, -1, 0), // 1 month ago
		ExpectedActivity: "International trade",
		BusinessPurpose:  "Import/Export business",
	}

	// Register customers
	customers := []*accounting.AMLCustomer{normalCustomer, pepCustomer, sanctionedCustomer}
	for _, customer := range customers {
		if err := amlService.RegisterCustomer(customer); err != nil {
			log.Fatalf("Failed to register customer %s: %v", customer.Name, err)
		}
	}

	fmt.Printf("   ✓ Registered %d customers with varying risk levels\n", len(customers))

	// 3. Create Test Transactions and Monitor for AML Violations
	fmt.Println("\n3. ✅ Creating Transactions and Monitoring AML...")

	// Create customer info map for monitoring
	customerInfo := map[string]*accounting.AMLCustomer{
		"customer_001": normalCustomer,
		"customer_002": pepCustomer,
		"customer_003": sanctionedCustomer,
	}

	// Transaction 1: Normal transaction (should not trigger alerts)
	normalTxn := &accounting.Transaction{
		ValidTime:   time.Now(),
		Description: "Salary deposit - SAL-001",
		SourceRef:   "SAL-001",
		Entries: []accounting.Entry{
			{
				AccountID: "cash",
				Type:      accounting.Debit,
				Amount:    accounting.Amount{Value: 500000, Currency: "USD"}, // $5,000
			},
			{
				AccountID: "revenue",
				Type:      accounting.Credit,
				Amount:    accounting.Amount{Value: 500000, Currency: "USD"},
			},
		},
	}

	if err := engine.CreateTransaction(normalTxn, "system"); err != nil {
		log.Fatalf("Failed to create normal transaction: %v", err)
	}

	// Monitor the normal transaction
	alerts1, err := amlService.MonitorTransaction(normalTxn, customerInfo)
	if err != nil {
		log.Fatalf("Failed to monitor normal transaction: %v", err)
	}
	fmt.Printf("   Normal transaction ($5,000): %d alerts\n", len(alerts1))

	// Transaction 2: Large cash transaction (should trigger CTR)
	largeCashTxn := &accounting.Transaction{
		ValidTime:   time.Now(),
		Description: "Large cash deposit - CASH-001",
		SourceRef:   "CASH-001",
		Entries: []accounting.Entry{
			{
				AccountID: "cash",
				Type:      accounting.Debit,
				Amount:    accounting.Amount{Value: 1500000, Currency: "USD"}, // $15,000
			},
			{
				AccountID: "revenue",
				Type:      accounting.Credit,
				Amount:    accounting.Amount{Value: 1500000, Currency: "USD"},
			},
		},
	}

	if err := engine.CreateTransaction(largeCashTxn, "system"); err != nil {
		log.Fatalf("Failed to create large cash transaction: %v", err)
	}

	// Monitor the large cash transaction
	alerts2, err := amlService.MonitorTransaction(largeCashTxn, customerInfo)
	if err != nil {
		log.Fatalf("Failed to monitor large cash transaction: %v", err)
	}
	fmt.Printf("   Large cash transaction ($15,000): %d alerts\n", len(alerts2))

	// Transaction 3: Suspicious round amount transaction
	suspiciousTxn := &accounting.Transaction{
		ValidTime:   time.Now(),
		Description: "Wire transfer - WIRE-001",
		SourceRef:   "WIRE-001",
		Entries: []accounting.Entry{
			{
				AccountID: "cash",
				Type:      accounting.Debit,
				Amount:    accounting.Amount{Value: 1000000, Currency: "USD"}, // Exactly $10,000 (round amount)
			},
			{
				AccountID: "accounts_payable",
				Type:      accounting.Credit,
				Amount:    accounting.Amount{Value: 1000000, Currency: "USD"},
			},
		},
	}

	if err := engine.CreateTransaction(suspiciousTxn, "system"); err != nil {
		log.Fatalf("Failed to create suspicious transaction: %v", err)
	}

	// Monitor the suspicious transaction
	alerts3, err := amlService.MonitorTransaction(suspiciousTxn, customerInfo)
	if err != nil {
		log.Fatalf("Failed to monitor suspicious transaction: %v", err)
	}
	fmt.Printf("   Suspicious round amount ($10,000): %d alerts\n", len(alerts3))

	// Transaction 4: Sanctions violation
	sanctionsTxn := &accounting.Transaction{
		ValidTime:   time.Now(),
		Description: "International wire - INTL-001",
		SourceRef:   "INTL-001",
		Entries: []accounting.Entry{
			{
				AccountID: "cash",
				Type:      accounting.Debit,
				Amount:    accounting.Amount{Value: 750000, Currency: "USD"}, // $7,500
			},
			{
				AccountID: "accounts_payable",
				Type:      accounting.Credit,
				Amount:    accounting.Amount{Value: 750000, Currency: "USD"},
			},
		},
	}

	if err := engine.CreateTransaction(sanctionsTxn, "system"); err != nil {
		log.Fatalf("Failed to create sanctions transaction: %v", err)
	}

	// Monitor the sanctions transaction
	alerts4, err := amlService.MonitorTransaction(sanctionsTxn, customerInfo)
	if err != nil {
		log.Fatalf("Failed to monitor sanctions transaction: %v", err)
	}
	fmt.Printf("   Sanctions violation transaction: %d alerts\n", len(alerts4))

	// 4. Display AML Alerts
	fmt.Println("\n4. ✅ AML Alerts Summary...")

	allAlerts, err := amlService.GetAMLAlerts("", "", 0) // Get all alerts
	if err != nil {
		log.Fatalf("Failed to get AML alerts: %v", err)
	}

	fmt.Printf("   Total AML alerts generated: %d\n", len(allAlerts))

	for i, alert := range allAlerts {
		fmt.Printf("\n   Alert #%d:\n", i+1)
		fmt.Printf("     • Type: %s\n", alert.RuleType)
		fmt.Printf("     • Risk Level: %s\n", alert.RiskLevel)
		fmt.Printf("     • Title: %s\n", alert.Title)
		fmt.Printf("     • Description: %s\n", alert.Description)
		fmt.Printf("     • Amount: %s %.2f\n", alert.Currency, float64(alert.Amount.Value)/100)
		fmt.Printf("     • Status: %s\n", alert.Status)
		fmt.Printf("     • Evidence Count: %d\n", len(alert.Evidence))
	}

	// 5. Demonstrate Alert Management
	fmt.Println("\n5. ✅ Alert Management Demo...")

	if len(allAlerts) > 0 {
		// Create investigation for first high-risk alert
		var highRiskAlert *accounting.AMLAlert
		for _, alert := range allAlerts {
			if alert.RiskLevel == accounting.RiskHigh || alert.RiskLevel == accounting.RiskCritical {
				highRiskAlert = alert
				break
			}
		}

		if highRiskAlert != nil {
			investigation, err := amlService.CreateInvestigation(highRiskAlert.ID, "analyst_001")
			if err != nil {
				log.Fatalf("Failed to create investigation: %v", err)
			}
			fmt.Printf("   ✓ Created investigation %s for alert %s\n", investigation.ID, highRiskAlert.ID)

			// Add investigation note
			err = amlService.AddInvestigationNote(highRiskAlert.ID, "Initial review completed. Requires enhanced due diligence.", "analyst_001")
			if err != nil {
				log.Fatalf("Failed to add investigation note: %v", err)
			}
			fmt.Printf("   ✓ Added investigation note\n")
		}
	}

	// 6. KYC Management
	fmt.Println("\n6. ✅ KYC (Know Your Customer) Management...")

	// Perform KYC for high-risk customer
	if err := amlService.PerformKYC("cust_002", "compliance_officer"); err != nil {
		log.Fatalf("Failed to perform KYC: %v", err)
	}
	fmt.Printf("   ✓ Performed KYC for high-risk customer\n")

	// Check customers needing review
	needReview, err := amlService.GetCustomersForReview()
	if err != nil {
		log.Fatalf("Failed to get customers for review: %v", err)
	}
	fmt.Printf("   Customers needing KYC review: %d\n", len(needReview))

	// 7. Generate AML Reports
	fmt.Println("\n7. ✅ AML Reporting...")

	startDate := time.Now().AddDate(0, 0, -1) // Yesterday
	endDate := time.Now().AddDate(0, 0, 1)    // Tomorrow

	// Generate alerts summary report
	alertsReport, err := amlService.GenerateAMLReport("ALERTS_SUMMARY", startDate, endDate)
	if err != nil {
		log.Fatalf("Failed to generate alerts report: %v", err)
	}

	if reportMap, ok := alertsReport.(map[string]interface{}); ok {
		fmt.Printf("   Alerts Summary Report:\n")
		fmt.Printf("     • Period: %v\n", reportMap["period"])
		fmt.Printf("     • Total Alerts: %v\n", reportMap["total_alerts"])
		fmt.Printf("     • Risk Level Distribution: %v\n", reportMap["by_risk_level"])
		fmt.Printf("     • Rule Type Distribution: %v\n", reportMap["by_rule_type"])
	}

	// Generate risk assessment report
	riskReport, err := amlService.GenerateAMLReport("RISK_ASSESSMENT", startDate, endDate)
	if err != nil {
		log.Fatalf("Failed to generate risk report: %v", err)
	}

	if reportMap, ok := riskReport.(map[string]interface{}); ok {
		fmt.Printf("\n   Risk Assessment Report:\n")
		fmt.Printf("     • Total Customers: %v\n", reportMap["total_customers"])
		fmt.Printf("     • PEP Count: %v\n", reportMap["pep_count"])
		fmt.Printf("     • Sanctions Matches: %v\n", reportMap["sanctions_matches"])
		fmt.Printf("     • Risk Distribution: %v\n", reportMap["risk_distribution"])
	}

	// Generate CTR candidates report
	ctrReport, err := amlService.GenerateAMLReport("CTR_REPORT", startDate, endDate)
	if err != nil {
		log.Fatalf("Failed to generate CTR report: %v", err)
	}

	if ctrList, ok := ctrReport.([]map[string]interface{}); ok {
		fmt.Printf("\n   CTR (Currency Transaction Report) Candidates: %d\n", len(ctrList))
		for i, ctr := range ctrList {
			fmt.Printf("     CTR #%d: Transaction %v, Amount: %v\n", i+1, ctr["transaction_id"], ctr["amount"])
		}
	}

	// Generate SAR candidates report
	sarReport, err := amlService.GenerateAMLReport("SAR_CANDIDATES", startDate, endDate)
	if err != nil {
		log.Fatalf("Failed to generate SAR report: %v", err)
	}

	if sarList, ok := sarReport.([]map[string]interface{}); ok {
		fmt.Printf("\n   SAR (Suspicious Activity Report) Candidates: %d\n", len(sarList))
		for i, sar := range sarList {
			fmt.Printf("     SAR #%d: %v (%v) - %v\n", i+1, sar["rule_type"], sar["risk_level"], sar["description"])
		}
	}

	fmt.Println("\n🎉 AML Demo completed successfully!")
	fmt.Println("\nKey AML Features Demonstrated:")
	fmt.Println("✓ Multi-framework rule setup (BSA, FATF, AMLD)")
	fmt.Println("✓ Customer risk profiling and management")
	fmt.Println("✓ Real-time transaction monitoring")
	fmt.Println("✓ Automated alert generation")
	fmt.Println("✓ Investigation workflow")
	fmt.Println("✓ KYC/CDD management")
	fmt.Println("✓ Comprehensive AML reporting")
	fmt.Println("✓ Sanctions screening")
	fmt.Println("✓ CTR and SAR identification")
}
