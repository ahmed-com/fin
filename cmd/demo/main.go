package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"accounting"
)

func main() {
	fmt.Println("🏦 Advanced Accounting System Demo")
	fmt.Println("===================================")

	// Clean up any existing database
	dbFile := "demo_accounting.db"
	os.Remove(dbFile)

	// Initialize the accounting engine
	engine, err := accounting.NewAccountingEngine(dbFile)
	if err != nil {
		log.Fatalf("Failed to create accounting engine: %v", err)
	}
	defer engine.Close()
	defer os.Remove(dbFile) // Clean up after demo

	userID := "demo_user"

	// Demo 1: Create Chart of Accounts
	fmt.Println("\n📊 Step 1: Creating Chart of Accounts")
	if err := engine.CreateStandardAccounts(userID); err != nil {
		log.Fatalf("Failed to create accounts: %v", err)
	}
	fmt.Println("✅ Standard chart of accounts created")

	// Demo 2: Record a Sale Transaction
	fmt.Println("\n💰 Step 2: Recording a Sale Transaction")
	saleTransaction := &accounting.Transaction{
		Description: "Sale of consulting services",
		ValidTime:   time.Now(),
		Entries: []accounting.Entry{
			{
				AccountID: "cash",
				Type:      accounting.Debit,
				Amount: accounting.Amount{
					Value:    250000, // $2,500.00 in cents
					Currency: "USD",
				},
				Dimensions: []accounting.Dimension{
					{Key: accounting.DimDepartment, Value: "consulting"},
					{Key: accounting.DimProject, Value: "project_alpha"},
				},
			},
			{
				AccountID: "revenue",
				Type:      accounting.Credit,
				Amount: accounting.Amount{
					Value:    250000,
					Currency: "USD",
				},
				Dimensions: []accounting.Dimension{
					{Key: accounting.DimDepartment, Value: "consulting"},
					{Key: accounting.DimProject, Value: "project_alpha"},
				},
			},
		},
	}

	if err := engine.CreateTransaction(saleTransaction, userID); err != nil {
		log.Fatalf("Failed to create transaction: %v", err)
	}

	if err := engine.PostTransaction(saleTransaction.ID, userID); err != nil {
		log.Fatalf("Failed to post transaction: %v", err)
	}
	fmt.Printf("✅ Sale transaction posted: %s\n", saleTransaction.ID[:8]+"...")

	// Demo 3: Record an Expense Transaction
	fmt.Println("\n📝 Step 3: Recording an Expense Transaction")
	expenseTransaction := &accounting.Transaction{
		Description: "Office supplies purchase",
		ValidTime:   time.Now(),
		Entries: []accounting.Entry{
			{
				AccountID: "expenses",
				Type:      accounting.Debit,
				Amount: accounting.Amount{
					Value:    15000, // $150.00
					Currency: "USD",
				},
				Dimensions: []accounting.Dimension{
					{Key: accounting.DimDepartment, Value: "admin"},
					{Key: accounting.DimCostCenter, Value: "office"},
				},
			},
			{
				AccountID: "cash",
				Type:      accounting.Credit,
				Amount: accounting.Amount{
					Value:    15000,
					Currency: "USD",
				},
			},
		},
	}

	if err := engine.CreateTransaction(expenseTransaction, userID); err != nil {
		log.Fatalf("Failed to create expense transaction: %v", err)
	}

	if err := engine.PostTransaction(expenseTransaction.ID, userID); err != nil {
		log.Fatalf("Failed to post expense transaction: %v", err)
	}
	fmt.Printf("✅ Expense transaction posted: %s\n", expenseTransaction.ID[:8]+"...")

	// Demo 4: Check Account Balances
	fmt.Println("\n📊 Step 4: Checking Account Balances")
	accounts := []string{"cash", "revenue", "expenses"}
	for _, accountID := range accounts {
		balance, err := engine.GetAccountBalance(accountID, time.Now())
		if err != nil {
			fmt.Printf("❌ Error getting %s balance: %v\n", accountID, err)
			continue
		}
		fmt.Printf("   %s: $%.2f %s\n",
			balance.AccountName,
			float64(balance.Balance.Value)/100,
			balance.Balance.Currency)
	}

	// Demo 5: Create Subscription with Accrual Schedule
	fmt.Println("\n📅 Step 5: Creating Subscription with Accrual Recognition")
	subscriptionTransaction := &accounting.Transaction{
		Description: "Annual software subscription - prepaid",
		ValidTime:   time.Now(),
		Entries: []accounting.Entry{
			{
				AccountID: "cash",
				Type:      accounting.Debit,
				Amount: accounting.Amount{
					Value:    120000, // $1,200.00 for 12 months
					Currency: "USD",
				},
			},
			{
				AccountID: "unearned_revenue",
				Type:      accounting.Credit,
				Amount: accounting.Amount{
					Value:    120000,
					Currency: "USD",
				},
			},
		},
	}

	if err := engine.CreateTransaction(subscriptionTransaction, userID); err != nil {
		log.Fatalf("Failed to create subscription transaction: %v", err)
	}

	if err := engine.PostTransaction(subscriptionTransaction.ID, userID); err != nil {
		log.Fatalf("Failed to post subscription transaction: %v", err)
	}

	// Create accrual schedule for monthly revenue recognition
	template := &accounting.AccrualTemplate{
		ID:               "revenue_recognition",
		Name:             "Monthly Revenue Recognition",
		AccrualType:      accounting.DeferralRevenue,
		RevenueAccountID: "revenue",
	}

	schedule, err := engine.CreateAccrualSchedule(
		subscriptionTransaction.ID,
		&accounting.Amount{Value: 120000, Currency: "USD"},
		accounting.Monthly,
		12,
		time.Now(),
		template,
		userID,
	)
	if err != nil {
		log.Fatalf("Failed to create accrual schedule: %v", err)
	}
	fmt.Printf("✅ Subscription recorded with 12-month accrual schedule: %s\n", schedule.ID[:8]+"...")

	// Demo 6: Generate Trial Balance
	fmt.Println("\n📋 Step 6: Generating Trial Balance")
	trialBalance, err := engine.GetTrialBalance(time.Now(), nil)
	if err != nil {
		log.Fatalf("Failed to generate trial balance: %v", err)
	}

	fmt.Println("\n   TRIAL BALANCE")
	fmt.Println("   Account                     Type        Balance")
	fmt.Println("   =============================================")

	totalDebits := int64(0)
	totalCredits := int64(0)

	for _, balance := range trialBalance {
		fmt.Printf("   %-25s %-10s $%8.2f\n",
			balance.AccountName,
			balance.AccountType,
			float64(balance.Balance.Value)/100)

		// Calculate totals based on normal balance sides
		switch balance.AccountType {
		case accounting.Asset, accounting.Expense:
			if balance.Balance.Value > 0 {
				totalDebits += balance.Balance.Value
			} else {
				totalCredits += -balance.Balance.Value
			}
		case accounting.Liability, accounting.Equity, accounting.Income:
			if balance.Balance.Value > 0 {
				totalCredits += balance.Balance.Value
			} else {
				totalDebits += -balance.Balance.Value
			}
		}
	}

	fmt.Println("   =============================================")
	fmt.Printf("   Total Debits:               $%8.2f\n", float64(totalDebits)/100)
	fmt.Printf("   Total Credits:              $%8.2f\n", float64(totalCredits)/100)
	fmt.Printf("   Difference:                 $%8.2f\n", float64(totalDebits-totalCredits)/100)

	// Demo 7: Create and Close Accounting Period
	fmt.Println("\n📆 Step 7: Managing Accounting Periods")
	period := &accounting.Period{
		Name:  "2025-Q2",
		Start: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2025, 6, 30, 23, 59, 59, 0, time.UTC),
	}

	if err := engine.CreatePeriod(period, userID); err != nil {
		log.Fatalf("Failed to create period: %v", err)
	}
	fmt.Printf("✅ Period created: %s (%s to %s)\n",
		period.Name,
		period.Start.Format("2006-01-02"),
		period.End.Format("2006-01-02"))

	// Soft close the period
	if err := engine.ClosePeriod(period.ID, true, userID); err != nil {
		log.Fatalf("Failed to close period: %v", err)
	}
	fmt.Println("✅ Period soft-closed (still editable)")

	// Demo 8: Show Event Sourcing Capabilities
	fmt.Println("\n📜 Step 8: Event Sourcing - Audit Trail")
	events, err := engine.GetEvents(time.Now().Add(-1*time.Hour), time.Now())
	if err != nil {
		log.Fatalf("Failed to get events: %v", err)
	}

	fmt.Printf("✅ Retrieved %d events from the audit trail:\n", len(events))
	for i, event := range events {
		if i < 5 { // Show first 5 events
			fmt.Printf("   %d. %s at %s by %s\n",
				i+1,
				event.EventType,
				event.TransactionTime.Format("15:04:05"),
				event.UserID)
		}
	}
	if len(events) > 5 {
		fmt.Printf("   ... and %d more events\n", len(events)-5)
	}

	// Demo 9: Demonstrate Transaction Reversal
	fmt.Println("\n🔄 Step 9: Transaction Reversal")
	reversalTxn, err := engine.ReverseTransaction(
		expenseTransaction.ID,
		"Reversal of office supplies - returned items",
		userID,
	)
	if err != nil {
		log.Fatalf("Failed to reverse transaction: %v", err)
	}
	fmt.Printf("✅ Transaction reversed: %s\n", reversalTxn.ID[:8]+"...")

	// Show updated cash balance
	cashBalance, err := engine.GetAccountBalance("cash", time.Now())
	if err != nil {
		log.Fatalf("Failed to get cash balance: %v", err)
	}
	fmt.Printf("✅ Updated cash balance: $%.2f\n", float64(cashBalance.Balance.Value)/100)

	// Demo 10: Summary
	fmt.Println("\n🎉 Demo Completed Successfully!")
	fmt.Println("===================================")
	fmt.Println("The accounting system demonstrated:")
	fmt.Println("✅ Double-entry bookkeeping")
	fmt.Println("✅ Multi-dimensional analytics")
	fmt.Println("✅ Event sourcing & audit trails")
	fmt.Println("✅ Accrual/deferral recognition")
	fmt.Println("✅ Period management")
	fmt.Println("✅ Transaction reversals")
	fmt.Println("✅ Trial balance reporting")
	fmt.Println("✅ Bi-temporal data tracking")
	fmt.Println("\n🚀 Ready for production use!")
}
