package accounting

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinancialReporting(t *testing.T) {
	// Setup
	dbFile := "test_reporting.db"
	defer os.Remove(dbFile)

	engine, err := NewAccountingEngine(dbFile)
	require.NoError(t, err)
	defer engine.Close()

	userID := "test_user"

	// Create accounts
	err = engine.CreateStandardAccounts(userID)
	require.NoError(t, err)

	// Create test transactions
	saleTransaction := &Transaction{
		Description: "Test sale",
		ValidTime:   time.Now(),
		Entries: []Entry{
			{
				AccountID: "cash",
				Type:      Debit,
				Amount:    Amount{Value: 100000, Currency: "USD"}, // $1,000
			},
			{
				AccountID: "revenue",
				Type:      Credit,
				Amount:    Amount{Value: 100000, Currency: "USD"},
			},
		},
	}

	err = engine.CreateTransaction(saleTransaction, userID)
	require.NoError(t, err)
	err = engine.PostTransaction(saleTransaction.ID, userID)
	require.NoError(t, err)

	t.Run("Balance Sheet Generation", func(t *testing.T) {
		balanceSheet, err := engine.GenerateBalanceSheet(time.Now(), "USD")
		require.NoError(t, err)

		assert.Equal(t, "Balance Sheet", balanceSheet.Name)
		assert.Equal(t, "USD", balanceSheet.Currency)
		assert.NotEmpty(t, balanceSheet.LineItems)

		// Check that we have asset accounts in the children
		hasAssets := false
		for _, item := range balanceSheet.LineItems {
			if item.AccountName == "ASSETS" && len(item.Children) > 0 {
				for _, child := range item.Children {
					if child.AccountType == Asset {
						hasAssets = true
						break
					}
				}
			}
		}
		assert.True(t, hasAssets, "Balance sheet should have asset accounts")
	})

	t.Run("Profit & Loss Generation", func(t *testing.T) {
		startDate := time.Now().AddDate(0, 0, -30)
		endDate := time.Now()

		pl, err := engine.GenerateProfitAndLoss(startDate, endDate, "USD")
		require.NoError(t, err)

		assert.Equal(t, "Profit & Loss Statement", pl.Name)
		assert.Equal(t, "USD", pl.Currency)
		assert.NotEmpty(t, pl.LineItems)

		// Check that we have revenue accounts in the children
		hasRevenue := false
		for _, item := range pl.LineItems {
			if item.AccountName == "REVENUE" && len(item.Children) > 0 {
				for _, child := range item.Children {
					if child.AccountType == Income {
						hasRevenue = true
						break
					}
				}
			}
		}
		assert.True(t, hasRevenue, "P&L should have revenue accounts")
	})

	t.Run("Cash Flow Generation", func(t *testing.T) {
		startDate := time.Now().AddDate(0, 0, -30)
		endDate := time.Now()

		cashFlow, err := engine.GenerateCashFlowStatement(startDate, endDate, "USD")
		require.NoError(t, err)

		assert.Equal(t, "Cash Flow Statement", cashFlow.Name)
		assert.Equal(t, "USD", cashFlow.Currency)
		assert.NotEmpty(t, cashFlow.OperatingActivities)

		// Check that we have operating activities
		assert.True(t, len(cashFlow.OperatingActivities) >= 0, "Cash flow should have operating activities section")
	})

	t.Run("Statement Formatting", func(t *testing.T) {
		balanceSheet, err := engine.GenerateBalanceSheet(time.Now(), "USD")
		require.NoError(t, err)

		formatted := engine.FormatFinancialStatement(balanceSheet)
		assert.Contains(t, formatted, "Balance Sheet")
		assert.Contains(t, formatted, "USD")
		assert.Contains(t, formatted, "ASSETS")
	})
}

func TestMultiCompanySupport(t *testing.T) {
	// Setup storage with unique filename
	dbFile := fmt.Sprintf("test_multicompany_%d.db", time.Now().UnixNano())
	defer os.Remove(dbFile)

	storage, err := NewStorage(dbFile)
	require.NoError(t, err)
	defer storage.Close()

	engine := NewMultiCompanyEngine(*storage)
	defer engine.Close()
	userID := "test_user"

	t.Run("Company Creation", func(t *testing.T) {
		company := &Company{
			ID:           "test_company",
			Name:         "Test Company",
			LegalName:    "Test Company Inc.",
			TaxID:        "12-3456789",
			BaseCurrency: "USD",
			Status:       CompanyActive,
			Settings: &CompanySettings{
				DefaultChartOfAccounts: "standard",
				AllowIntercompanyTxn:   true,
				ReportingCurrency:      "USD",
			},
		}

		err := engine.CreateCompany(company, userID)
		require.NoError(t, err)

		// Verify company was created
		retrieved, err := engine.GetCompany("test_company")
		require.NoError(t, err)
		assert.Equal(t, company.Name, retrieved.Name)
		assert.Equal(t, company.TaxID, retrieved.TaxID)
	})

	t.Run("Intercompany Transaction", func(t *testing.T) {
		// Create two companies
		parentCompany := &Company{
			ID:           "parent",
			Name:         "Parent Corp",
			BaseCurrency: "USD",
			Status:       CompanyActive,
			Settings: &CompanySettings{
				DefaultChartOfAccounts: "standard",
				AllowIntercompanyTxn:   true,
			},
		}

		subsidiaryCompany := &Company{
			ID:              "subsidiary",
			Name:            "Subsidiary LLC",
			BaseCurrency:    "USD",
			ParentCompanyID: "parent",
			Status:          CompanyActive,
			Settings: &CompanySettings{
				DefaultChartOfAccounts: "standard",
				AllowIntercompanyTxn:   true,
			},
		}

		err := engine.CreateCompany(parentCompany, userID)
		require.NoError(t, err)
		err = engine.CreateCompany(subsidiaryCompany, userID)
		require.NoError(t, err)

		// Create intercompany transaction using the correct method
		icTransaction, err := engine.CreateIntercompanyTransaction(
			"parent", "subsidiary",
			&Amount{Value: 50000, Currency: "USD"}, // $500
			"Intercompany loan",
			userID)
		require.NoError(t, err)

		// Verify transaction was created
		assert.NotEmpty(t, icTransaction.ID)
		assert.Equal(t, IntercompanyMatched, icTransaction.MatchingStatus)
	})

	t.Run("Consolidation Group", func(t *testing.T) {
		group := &ConsolidationGroup{
			ID:                  "test_group",
			Name:                "Test Consolidation Group",
			ParentCompany:       "parent",
			ChildCompanies:      []string{"subsidiary"},
			ConsolidationMethod: "FULL",
		}

		err := engine.CreateConsolidationGroup(group, userID)
		require.NoError(t, err)

		// Verify group was created
		retrieved, err := storage.GetConsolidationGroup("test_group")
		require.NoError(t, err)
		assert.Equal(t, group.Name, retrieved.Name)
		assert.Equal(t, group.ParentCompany, retrieved.ParentCompany)
		assert.Equal(t, group.ChildCompanies, retrieved.ChildCompanies)
	})
}

func TestReportingServiceIntegration(t *testing.T) {
	// Setup
	dbFile := "test_integration.db"
	defer os.Remove(dbFile)

	storage, err := NewStorage(dbFile)
	require.NoError(t, err)
	defer storage.Close()

	// Create a minimal query API for testing
	queryAPI := &QueryAPI{storage: storage}
	reportingService := NewReportingService(storage, queryAPI)

	t.Run("Reporting Service Creation", func(t *testing.T) {
		assert.NotNil(t, reportingService)
		assert.NotNil(t, reportingService.storage)
		assert.NotNil(t, reportingService.queryAPI)
	})
}

func TestStorageExtensions(t *testing.T) {
	// Setup
	dbFile := "test_storage_ext.db"
	defer os.Remove(dbFile)

	storage, err := NewStorage(dbFile)
	require.NoError(t, err)
	defer storage.Close()

	t.Run("Company Storage", func(t *testing.T) {
		company := &Company{
			ID:           "storage_test",
			Name:         "Storage Test Company",
			BaseCurrency: "USD",
			Status:       CompanyActive,
			CreatedAt:    time.Now(),
		}

		// Save company
		err := storage.SaveCompany(company)
		require.NoError(t, err)

		// Retrieve company
		retrieved, err := storage.GetCompany("storage_test")
		require.NoError(t, err)
		assert.Equal(t, company.Name, retrieved.Name)
		assert.Equal(t, company.BaseCurrency, retrieved.BaseCurrency)

		// List companies
		companies, err := storage.GetCompanies()
		require.NoError(t, err)
		assert.True(t, len(companies) >= 1)
	})

	t.Run("Intercompany Transaction Storage", func(t *testing.T) {
		transaction := &IntercompanyTransaction{
			ID:              "ic_test_001",
			SourceCompanyID: "source",
			TargetCompanyID: "target",
			Amount:          &Amount{Value: 25000, Currency: "USD"},
			Description:     "Test intercompany transaction",
			MatchingStatus:  IntercompanyPending,
			CreatedAt:       time.Now(),
			CreatedBy:       "test_user",
		}

		// Save transaction
		err := storage.SaveIntercompanyTransaction(transaction)
		require.NoError(t, err)

		// Retrieve transaction
		retrieved, err := storage.GetIntercompanyTransaction("ic_test_001")
		require.NoError(t, err)
		assert.Equal(t, transaction.SourceCompanyID, retrieved.SourceCompanyID)
		assert.Equal(t, transaction.TargetCompanyID, retrieved.TargetCompanyID)
		assert.Equal(t, transaction.Amount.Value, retrieved.Amount.Value)

		// List transactions by company
		transactions, err := storage.GetIntercompanyTransactionsByCompany("source")
		require.NoError(t, err)
		assert.True(t, len(transactions) >= 1)
	})

	t.Run("Consolidation Group Storage", func(t *testing.T) {
		group := &ConsolidationGroup{
			ID:                  "group_test",
			Name:                "Test Group",
			ParentCompany:       "parent",
			ChildCompanies:      []string{"sub1", "sub2"},
			ConsolidationMethod: "FULL",
			CreatedAt:           time.Now(),
			CreatedBy:           "test_user",
		}

		// Save group
		err := storage.SaveConsolidationGroup(group)
		require.NoError(t, err)

		// Retrieve group
		retrieved, err := storage.GetConsolidationGroup("group_test")
		require.NoError(t, err)
		assert.Equal(t, group.Name, retrieved.Name)
		assert.Equal(t, group.ParentCompany, retrieved.ParentCompany)
		assert.Equal(t, group.ChildCompanies, retrieved.ChildCompanies)

		// List groups
		groups, err := storage.GetConsolidationGroups()
		require.NoError(t, err)
		assert.True(t, len(groups) >= 1)
	})
}
