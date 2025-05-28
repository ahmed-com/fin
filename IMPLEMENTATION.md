# Accounting System Implementation

This is a comprehensive implementation of an accounting software system in Go that supports advanced features like event sourcing, double-entry bookkeeping, multi-currency transactions, accrual/deferral recognition, reconciliation, financial reporting, and multi-company consolidation.

## Features Implemented

### ✅ Core Accounting Features

- **Double-Entry Bookkeeping**: Every transaction maintains the fundamental accounting equation (Assets = Liabilities + Equity)
- **Chart of Accounts**: Hierarchical account structure with configurable account types
- **Multi-Currency Support**: Handle transactions in different currencies with exchange rate tracking
- **Bi-Temporal Data**: Track both valid time (business date) and transaction time (when recorded)

### ✅ Advanced Features

- **Event Sourcing**: Immutable event log with ability to replay and rebuild state
- **Multidimensional Analytics**: Tag transactions with dimensions for OLAP-style reporting
- **Period Management**: Soft and hard period closing for compliance
- **Accrual & Deferral Recognition**: Automated recognition schedules for revenue and expenses
- **Bank Reconciliation**: Automated matching of bank statements with internal transactions
- **Audit Trail**: Complete audit trail with user tracking and timestamps

### ✅ Financial Reporting System

- **Balance Sheet Generation**: Hierarchical balance sheet with proper asset, liability, and equity categorization
- **Profit & Loss Statements**: Period-based P&L with revenue and expense categorization
- **Cash Flow Statements**: Operating, investing, and financing activity categorization
- **Financial Statement Formatting**: Professional formatted output for display and reporting
- **Period-Based Calculations**: Accurate balance calculations for any date or period

### ✅ Multi-Company Support

- **Company Entity Management**: Create and manage multiple business entities
- **Intercompany Transactions**: Automated processing with contra entries
- **Consolidation Groups**: Define and manage consolidation hierarchies
- **Elimination Rules**: Automatic elimination of intercompany balances
- **Consolidated Reporting**: Generate consolidated trial balances and reports
- **Multi-Entity Coordination**: Coordinated accounting across multiple companies

### ✅ Regulatory & Tax Compliance

- **Compliance Rule Framework**: Configurable rules for GAAP, IFRS, and SOX compliance
- **Automated Rule Validation**: Real-time transaction validation against compliance rules
- **Tax Calculation Engine**: Multi-jurisdiction tax calculation with configurable rules
- **Tax Return Generation**: Automated tax return calculation and filing preparation
- **Violation Tracking**: Comprehensive compliance violation detection and resolution
- **Audit Trail Integration**: Full audit trails for regulatory compliance reporting
- **Multi-Jurisdiction Support**: US Federal, EU VAT, UK VAT, Canada GST, and more

### ✅ Technical Features

- **BoltDB Storage**: Embedded database for persistence with multi-company extensions
- **Transaction Validation**: Balance checking and business rule validation
- **Query API**: Flexible querying with filters and aggregations
- **Trial Balance**: Generate trial balance reports
- **Account Hierarchy**: Build and query account hierarchies

## Architecture

```
┌─────────────────┐
│ AccountingEngine│  ← Main entry point with reporting integration
├─────────────────┤
│ ReportingService│  ← Financial reporting (Balance Sheet, P&L, Cash Flow)
│ MultiCompanyEng │  ← Multi-company and consolidation management
│ ComplianceService│ ← Regulatory compliance and tax calculation
│ PostingEngine   │  ← Transaction posting and validation
│ QueryAPI        │  ← Data retrieval and reporting
│ ReconcileSvc    │  ← Bank reconciliation
│ AccrualService  │  ← Accrual/deferral recognition
├─────────────────┤
│ EventStore      │  ← Event sourcing
│ EventProcessor  │  ← Event processing
├─────────────────┤
│ Storage         │  ← BoltDB persistence layer with multi-company support
└─────────────────┘
```

## Usage Example

```go
package main

import (
    "fmt"
    "time"
    "accounting"
)

func main() {
    // Initialize the accounting engine
    engine, err := accounting.NewAccountingEngine("accounting.db")
    if err != nil {
        panic(err)
    }
    defer engine.Close()

    userID := "user_123"

    // Create standard chart of accounts
    engine.CreateStandardAccounts(userID)

    // Create a sales transaction
    transaction := &accounting.Transaction{
        Description: "Product Sale",
        ValidTime:   time.Now(),
        Entries: []accounting.Entry{
            {
                AccountID: "cash",
                Type:      accounting.Debit,
                Amount:    accounting.Amount{Value: 100000, Currency: "USD"}, // $1000.00
            },
            {
                AccountID: "revenue",
                Type:      accounting.Credit,
                Amount:    accounting.Amount{Value: 100000, Currency: "USD"},
            },
        },
    }

    // Create and post the transaction
    engine.CreateTransaction(transaction, userID)
    engine.PostTransaction(transaction.ID, userID)

    // Get account balance
    balance, _ := engine.GetAccountBalance("cash", time.Now())
    fmt.Printf("Cash balance: $%.2f\n", float64(balance.Balance.Value)/100)
}
```

## Running Tests

```bash
go test -v
```

## Key Data Structures

### Account
```go
type Account struct {
    ID         string      `json:"id"`
    ParentID   string      `json:"parent_id,omitempty"`
    Code       string      `json:"code"`
    Name       string      `json:"name"`
    Type       AccountType `json:"type"`
    Currency   Currency    `json:"currency,omitempty"`
    Dimensions []Dimension `json:"dimensions,omitempty"`
}
```

### Transaction
```go
type Transaction struct {
    ID              string            `json:"id"`
    Description     string            `json:"description,omitempty"`
    ValidTime       time.Time         `json:"valid_time"`
    TransactionTime time.Time         `json:"transaction_time"`
    Status          TransactionStatus `json:"status"`
    Entries         []Entry           `json:"entries"`
}
```

### Entry
```go
type Entry struct {
    ID            string      `json:"id"`
    TransactionID string      `json:"transaction_id"`
    AccountID     string      `json:"account_id"`
    Type          EntryType   `json:"type"`
    Amount        Amount      `json:"amount"`
    Dimensions    []Dimension `json:"dimensions,omitempty"`
}
```

## Supported Operations

### Account Management
- `CreateAccount()` - Create new accounts
- `GetAccountBalance()` - Get current account balance
- `GetAccountHierarchy()` - Get account tree structure

### Transaction Management
- `CreateTransaction()` - Create new transaction
- `PostTransaction()` - Post transaction to ledger
- `ReverseTransaction()` - Create reversing entry

### Reporting
- `GetTrialBalance()` - Generate trial balance
- `GetDimensionRollup()` - Get aggregated analytics
- `GenerateBalanceSheet()` - Generate balance sheet
- `GenerateProfitAndLoss()` - Generate P&L statement
- `GenerateCashFlowStatement()` - Generate cash flow statement
- `FormatFinancialStatement()` - Format statements for display

### Multi-Company Operations
- `CreateCompany()` - Create new business entity
- `ProcessIntercompanyTransaction()` - Process intercompany transactions
- `CreateConsolidationGroup()` - Create consolidation group
- `GenerateConsolidatedTrialBalance()` - Generate consolidated reporting

### Reconciliation
- `AutoReconcile()` - Automatic bank reconciliation
- `ConfirmReconciliation()` - Confirm reconciliation matches

### Accrual Management
- `CreateAccrualSchedule()` - Create recognition schedule
- `ProcessAccruals()` - Process pending recognitions

### Event Sourcing
- `GetEvents()` - Retrieve event history
- `ReplayEvents()` - Replay events to rebuild state

## Advanced Usage Examples

### Financial Reporting
```go
// Generate Balance Sheet
balanceSheet, err := engine.GenerateBalanceSheet(time.Now(), "USD")
if err != nil {
    log.Fatal(err)
}
fmt.Println(engine.FormatFinancialStatement(balanceSheet))

// Generate Profit & Loss
startDate := time.Now().AddDate(0, -1, 0) // Last month
endDate := time.Now()
pl, err := engine.GenerateProfitAndLoss(startDate, endDate, "USD")
if err != nil {
    log.Fatal(err)
}
fmt.Println(engine.FormatFinancialStatement(pl))

// Generate Cash Flow Statement
cashFlow, err := engine.GenerateCashFlowStatement(startDate, endDate, "USD")
if err != nil {
    log.Fatal(err)
}
fmt.Println(engine.FormatCashFlowStatement(cashFlow))
```

### Multi-Company Operations
```go
// Create storage and multi-company engine
storage, err := accounting.NewStorage("multicompany.db")
if err != nil {
    log.Fatal(err)
}
defer storage.Close()

multiEngine := accounting.NewMultiCompanyEngine(*storage)

// Create parent company
parentCompany := &accounting.Company{
    ID:           "parent_corp",
    Name:         "Parent Corporation",
    BaseCurrency: "USD",
    Status:       accounting.CompanyActive,
    Settings: &accounting.CompanySettings{
        AllowIntercompanyTxn: true,
        ReportingCurrency:   "USD",
    },
}

err = multiEngine.CreateCompany(parentCompany, "admin_user")
if err != nil {
    log.Fatal(err)
}

// Process intercompany transaction
icTransaction := &accounting.IntercompanyTransaction{
    SourceCompanyID: "parent_corp",
    TargetCompanyID: "subsidiary_llc",
    Amount:          &accounting.Amount{Value: 100000, Currency: "USD"},
    Description:     "Intercompany loan",
    ValidTime:       time.Now(),
}

err = multiEngine.ProcessIntercompanyTransaction(icTransaction, "admin_user")
if err != nil {
    log.Fatal(err)
}
```

## Design Principles

1. **Immutability**: Events are append-only, transactions cannot be modified after posting
2. **Auditability**: Complete audit trail with user and timestamp tracking
3. **Consistency**: Double-entry bookkeeping ensures books always balance
4. **Flexibility**: Multidimensional tagging for complex analytics
5. **Compliance**: Period locking and proper audit trails for regulatory compliance
6. **Enterprise Scale**: Multi-company support with consolidation capabilities
7. **Professional Reporting**: Financial statements that meet GAAP standards
8. **Separation of Concerns**: Clear separation between companies and elimination of intercompany transactions

## Enterprise Features

### Financial Reporting
- **Balance Sheet**: Complete asset, liability, and equity reporting with hierarchical categorization
- **Profit & Loss**: Period-based income statement with proper revenue and expense classification
- **Cash Flow**: Operating, investing, and financing activity analysis
- **Statement Formatting**: Professional formatted output suitable for executive reporting

### Multi-Company Consolidation
- **Entity Management**: Create and manage multiple business entities with individual settings
- **Intercompany Processing**: Automated intercompany transaction processing with matching entries
- **Consolidation Groups**: Define consolidation hierarchies for complex corporate structures
- **Elimination Rules**: Automatic elimination of intercompany balances for consolidated reporting

### Advanced Analytics
- **Financial Ratios**: Automated calculation of key financial metrics (ROA, profit margin, debt-to-equity)
- **Dimension Analysis**: Multi-dimensional reporting by department, product, cost center, etc.
- **Trend Analysis**: Period-over-period comparisons and variance analysis
- **Performance Metrics**: KPI calculations and business intelligence integration

## Demos

### Basic Demo (`cmd/demo/main.go`)
Demonstrates core accounting features including:
- Chart of accounts creation
- Transaction posting
- Trial balance generation
- Event sourcing and audit trails
- Accrual recognition
- Period management

### Advanced Demo (`cmd/advanced_demo/main.go`)
Showcases enterprise features including:
- Financial statement generation (Balance Sheet, P&L, Cash Flow)
- Multi-company operations
- Intercompany transaction processing
- Advanced analytics and financial ratios
- Consolidation group management

## Future Enhancements

- ✅ ~~Multi-company/multi-entity support~~ **COMPLETED**
- ✅ ~~Advanced reporting (P&L, Balance Sheet, Cash Flow)~~ **COMPLETED**
- ✅ ~~Tax calculation and regulatory compliance~~ **COMPLETED**
- Integration with external systems (banks, ERPs)
- Real-time posting with failure recovery
- Advanced reconciliation algorithms
- Machine learning for transaction categorization
- Web-based UI for enterprise management
- API layer for external integrations
- Performance optimization for large-scale operations

## Dependencies

- `go.etcd.io/bbolt` - Embedded key-value database
- `github.com/google/uuid` - UUID generation
- `github.com/stretchr/testify` - Testing framework (for comprehensive test suite)

## License

This implementation is provided as an educational example of advanced accounting system design patterns and enterprise-grade features.
