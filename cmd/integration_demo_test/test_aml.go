package main_test

import (
	"accounting"
	"fmt"
	"os"
	"testing"
	"time"
)

// TestAMLIntegration tests AML functionality integration with the accounting system
func TestAMLIntegration(t *testing.T) {
	// Clean up any existing test database
	dbPath := "aml_integration_test.db"
	os.Remove(dbPath)
	defer os.Remove(dbPath)

	// Initialize accounting engine
	engine, err := accounting.NewAccountingEngine(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize engine: %v", err)
	}
	defer engine.Close()

	// Set up standard chart of accounts
	if err := engine.CreateStandardAccounts("admin"); err != nil {
		t.Fatalf("Failed to create chart of accounts: %v", err)
	}

	// Get AML service
	amlService := engine.GetAMLService()

	fmt.Println("🛡️ AML Integration Test")
	fmt.Println("=======================")

	// 1. Setup AML Rules
	fmt.Println("\n1. ✅ Setting up AML Rules...")

	if err := amlService.SetupStandardAMLRules(accounting.BSA_Framework); err != nil {
		t.Fatalf("Failed to setup BSA rules: %v", err)
	}

	if err := amlService.SetupStandardAMLRules(accounting.FATF_Framework); err != nil {
		t.Fatalf("Failed to setup FATF rules: %v", err)
	}

	// 2. Register test customers
	fmt.Println("\n2. ✅ Registering Test Customers...")

	testCustomers := []*accounting.AMLCustomer{
		{
			ID:               "cust_001",
			CustomerID:       "customer_001",
			Name:             "Test Customer 1",
			Type:             "INDIVIDUAL",
			RiskLevel:        accounting.RiskLow,
			Country:          "US",
			IsPEP:            false,
			IsHighRisk:       false,
			SanctionsMatch:   false,
			OnboardingDate:   time.Now().AddDate(-1, 0, 0),
			ExpectedActivity: "Normal banking",
			BusinessPurpose:  "Personal finance",
		},
		{
			ID:               "cust_002",
			CustomerID:       "customer_002",
			Name:             "High Risk Customer",
			Type:             "INDIVIDUAL",
			RiskLevel:        accounting.RiskHigh,
			Country:          "US",
			IsPEP:            true,
			IsHighRisk:       true,
			SanctionsMatch:   false,
			OnboardingDate:   time.Now().AddDate(0, -3, 0),
			ExpectedActivity: "Large transfers",
			BusinessPurpose:  "Investment",
		},
		{
			ID:               "cust_003",
			CustomerID:       "customer_003",
			Name:             "Sanctioned Entity",
			Type:             "BUSINESS",
			RiskLevel:        accounting.RiskCritical,
			Country:          "IR",
			IsPEP:            false,
			IsHighRisk:       true,
			SanctionsMatch:   true,
			OnboardingDate:   time.Now().AddDate(0, -1, 0),
			ExpectedActivity: "Trade",
			BusinessPurpose:  "Import/Export",
		},
	}

	for _, customer := range testCustomers {
		if err := amlService.RegisterCustomer(customer); err != nil {
			t.Fatalf("Failed to register customer %s: %v", customer.Name, err)
		}
	}

	customerInfo := map[string]*accounting.AMLCustomer{
		"customer_001": testCustomers[0],
		"customer_002": testCustomers[1],
		"customer_003": testCustomers[2],
	}

	// 3. Test transactions and monitoring
	fmt.Println("\n3. ✅ Creating Test Transactions...")

	// Normal transaction - should not trigger alerts
	normalTxn := &accounting.Transaction{
		ValidTime:   time.Now(),
		Description: "Normal payment",
		SourceRef:   "NORM-001",
		Entries: []accounting.Entry{
			{
				AccountID: "cash",
				Type:      accounting.Debit,
				Amount:    accounting.Amount{Value: 100000, Currency: "USD"}, // $1,000
			},
			{
				AccountID: "revenue",
				Type:      accounting.Credit,
				Amount:    accounting.Amount{Value: 100000, Currency: "USD"},
			},
		},
	}

	if err := engine.CreateTransaction(normalTxn, "system"); err != nil {
		t.Fatalf("Failed to create normal transaction: %v", err)
	}

	// Monitor normal transaction
	normalAlerts, err := amlService.MonitorTransaction(normalTxn, customerInfo)
	if err != nil {
		t.Fatalf("Failed to monitor normal transaction: %v", err)
	}

	// Large cash transaction - should trigger CTR alert
	largeCashTxn := &accounting.Transaction{
		ValidTime:   time.Now(),
		Description: "Large cash deposit",
		SourceRef:   "CASH-001",
		Entries: []accounting.Entry{
			{
				AccountID: "cash",
				Type:      accounting.Debit,
				Amount:    accounting.Amount{Value: 1200000, Currency: "USD"}, // $12,000
			},
			{
				AccountID: "revenue",
				Type:      accounting.Credit,
				Amount:    accounting.Amount{Value: 1200000, Currency: "USD"},
			},
		},
	}

	if err := engine.CreateTransaction(largeCashTxn, "system"); err != nil {
		t.Fatalf("Failed to create large cash transaction: %v", err)
	}

	// Monitor large cash transaction
	ctrAlerts, err := amlService.MonitorTransaction(largeCashTxn, customerInfo)
	if err != nil {
		t.Fatalf("Failed to monitor large cash transaction: %v", err)
	}

	// Suspicious round amount transaction
	suspiciousTxn := &accounting.Transaction{
		ValidTime:   time.Now(),
		Description: "Suspicious wire",
		SourceRef:   "WIRE-001",
		Entries: []accounting.Entry{
			{
				AccountID: "cash",
				Type:      accounting.Debit,
				Amount:    accounting.Amount{Value: 1000000, Currency: "USD"}, // Exactly $10,000
			},
			{
				AccountID: "accounts_payable",
				Type:      accounting.Credit,
				Amount:    accounting.Amount{Value: 1000000, Currency: "USD"},
			},
		},
	}

	if err := engine.CreateTransaction(suspiciousTxn, "system"); err != nil {
		t.Fatalf("Failed to create suspicious transaction: %v", err)
	}

	// Monitor suspicious transaction
	suspiciousAlerts, err := amlService.MonitorTransaction(suspiciousTxn, customerInfo)
	if err != nil {
		t.Fatalf("Failed to monitor suspicious transaction: %v", err)
	}

	// Sanctions violation transaction
	sanctionsTxn := &accounting.Transaction{
		ValidTime:   time.Now(),
		Description: "International payment",
		SourceRef:   "INTL-001",
		Entries: []accounting.Entry{
			{
				AccountID: "cash",
				Type:      accounting.Debit,
				Amount:    accounting.Amount{Value: 500000, Currency: "USD"}, // $5,000
			},
			{
				AccountID: "accounts_payable",
				Type:      accounting.Credit,
				Amount:    accounting.Amount{Value: 500000, Currency: "USD"},
			},
		},
	}

	if err := engine.CreateTransaction(sanctionsTxn, "system"); err != nil {
		t.Fatalf("Failed to create sanctions transaction: %v", err)
	}

	// Monitor sanctions transaction
	sanctionsAlerts, err := amlService.MonitorTransaction(sanctionsTxn, customerInfo)
	if err != nil {
		t.Fatalf("Failed to monitor sanctions transaction: %v", err)
	}

	// 4. Verify alert generation
	fmt.Println("\n4. ✅ Verifying AML Alert Generation...")

	fmt.Printf("   Normal transaction alerts: %d\n", len(normalAlerts))
	fmt.Printf("   CTR alerts: %d\n", len(ctrAlerts))
	fmt.Printf("   Suspicious transaction alerts: %d\n", len(suspiciousAlerts))
	fmt.Printf("   Sanctions alerts: %d\n", len(sanctionsAlerts))

	// Verify we have appropriate alerts
	if len(normalAlerts) != 0 {
		t.Errorf("Expected 0 alerts for normal transaction, got %d", len(normalAlerts))
	}

	// Should have at least 1 alert for large cash (CTR or SAR)
	if len(ctrAlerts) == 0 {
		t.Errorf("Expected at least 1 alert for large cash transaction")
	}

	// Should have at least 1 alert for sanctions
	if len(sanctionsAlerts) == 0 {
		t.Errorf("Expected at least 1 alert for sanctions transaction")
	}

	// 5. Test alert retrieval and management
	fmt.Println("\n5. ✅ Testing Alert Management...")

	allAlerts, err := amlService.GetAMLAlerts("", "", 0)
	if err != nil {
		t.Fatalf("Failed to get all alerts: %v", err)
	}

	fmt.Printf("   Total alerts in system: %d\n", len(allAlerts))

	// Test filtering by risk level
	highRiskAlerts, err := amlService.GetAMLAlerts("", accounting.RiskHigh, 0)
	if err != nil {
		t.Fatalf("Failed to get high-risk alerts: %v", err)
	}

	fmt.Printf("   High-risk alerts: %d\n", len(highRiskAlerts))

	// Test alert status update
	if len(allAlerts) > 0 {
		testAlert := allAlerts[0]
		if err := amlService.UpdateAlertStatus(testAlert.ID, "INVESTIGATING", "analyst_001"); err != nil {
			t.Fatalf("Failed to update alert status: %v", err)
		}
		fmt.Printf("   ✓ Updated alert status to INVESTIGATING\n")
	}

	// 6. Test investigation workflow
	fmt.Println("\n6. ✅ Testing Investigation Workflow...")

	if len(allAlerts) > 0 {
		testAlert := allAlerts[0]

		// Create investigation
		investigation, err := amlService.CreateInvestigation(testAlert.ID, "investigator_001")
		if err != nil {
			t.Fatalf("Failed to create investigation: %v", err)
		}
		fmt.Printf("   ✓ Created investigation: %s\n", investigation.ID)

		// Add investigation note
		if err := amlService.AddInvestigationNote(testAlert.ID, "Initial review completed", "investigator_001"); err != nil {
			t.Fatalf("Failed to add investigation note: %v", err)
		}
		fmt.Printf("   ✓ Added investigation note\n")
	}

	// 7. Test KYC functionality
	fmt.Println("\n7. ✅ Testing KYC Management...")

	// Perform KYC
	if err := amlService.PerformKYC("cust_002", "compliance_officer"); err != nil {
		t.Fatalf("Failed to perform KYC: %v", err)
	}
	fmt.Printf("   ✓ Performed KYC for high-risk customer\n")

	// Test customer risk update
	if err := amlService.UpdateCustomerRisk("cust_001", accounting.RiskMedium, "Increased activity volume"); err != nil {
		t.Fatalf("Failed to update customer risk: %v", err)
	}
	fmt.Printf("   ✓ Updated customer risk level\n")

	// 8. Test reporting
	fmt.Println("\n8. ✅ Testing AML Reporting...")

	startDate := time.Now().AddDate(0, 0, -1)
	endDate := time.Now().AddDate(0, 0, 1)

	// Generate alerts summary
	alertsReport, err := amlService.GenerateAMLReport("ALERTS_SUMMARY", startDate, endDate)
	if err != nil {
		t.Fatalf("Failed to generate alerts report: %v", err)
	}

	if reportMap, ok := alertsReport.(map[string]interface{}); ok {
		fmt.Printf("   ✓ Generated alerts summary: %v total alerts\n", reportMap["total_alerts"])
	}

	// Generate risk assessment
	riskReport, err := amlService.GenerateAMLReport("RISK_ASSESSMENT", startDate, endDate)
	if err != nil {
		t.Fatalf("Failed to generate risk assessment: %v", err)
	}

	if reportMap, ok := riskReport.(map[string]interface{}); ok {
		fmt.Printf("   ✓ Generated risk assessment: %v customers\n", reportMap["total_customers"])
	}

	fmt.Println("\n🎉 AML Integration Test completed successfully!")
	fmt.Println("\nAML Features Tested:")
	fmt.Println("✓ Multi-framework rule setup")
	fmt.Println("✓ Customer registration and risk management")
	fmt.Println("✓ Real-time transaction monitoring")
	fmt.Println("✓ Alert generation and filtering")
	fmt.Println("✓ Investigation workflow")
	fmt.Println("✓ KYC management")
	fmt.Println("✓ AML reporting")
}
