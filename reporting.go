package accounting

import (
	"fmt"
	"time"
)

// FinancialStatement represents a financial statement
type FinancialStatement struct {
	Name        string               `json:"name"`
	AsOfDate    time.Time            `json:"as_of_date"`
	FromDate    *time.Time           `json:"from_date,omitempty"` // For P&L and Cash Flow
	Currency    string               `json:"currency"`
	LineItems   []*FinancialLineItem `json:"line_items"`
	TotalAssets *Amount              `json:"total_assets,omitempty"`
	TotalLiabs  *Amount              `json:"total_liabilities,omitempty"`
	TotalEquity *Amount              `json:"total_equity,omitempty"`
	NetIncome   *Amount              `json:"net_income,omitempty"`
}

// FinancialLineItem represents a line item in a financial statement
type FinancialLineItem struct {
	AccountID   string               `json:"account_id"`
	AccountName string               `json:"account_name"`
	AccountType AccountType          `json:"account_type"`
	Amount      *Amount              `json:"amount"`
	Level       int                  `json:"level"` // For hierarchy display
	IsSubtotal  bool                 `json:"is_subtotal"`
	Children    []*FinancialLineItem `json:"children,omitempty"`
}

// CashFlowItem represents an item in cash flow statement
type CashFlowItem struct {
	Description string           `json:"description"`
	Amount      *Amount          `json:"amount"`
	Category    CashFlowCategory `json:"category"`
}

// CashFlowCategory represents categories in cash flow statement
type CashFlowCategory string

const (
	CashFlowOperating CashFlowCategory = "OPERATING"
	CashFlowInvesting CashFlowCategory = "INVESTING"
	CashFlowFinancing CashFlowCategory = "FINANCING"
)

// CashFlowStatement represents a cash flow statement
type CashFlowStatement struct {
	Name                string          `json:"name"`
	FromDate            time.Time       `json:"from_date"`
	ToDate              time.Time       `json:"to_date"`
	Currency            string          `json:"currency"`
	OperatingActivities []*CashFlowItem `json:"operating_activities"`
	InvestingActivities []*CashFlowItem `json:"investing_activities"`
	FinancingActivities []*CashFlowItem `json:"financing_activities"`
	NetCashFlow         *Amount         `json:"net_cash_flow"`
	BeginningCash       *Amount         `json:"beginning_cash"`
	EndingCash          *Amount         `json:"ending_cash"`
}

// ReportingService handles financial statement generation
type ReportingService struct {
	storage  *Storage
	queryAPI *QueryAPI
}

// NewReportingService creates a new reporting service
func NewReportingService(storage *Storage, queryAPI *QueryAPI) *ReportingService {
	return &ReportingService{
		storage:  storage,
		queryAPI: queryAPI,
	}
}

// GenerateBalanceSheet generates a balance sheet as of a specific date
func (rs *ReportingService) GenerateBalanceSheet(asOfDate time.Time, currency string) (*FinancialStatement, error) {
	// Get all account balances
	trialBalance, err := rs.queryAPI.GetTrialBalance(asOfDate, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get trial balance: %w", err)
	}

	bs := &FinancialStatement{
		Name:     "Balance Sheet",
		AsOfDate: asOfDate,
		Currency: currency,
	}

	// Group accounts by type
	assets := make([]*FinancialLineItem, 0)
	liabilities := make([]*FinancialLineItem, 0)
	equity := make([]*FinancialLineItem, 0)

	totalAssets := &Amount{Value: 0, Currency: Currency(currency)}
	totalLiabs := &Amount{Value: 0, Currency: Currency(currency)}
	totalEquity := &Amount{Value: 0, Currency: Currency(currency)}

	for _, balance := range trialBalance {
		lineItem := &FinancialLineItem{
			AccountID:   balance.AccountID,
			AccountName: balance.AccountName,
			AccountType: balance.AccountType,
			Amount:      balance.Balance,
			Level:       1,
		}

		switch balance.AccountType {
		case Asset:
			assets = append(assets, lineItem)
			totalAssets.Value += balance.Balance.Value
		case Liability:
			liabilities = append(liabilities, lineItem)
			totalLiabs.Value += balance.Balance.Value
		case Equity:
			equity = append(equity, lineItem)
			totalEquity.Value += balance.Balance.Value
		}
	}

	// Build line items with subtotals
	bs.LineItems = append(bs.LineItems, &FinancialLineItem{
		AccountName: "ASSETS",
		Level:       0,
		IsSubtotal:  true,
		Children:    assets,
	})

	bs.LineItems = append(bs.LineItems, &FinancialLineItem{
		AccountName: "LIABILITIES",
		Level:       0,
		IsSubtotal:  true,
		Children:    liabilities,
	})

	bs.LineItems = append(bs.LineItems, &FinancialLineItem{
		AccountName: "EQUITY",
		Level:       0,
		IsSubtotal:  true,
		Children:    equity,
	})

	bs.TotalAssets = totalAssets
	bs.TotalLiabs = totalLiabs
	bs.TotalEquity = totalEquity

	return bs, nil
}

// GenerateProfitAndLoss generates a P&L statement for a period
func (rs *ReportingService) GenerateProfitAndLoss(fromDate, toDate time.Time, currency string) (*FinancialStatement, error) {
	// Get all account balances for the period
	trialBalance, err := rs.queryAPI.GetTrialBalance(toDate, []AccountType{Income, Expense})
	if err != nil {
		return nil, fmt.Errorf("failed to get trial balance: %w", err)
	}

	pl := &FinancialStatement{
		Name:     "Profit & Loss Statement",
		AsOfDate: toDate,
		FromDate: &fromDate,
		Currency: currency,
	}

	// Group accounts by type
	revenue := make([]*FinancialLineItem, 0)
	expenses := make([]*FinancialLineItem, 0)

	totalRevenue := &Amount{Value: 0, Currency: Currency(currency)}
	totalExpenses := &Amount{Value: 0, Currency: Currency(currency)}

	for _, balance := range trialBalance {
		// Filter transactions by date range
		periodBalance, err := rs.calculatePeriodBalance(balance.AccountID, fromDate, toDate)
		if err != nil {
			continue // Skip on error
		}

		lineItem := &FinancialLineItem{
			AccountID:   balance.AccountID,
			AccountName: balance.AccountName,
			AccountType: balance.AccountType,
			Amount:      periodBalance,
			Level:       1,
		}

		switch balance.AccountType {
		case Income:
			revenue = append(revenue, lineItem)
			totalRevenue.Value += periodBalance.Value
		case Expense:
			expenses = append(expenses, lineItem)
			totalExpenses.Value += periodBalance.Value
		}
	}

	// Calculate net income
	netIncome := &Amount{
		Value:    totalRevenue.Value - totalExpenses.Value,
		Currency: Currency(currency),
	}

	// Build line items
	pl.LineItems = append(pl.LineItems, &FinancialLineItem{
		AccountName: "REVENUE",
		Level:       0,
		IsSubtotal:  true,
		Amount:      totalRevenue,
		Children:    revenue,
	})

	pl.LineItems = append(pl.LineItems, &FinancialLineItem{
		AccountName: "EXPENSES",
		Level:       0,
		IsSubtotal:  true,
		Amount:      totalExpenses,
		Children:    expenses,
	})

	pl.LineItems = append(pl.LineItems, &FinancialLineItem{
		AccountName: "NET INCOME",
		Level:       0,
		IsSubtotal:  true,
		Amount:      netIncome,
	})

	pl.NetIncome = netIncome

	return pl, nil
}

// GenerateCashFlowStatement generates a cash flow statement for a period
func (rs *ReportingService) GenerateCashFlowStatement(fromDate, toDate time.Time, currency string) (*CashFlowStatement, error) {
	// Get cash account movements
	cashEntries, err := rs.storage.GetEntriesByAccount("cash")
	if err != nil {
		return nil, fmt.Errorf("failed to get cash entries: %w", err)
	}

	cf := &CashFlowStatement{
		Name:     "Cash Flow Statement",
		FromDate: fromDate,
		ToDate:   toDate,
		Currency: currency,
	}

	// Get beginning cash balance
	beginningCash, err := rs.queryAPI.GetAccountBalance("cash", fromDate)
	if err == nil {
		cf.BeginningCash = beginningCash.Balance
	} else {
		cf.BeginningCash = &Amount{Value: 0, Currency: Currency(currency)}
	}

	// Categorize cash flows
	operatingCashFlow := int64(0)
	investingCashFlow := int64(0)
	financingCashFlow := int64(0)

	for _, entry := range cashEntries {
		// Get transaction to check date and get contra accounts
		txn, err := rs.storage.GetTransaction(entry.TransactionID)
		if err != nil {
			continue
		}

		// Filter by date range
		if txn.ValidTime.Before(fromDate) || txn.ValidTime.After(toDate) || txn.Status != Posted {
			continue
		}

		// Determine cash flow category based on contra accounts
		category := rs.categorizeTransactionForCashFlow(txn, entry.AccountID)

		amount := entry.Amount.Value
		if entry.Type == Credit {
			amount = -amount // Cash going out
		}

		item := &CashFlowItem{
			Description: txn.Description,
			Amount:      &Amount{Value: amount, Currency: entry.Amount.Currency},
			Category:    category,
		}

		switch category {
		case CashFlowOperating:
			cf.OperatingActivities = append(cf.OperatingActivities, item)
			operatingCashFlow += amount
		case CashFlowInvesting:
			cf.InvestingActivities = append(cf.InvestingActivities, item)
			investingCashFlow += amount
		case CashFlowFinancing:
			cf.FinancingActivities = append(cf.FinancingActivities, item)
			financingCashFlow += amount
		}
	}

	// Calculate net cash flow
	netCashFlow := operatingCashFlow + investingCashFlow + financingCashFlow
	cf.NetCashFlow = &Amount{Value: netCashFlow, Currency: Currency(currency)}

	// Calculate ending cash
	cf.EndingCash = &Amount{
		Value:    cf.BeginningCash.Value + netCashFlow,
		Currency: Currency(currency),
	}

	return cf, nil
}

// calculatePeriodBalance calculates account balance for a specific period
func (rs *ReportingService) calculatePeriodBalance(accountID string, fromDate, toDate time.Time) (*Amount, error) {
	entries, err := rs.storage.GetEntriesByAccount(accountID)
	if err != nil {
		return nil, err
	}

	account, err := rs.storage.GetAccount(accountID)
	if err != nil {
		return nil, err
	}

	balance := &Amount{Value: 0, Currency: account.Currency}

	for _, entry := range entries {
		txn, err := rs.storage.GetTransaction(entry.TransactionID)
		if err != nil {
			continue
		}

		// Filter by date range
		if txn.ValidTime.Before(fromDate) || txn.ValidTime.After(toDate) || txn.Status != Posted {
			continue
		}

		// Calculate balance based on account type and entry type
		multiplier := 1
		if (account.Type == Income || account.Type == Liability || account.Type == Equity) && entry.Type == Credit {
			multiplier = 1
		} else if (account.Type == Asset || account.Type == Expense) && entry.Type == Debit {
			multiplier = 1
		} else {
			multiplier = -1
		}

		balance.Value += entry.Amount.Value * int64(multiplier)
	}

	return balance, nil
}

// categorizeTransactionForCashFlow categorizes a transaction for cash flow purposes
func (rs *ReportingService) categorizeTransactionForCashFlow(txn *Transaction, cashAccountID string) CashFlowCategory {
	// Look at contra accounts to determine category
	for _, entry := range txn.Entries {
		if entry.AccountID == cashAccountID {
			continue // Skip the cash account itself
		}

		// Get account to determine its type
		account, err := rs.storage.GetAccount(entry.AccountID)
		if err != nil {
			continue
		}

		// Categorize based on account type and common patterns
		switch account.Type {
		case Income, Expense:
			return CashFlowOperating
		case Asset:
			if account.ID == "inventory" || account.ID == "equipment" || account.ID == "property" {
				return CashFlowInvesting
			}
			return CashFlowOperating
		case Liability:
			if account.ID == "notes_payable" || account.ID == "bonds_payable" {
				return CashFlowFinancing
			}
			return CashFlowOperating
		case Equity:
			return CashFlowFinancing
		}
	}

	// Default to operating if we can't determine
	return CashFlowOperating
}

// FormatFinancialStatement formats a financial statement for display
func (rs *ReportingService) FormatFinancialStatement(statement *FinancialStatement) string {
	var output string

	output += fmt.Sprintf("\n%s\n", statement.Name)
	if statement.FromDate != nil {
		output += fmt.Sprintf("Period: %s to %s\n",
			statement.FromDate.Format("2006-01-02"),
			statement.AsOfDate.Format("2006-01-02"))
	} else {
		output += fmt.Sprintf("As of: %s\n", statement.AsOfDate.Format("2006-01-02"))
	}
	output += fmt.Sprintf("Currency: %s\n", statement.Currency)
	output += "==========================================\n"

	for _, lineItem := range statement.LineItems {
		output += rs.formatLineItem(lineItem, 0)
	}

	if statement.NetIncome != nil {
		output += fmt.Sprintf("\nNET INCOME: $%.2f\n", float64(statement.NetIncome.Value)/100)
	}

	if statement.TotalAssets != nil {
		output += fmt.Sprintf("\nTOTAL ASSETS: $%.2f\n", float64(statement.TotalAssets.Value)/100)
		output += fmt.Sprintf("TOTAL LIAB + EQUITY: $%.2f\n",
			float64(statement.TotalLiabs.Value+statement.TotalEquity.Value)/100)
	}

	return output
}

// formatLineItem formats a single line item
func (rs *ReportingService) formatLineItem(item *FinancialLineItem, indent int) string {
	var output string
	indentStr := ""
	for i := 0; i < indent; i++ {
		indentStr += "  "
	}

	if item.IsSubtotal {
		output += fmt.Sprintf("%s%s\n", indentStr, item.AccountName)
		output += fmt.Sprintf("%s%s\n", indentStr, "--------------------")

		// Show children
		for _, child := range item.Children {
			output += rs.formatLineItem(child, indent+1)
		}

		if item.Amount != nil {
			output += fmt.Sprintf("%s%s: $%.2f\n", indentStr, "TOTAL", float64(item.Amount.Value)/100)
		}
		output += "\n"
	} else {
		output += fmt.Sprintf("%s%-20s $%8.2f\n",
			indentStr,
			item.AccountName,
			float64(item.Amount.Value)/100)
	}

	return output
}

// FormatCashFlowStatement formats a cash flow statement for display
func (rs *ReportingService) FormatCashFlowStatement(cf *CashFlowStatement) string {
	var output string

	output += fmt.Sprintf("\n%s\n", cf.Name)
	output += fmt.Sprintf("Period: %s to %s\n",
		cf.FromDate.Format("2006-01-02"),
		cf.ToDate.Format("2006-01-02"))
	output += fmt.Sprintf("Currency: %s\n", cf.Currency)
	output += "==========================================\n"

	// Operating Activities
	output += "OPERATING ACTIVITIES:\n"
	operatingTotal := int64(0)
	for _, item := range cf.OperatingActivities {
		output += fmt.Sprintf("  %-30s $%8.2f\n", item.Description, float64(item.Amount.Value)/100)
		operatingTotal += item.Amount.Value
	}
	output += fmt.Sprintf("  Net Cash from Operations: $%8.2f\n\n", float64(operatingTotal)/100)

	// Investing Activities
	output += "INVESTING ACTIVITIES:\n"
	investingTotal := int64(0)
	for _, item := range cf.InvestingActivities {
		output += fmt.Sprintf("  %-30s $%8.2f\n", item.Description, float64(item.Amount.Value)/100)
		investingTotal += item.Amount.Value
	}
	output += fmt.Sprintf("  Net Cash from Investing: $%8.2f\n\n", float64(investingTotal)/100)

	// Financing Activities
	output += "FINANCING ACTIVITIES:\n"
	financingTotal := int64(0)
	for _, item := range cf.FinancingActivities {
		output += fmt.Sprintf("  %-30s $%8.2f\n", item.Description, float64(item.Amount.Value)/100)
		financingTotal += item.Amount.Value
	}
	output += fmt.Sprintf("  Net Cash from Financing: $%8.2f\n\n", float64(financingTotal)/100)

	// Summary
	output += "CASH FLOW SUMMARY:\n"
	output += fmt.Sprintf("  Beginning Cash:      $%8.2f\n", float64(cf.BeginningCash.Value)/100)
	output += fmt.Sprintf("  Net Cash Flow:       $%8.2f\n", float64(cf.NetCashFlow.Value)/100)
	output += fmt.Sprintf("  Ending Cash:         $%8.2f\n", float64(cf.EndingCash.Value)/100)

	return output
}
