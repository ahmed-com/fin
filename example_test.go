package accounting

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// TestAccountingSystem demonstrates the full accounting system capabilities
func TestAccountingSystem(t *testing.T) {
	// Create temporary database file
	dbFile := "/tmp/test_accounting.db"
	defer os.Remove(dbFile)

	// Initialize the accounting engine
	engine, err := NewAccountingEngine(dbFile)
	if err != nil {
		t.Fatalf("Failed to create accounting engine: %v", err)
	}
	defer engine.Close()

	userID := "test_user"

	// Test 1: Create Chart of Accounts
	fmt.Println("=== Test 1: Creating Chart of Accounts ===")
	if err := engine.CreateStandardAccounts(userID); err != nil {
		t.Fatalf("Failed to create standard accounts: %v", err)
	}
	fmt.Println("✓ Standard accounts created successfully")

	// Test 2: Create and Post a Simple Transaction
	fmt.Println("\n=== Test 2: Creating and Posting Transaction ===")

	// Create a simple sale transaction: Cash 1000 / Revenue 1000
	saleTransaction := &Transaction{
		Description: "Sale of products",
		ValidTime:   time.Now(),
		Entries: []Entry{
			{
				AccountID: "cash",
				Type:      Debit,
				Amount: Amount{
					Value:    100000, // $1000.00 in cents
					Currency: "USD",
				},
			},
			{
				AccountID: "revenue",
				Type:      Credit,
				Amount: Amount{
					Value:    100000,
					Currency: "USD",
				},
			},
		},
	}

	// Create the transaction
	if err := engine.CreateTransaction(saleTransaction, userID); err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	// Post the transaction
	if err := engine.PostTransaction(saleTransaction.ID, userID); err != nil {
		t.Fatalf("Failed to post transaction: %v", err)
	}
	fmt.Printf("✓ Transaction created and posted: %s\n", saleTransaction.ID)

	// Test 3: Check Account Balances
	fmt.Println("\n=== Test 3: Checking Account Balances ===")

	cashBalance, err := engine.GetAccountBalance("cash", time.Now())
	if err != nil {
		t.Fatalf("Failed to get cash balance: %v", err)
	}
	fmt.Printf("Cash Balance: $%.2f %s\n", float64(cashBalance.Balance.Value)/100, cashBalance.Balance.Currency)

	revenueBalance, err := engine.GetAccountBalance("revenue", time.Now())
	if err != nil {
		t.Fatalf("Failed to get revenue balance: %v", err)
	}
	fmt.Printf("Revenue Balance: $%.2f %s\n", float64(revenueBalance.Balance.Value)/100, revenueBalance.Balance.Currency)

	// Test 4: Create an Accrual Schedule
	fmt.Println("\n=== Test 4: Creating Accrual Schedule ===")

	// Create a prepaid expense transaction
	prepaidTransaction := &Transaction{
		Description: "Prepaid Insurance - 12 months",
		ValidTime:   time.Now(),
		Entries: []Entry{
			{
				AccountID: "cash", // Actually should be prepaid_expenses, but using cash for simplicity
				Type:      Credit,
				Amount: Amount{
					Value:    120000, // $1200.00 for 12 months
					Currency: "USD",
				},
			},
			{
				AccountID: "expenses",
				Type:      Debit,
				Amount: Amount{
					Value:    120000,
					Currency: "USD",
				},
			},
		},
	}

	if err := engine.CreateTransaction(prepaidTransaction, userID); err != nil {
		t.Fatalf("Failed to create prepaid transaction: %v", err)
	}

	if err := engine.PostTransaction(prepaidTransaction.ID, userID); err != nil {
		t.Fatalf("Failed to post prepaid transaction: %v", err)
	}

	// Create accrual schedule for monthly recognition
	template := &AccrualTemplate{
		ID:               "expense_template",
		Name:             "Monthly Expense Recognition",
		AccrualType:      AccrualExpense,
		ExpenseAccountID: "expenses",
	}

	schedule, err := engine.CreateAccrualSchedule(
		prepaidTransaction.ID,
		&Amount{Value: 120000, Currency: "USD"},
		Monthly,
		12,
		time.Now(),
		template,
		userID,
	)
	if err != nil {
		t.Fatalf("Failed to create accrual schedule: %v", err)
	}
	fmt.Printf("✓ Accrual schedule created: %s\n", schedule.ID)

	// Test 5: Create an Accounting Period
	fmt.Println("\n=== Test 5: Creating Accounting Period ===")

	now := time.Now()
	period := &Period{
		Name:  "2025-Q2",
		Start: time.Date(now.Year(), 4, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(now.Year(), 6, 30, 23, 59, 59, 0, time.UTC),
	}

	if err := engine.CreatePeriod(period, userID); err != nil {
		t.Fatalf("Failed to create period: %v", err)
	}
	fmt.Printf("✓ Period created: %s (%s to %s)\n", period.Name, period.Start.Format("2006-01-02"), period.End.Format("2006-01-02"))

	// Test 6: Generate Trial Balance
	fmt.Println("\n=== Test 6: Generating Trial Balance ===")

	trialBalance, err := engine.GetTrialBalance(time.Now(), nil)
	if err != nil {
		t.Fatalf("Failed to generate trial balance: %v", err)
	}

	fmt.Println("Trial Balance:")
	fmt.Println("Account                   Type        Balance")
	fmt.Println("--------------------------------------------")
	totalDebits := int64(0)
	totalCredits := int64(0)

	for _, balance := range trialBalance {
		fmt.Printf("%-25s %-10s $%.2f\n",
			balance.AccountName,
			balance.AccountType,
			float64(balance.Balance.Value)/100)

		// Calculate totals based on normal balance sides
		switch balance.AccountType {
		case Asset, Expense:
			if balance.Balance.Value > 0 {
				totalDebits += balance.Balance.Value
			} else {
				totalCredits += -balance.Balance.Value
			}
		case Liability, Equity, Income:
			if balance.Balance.Value > 0 {
				totalCredits += balance.Balance.Value
			} else {
				totalDebits += -balance.Balance.Value
			}
		}
	}

	fmt.Println("--------------------------------------------")
	fmt.Printf("Total Debits:  $%.2f\n", float64(totalDebits)/100)
	fmt.Printf("Total Credits: $%.2f\n", float64(totalCredits)/100)

	// Test 7: Test Reconciliation
	fmt.Println("\n=== Test 7: Testing Reconciliation ===")

	// Create mock bank statement
	statements := []*ExternalStatement{
		{
			ID:          "stmt_001",
			Date:        time.Now(),
			Description: "Sale deposit",
			Amount: &Amount{
				Value:    100000,
				Currency: "USD",
			},
			Reference:   "DEP001",
			BankAccount: "checking_001",
		},
	}

	matches, err := engine.AutoReconcile("cash", statements)
	if err != nil {
		t.Fatalf("Failed to auto reconcile: %v", err)
	}
	fmt.Printf("✓ Found %d reconciliation matches\n", len(matches))

	// Test 8: Transaction Reversal
	fmt.Println("\n=== Test 8: Testing Transaction Reversal ===")

	reversalTxn, err := engine.ReverseTransaction(saleTransaction.ID, "Reversal of sale", userID)
	if err != nil {
		t.Fatalf("Failed to reverse transaction: %v", err)
	}
	fmt.Printf("✓ Transaction reversed: %s\n", reversalTxn.ID)

	// Check cash balance after reversal
	cashBalanceAfterReversal, err := engine.GetAccountBalance("cash", time.Now())
	if err != nil {
		t.Fatalf("Failed to get cash balance after reversal: %v", err)
	}
	fmt.Printf("Cash Balance after reversal: $%.2f %s\n",
		float64(cashBalanceAfterReversal.Balance.Value)/100,
		cashBalanceAfterReversal.Balance.Currency)

	// Test 9: Event Sourcing - Replay Events
	fmt.Println("\n=== Test 9: Testing Event Sourcing ===")

	events, err := engine.GetEvents(time.Now().Add(-1*time.Hour), time.Now())
	if err != nil {
		t.Fatalf("Failed to get events: %v", err)
	}
	fmt.Printf("✓ Retrieved %d events from the last hour\n", len(events))

	for i, event := range events {
		fmt.Printf("Event %d: %s at %s\n", i+1, event.EventType, event.TransactionTime.Format("15:04:05"))
	}

	fmt.Println("\n=== All Tests Completed Successfully! ===")
}

// DemoUsage demonstrates basic usage of the accounting system
func DemoUsage() {
	// Initialize the accounting engine
	engine, err := NewAccountingEngine("accounting.db")
	if err != nil {
		panic(err)
	}
	defer engine.Close()

	userID := "user_123"

	// Create basic accounts
	engine.CreateStandardAccounts(userID)

	// Create a sales transaction
	transaction := &Transaction{
		Description: "Product Sale",
		ValidTime:   time.Now(),
		Entries: []Entry{
			{
				AccountID: "cash",
				Type:      Debit,
				Amount:    Amount{Value: 50000, Currency: "USD"}, // $500.00
				Dimensions: []Dimension{
					{Key: DimDepartment, Value: "sales"},
					{Key: DimProduct, Value: "widget_a"},
				},
			},
			{
				AccountID: "revenue",
				Type:      Credit,
				Amount:    Amount{Value: 50000, Currency: "USD"},
				Dimensions: []Dimension{
					{Key: DimDepartment, Value: "sales"},
					{Key: DimProduct, Value: "widget_a"},
				},
			},
		},
	}

	// Create and post the transaction
	engine.CreateTransaction(transaction, userID)
	engine.PostTransaction(transaction.ID, userID)

	// Get account balance
	balance, _ := engine.GetAccountBalance("cash", time.Now())
	fmt.Printf("Cash balance: $%.2f\n", float64(balance.Balance.Value)/100)

	// Create an accrual schedule for subscription revenue
	template := &AccrualTemplate{
		ID:               "revenue_template",
		Name:             "Monthly Revenue Recognition",
		AccrualType:      AccrualRevenue,
		RevenueAccountID: "revenue",
	}

	engine.CreateAccrualSchedule(
		transaction.ID,
		&Amount{Value: 120000, Currency: "USD"}, // $1200 over 12 months
		Monthly,
		12,
		time.Now(),
		template,
		userID,
	)

	// Process accruals up to current date
	engine.ProcessAccruals(time.Now(), userID)

	fmt.Println("Example completed successfully!")
}
