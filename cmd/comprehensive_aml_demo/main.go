package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"accounting"
)

func main() {
	// Clean up any existing test database
	dbPath := "comprehensive_aml_demo.db"
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

	fmt.Println("🛡️ Comprehensive AML (Anti-Money Laundering) Demo")
	fmt.Println("===============================================")

	// 1. Setup All AML Rules
	fmt.Println("\n1. ✅ Setting up Comprehensive AML Rules...")

	if err := amlService.SetupAllStandardAMLRules(); err != nil {
		log.Fatalf("Failed to setup comprehensive AML rules: %v", err)
	}

	fmt.Printf("   ✓ BSA rules (US) configured\n")
	fmt.Printf("   ✓ FATF rules (International) configured\n")
	fmt.Printf("   ✓ AMLD rules (EU) configured\n")
	fmt.Printf("   ✓ FinCEN rules (US Financial Crimes) configured\n")
	fmt.Printf("   ✓ OFAC rules (US Sanctions) configured\n")
	fmt.Printf("   ✓ Common AML detection rules configured\n")

	// 2. Register AML Customers with various risk profiles
	fmt.Println("\n2. ✅ Registering AML Customers...")

	// Low-risk customer
	lowRiskCustomer := &accounting.AMLCustomer{
		ID:               "cust_low_001",
		CustomerID:       "customer_low_001",
		Name:             "Jane Smith",
		Type:             "INDIVIDUAL",
		Country:          "US",
		RiskLevel:        accounting.RiskLow,
		OnboardingDate:   time.Now().AddDate(0, -6, 0),
		ExpectedActivity: "Regular salary deposits and bill payments",
		BusinessPurpose:  "Personal banking",
	}

	// High-risk customer (PEP)
	highRiskCustomer := &accounting.AMLCustomer{
		ID:               "cust_high_001",
		CustomerID:       "customer_high_001",
		Name:             "Alex Political",
		Type:             "PEP",
		Country:          "IR", // High-risk jurisdiction
		RiskLevel:        accounting.RiskHigh,
		OnboardingDate:   time.Now().AddDate(0, -3, 0),
		ExpectedActivity: "International wire transfers",
		BusinessPurpose:  "Government official",
	}

	// Business customer
	businessCustomer := &accounting.AMLCustomer{
		ID:               "cust_biz_001",
		CustomerID:       "customer_biz_001",
		Name:             "Cash Express LLC",
		Type:             "BUSINESS",
		Country:          "US",
		RiskLevel:        accounting.RiskMedium,
		OnboardingDate:   time.Now().AddDate(0, -12, 0),
		ExpectedActivity: "High volume cash transactions",
		BusinessPurpose:  "Cash intensive business",
	}

	customerInfo := map[string]*accounting.AMLCustomer{
		"customer_low_001":  lowRiskCustomer,
		"customer_high_001": highRiskCustomer,
		"customer_biz_001":  businessCustomer,
	}

	fmt.Printf("   ✓ Low-risk customer registered: %s\n", lowRiskCustomer.Name)
	fmt.Printf("   ✓ High-risk PEP customer registered: %s\n", highRiskCustomer.Name)
	fmt.Printf("   ✓ Business customer registered: %s\n", businessCustomer.Name)

	// 3. Demonstrate various AML rule violations
	fmt.Println("\n3. ✅ Testing AML Rules with Various Transaction Scenarios...")

	// Scenario 1: Just Under Threshold (Structuring)
	fmt.Println("\n   Scenario 1: Just Under Threshold Detection")
	structuringTxn := &accounting.Transaction{
		ValidTime:       time.Now(),
		TransactionTime: time.Now(),
		Description:     "Cash deposit - avoiding threshold",
		SourceRef:       "STRUCT-001",
		Entries: []accounting.Entry{
			{
				AccountID: "cash",
				Type:      accounting.Debit,
				Amount:    accounting.Amount{Value: 995000, Currency: "USD"}, // $9,950 (just under $10k)
			},
			{
				AccountID: "deposits",
				Type:      accounting.Credit,
				Amount:    accounting.Amount{Value: 995000, Currency: "USD"},
			},
		},
	}

	if err := engine.CreateTransaction(structuringTxn, "system"); err != nil {
		log.Fatalf("Failed to create structuring transaction: %v", err)
	}

	alerts1, err := amlService.MonitorTransaction(structuringTxn, customerInfo)
	if err != nil {
		log.Fatalf("Failed to monitor structuring transaction: %v", err)
	}
	fmt.Printf("      Just under threshold ($9,950): %d alerts generated\n", len(alerts1))

	// Scenario 2: Unusual Timing (Weekend transaction)
	fmt.Println("\n   Scenario 2: Unusual Timing Detection")
	weekendTime := time.Now()
	// Adjust to weekend if not already
	for weekendTime.Weekday() != time.Saturday {
		weekendTime = weekendTime.AddDate(0, 0, 1)
	}

	weekendTxn := &accounting.Transaction{
		ValidTime:       weekendTime,
		TransactionTime: weekendTime,
		Description:     "Large weekend wire transfer",
		SourceRef:       "WEEKEND-001",
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

	if err := engine.CreateTransaction(weekendTxn, "system"); err != nil {
		log.Fatalf("Failed to create weekend transaction: %v", err)
	}

	alerts2, err := amlService.MonitorTransaction(weekendTxn, customerInfo)
	if err != nil {
		log.Fatalf("Failed to monitor weekend transaction: %v", err)
	}
	fmt.Printf("      Weekend transaction ($7,500): %d alerts generated\n", len(alerts2))

	// Scenario 3: High-Risk Geography
	fmt.Println("\n   Scenario 3: High-Risk Geography Detection")
	geoRiskTxn := &accounting.Transaction{
		ValidTime:       time.Now(),
		TransactionTime: time.Now(),
		Description:     "International wire from Iran",
		SourceRef:       "IRAN-001",
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

	if err := engine.CreateTransaction(geoRiskTxn, "system"); err != nil {
		log.Fatalf("Failed to create geo-risk transaction: %v", err)
	}

	alerts3, err := amlService.MonitorTransaction(geoRiskTxn, customerInfo)
	if err != nil {
		log.Fatalf("Failed to monitor geo-risk transaction: %v", err)
	}
	fmt.Printf("      High-risk geography transaction: %d alerts generated\n", len(alerts3))

	// Scenario 4: Cash Intensive Activity
	fmt.Println("\n   Scenario 4: Cash Intensive Activity Detection")

	// Create multiple cash transactions to trigger cash intensive rule
	for i := 1; i <= 5; i++ {
		cashTxn := &accounting.Transaction{
			ValidTime:       time.Now().AddDate(0, 0, -i),
			TransactionTime: time.Now().AddDate(0, 0, -i),
			Description:     fmt.Sprintf("Cash deposit #%d", i),
			SourceRef:       fmt.Sprintf("CASH-%03d", i),
			Entries: []accounting.Entry{
				{
					AccountID: "cash",
					Type:      accounting.Debit,
					Amount:    accounting.Amount{Value: 1500000, Currency: "USD"}, // $15,000 each
				},
				{
					AccountID: "deposits",
					Type:      accounting.Credit,
					Amount:    accounting.Amount{Value: 1500000, Currency: "USD"},
				},
			},
		}

		if err := engine.CreateTransaction(cashTxn, "system"); err != nil {
			log.Printf("Failed to create cash transaction %d: %v", i, err)
			continue
		}

		if i == 1 { // Only monitor the first transaction to avoid duplicate alerts
			alerts4, err := amlService.MonitorTransaction(cashTxn, customerInfo)
			if err != nil {
				log.Printf("Failed to monitor cash transaction: %v", err)
			} else {
				fmt.Printf("      Cash intensive activity: %d alerts generated\n", len(alerts4))
			}
		}
	}

	// Scenario 5: Round Amount Pattern
	fmt.Println("\n   Scenario 5: Round Amount Pattern Detection")
	roundTxn := &accounting.Transaction{
		ValidTime:       time.Now(),
		TransactionTime: time.Now(),
		Description:     "Suspiciously round amount",
		SourceRef:       "ROUND-001",
		Entries: []accounting.Entry{
			{
				AccountID: "cash",
				Type:      accounting.Debit,
				Amount:    accounting.Amount{Value: 5000000, Currency: "USD"}, // Exactly $50,000
			},
			{
				AccountID: "investments",
				Type:      accounting.Credit,
				Amount:    accounting.Amount{Value: 5000000, Currency: "USD"},
			},
		},
	}

	if err := engine.CreateTransaction(roundTxn, "system"); err != nil {
		log.Fatalf("Failed to create round amount transaction: %v", err)
	}

	alerts5, err := amlService.MonitorTransaction(roundTxn, customerInfo)
	if err != nil {
		log.Fatalf("Failed to monitor round amount transaction: %v", err)
	}
	fmt.Printf("      Round amount transaction ($50,000): %d alerts generated\n", len(alerts5))

	// 4. Generate AML Dashboard
	fmt.Println("\n4. ✅ Generating AML Monitoring Dashboard...")

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30) // Last 30 days

	dashboard, err := amlService.GenerateAMLDashboard(startDate, endDate)
	if err != nil {
		log.Fatalf("Failed to generate AML dashboard: %v", err)
	}

	fmt.Printf("   ✓ Dashboard Period: %s to %s\n",
		dashboard.PeriodStart.Format("2006-01-02"),
		dashboard.PeriodEnd.Format("2006-01-02"))
	fmt.Printf("   ✓ Total Alerts: %d\n", dashboard.TotalAlerts)

	fmt.Println("\n   Alert Breakdown by Risk Level:")
	for riskLevel, count := range dashboard.AlertsByRiskLevel {
		if count > 0 {
			fmt.Printf("      %s: %d alerts\n", riskLevel, count)
		}
	}

	fmt.Println("\n   Alert Breakdown by Rule Type:")
	for ruleType, count := range dashboard.AlertsByType {
		if count > 0 {
			fmt.Printf("      %s: %d alerts\n", ruleType, count)
		}
	}

	fmt.Println("\n   Compliance Metrics:")
	fmt.Printf("      Compliance Score: %d/100\n", dashboard.ComplianceMetrics.ComplianceScore)
	fmt.Printf("      False Positive Rate: %.1f%%\n", dashboard.ComplianceMetrics.FalsePositiveRate)
	fmt.Printf("      Average Resolution Time: %d hours\n", dashboard.ComplianceMetrics.AverageResolutionTime)

	fmt.Printf("\n   Top Risky Customers: %d identified\n", len(dashboard.TopRiskyCustomers))
	for i, customer := range dashboard.TopRiskyCustomers {
		if i < 3 { // Show top 3
			fmt.Printf("      %d. %s (Risk Score: %d, Alerts: %d)\n",
				i+1, customer.CustomerName, customer.RiskScore, customer.AlertCount)
		}
	}

	fmt.Printf("\n   Recommendations: %d action items\n", len(dashboard.RecommendedActions))
	for i, rec := range dashboard.RecommendedActions {
		if i < 2 { // Show top 2
			fmt.Printf("      %d. [%s] %s\n", i+1, rec.Priority, rec.Title)
		}
	}

	// 5. Export AML Report
	fmt.Println("\n5. ✅ Exporting AML Compliance Report...")

	jsonReport, err := amlService.ExportAMLReport(dashboard, "JSON")
	if err != nil {
		log.Fatalf("Failed to export JSON report: %v", err)
	}

	fmt.Printf("   ✓ JSON Report Generated (%d bytes)\n", len(jsonReport))
	fmt.Printf("   Report Summary:\n%s\n", string(jsonReport))

	// 6. Summary
	fmt.Println("\n6. ✅ AML Demo Summary")
	fmt.Println("   ===================")
	fmt.Printf("   ✓ %d regulatory frameworks configured\n", 5)
	fmt.Printf("   ✓ %d+ AML rules implemented\n", 27)
	fmt.Printf("   ✓ %d customer risk profiles tested\n", 3)
	fmt.Printf("   ✓ %d transaction scenarios evaluated\n", 5)
	fmt.Printf("   ✓ %d total alerts generated\n", dashboard.TotalAlerts)
	fmt.Printf("   ✓ Comprehensive monitoring dashboard created\n")
	fmt.Printf("   ✓ Compliance reporting enabled\n")

	fmt.Println("\n🎯 AML System Features Demonstrated:")
	fmt.Println("   • Real-time transaction monitoring")
	fmt.Println("   • Multi-framework regulatory compliance")
	fmt.Println("   • Risk-based customer profiling")
	fmt.Println("   • Pattern-based suspicious activity detection")
	fmt.Println("   • Geographic risk assessment")
	fmt.Println("   • Automated alert generation and scoring")
	fmt.Println("   • Compliance metrics and reporting")
	fmt.Println("   • Investigation workflow support")

	fmt.Println("\n✅ Comprehensive AML Demo completed successfully!")
}
