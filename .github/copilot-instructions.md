# Copilot Instructions for Accounting System

## Repository Overview

This is a comprehensive accounting software library written in Go, designed as an embedded third-party package. It implements advanced accounting features including event sourcing, double-entry bookkeeping, bi-temporal data tracking, multi-currency support, and regulatory compliance.

## Technology Stack

- **Language**: Go 1.24.2
- **Storage**: BoltDB (etcd-io/bbolt) - embedded key-value database
- **Serialization**: Protocol Buffers (google.golang.org/protobuf)
- **Testing**: Standard Go testing + testify
- **Key Dependencies**: google/uuid

## Architecture & Key Concepts

### Core Accounting Principles

1. **Double-Entry Bookkeeping**: Every transaction must have balanced debits and credits. The fundamental accounting equation (Assets = Liabilities + Equity) must always hold.

2. **Bi-Temporal Data**: Track both:
   - **Valid Time**: When the transaction occurred in the business domain
   - **Transaction Time**: When it was recorded in the system

3. **Event Sourcing**: All changes are recorded as immutable events. The current state is derived by replaying events. Never modify existing transactions; create reversing/adjusting entries instead.

4. **Immutability**: Financial records are append-only. Use versioning and audit trails for all changes.

### System Components

The main entry point is `AccountingEngine` which coordinates these services:

- **PostingEngine**: Creates and validates transactions, ensures balanced entries
- **EventStore**: Manages immutable event log
- **QueryAPI**: Queries transactions with bi-temporal filters
- **AccrualService**: Handles revenue/expense recognition schedules
- **ReconciliationService**: Matches bank statements with internal records
- **ReportingService**: Generates financial statements (Balance Sheet, P&L, Cash Flow)
- **ZBBService**: Zero-based budgeting
- **ComplianceService**: GAAP, IFRS, SOX compliance validation
- **AMLService**: Anti-money laundering detection (27+ rule types)
- **ForensicService**: Financial forensics and fraud detection
- **MultiCompanyService**: Multi-entity consolidation

### Data Model

- **Account**: Chart of accounts with hierarchical structure
- **Transaction**: Container for multiple entries
- **Entry**: Individual debit/credit line with amount and currency
- **Amount**: Value in cents/smallest unit + currency code
- **Period**: Time period with soft/hard close status
- **Dimension**: Multi-dimensional tagging (department, project, region, etc.)
- **RecognitionSchedule**: For accruals and deferrals

### Protocol Buffers

All data types are serialized using Protocol Buffers. The `.proto` files are in `proto/accounting/`. After modifying proto files, regenerate with:

```bash
protoc --go_out=. --go_opt=paths=source_relative proto/accounting/*.proto
```

## Building & Testing

### Run Tests

```bash
# Run all tests
go test ./...

# Run specific test file
go test -v -run TestAccountingSystem

# Run with coverage
go test -cover ./...
```

### Build Demo Applications

```bash
# Build specific demo
go build ./cmd/demo
go build ./cmd/advanced_demo
go build ./cmd/aml_demo
go build ./cmd/integration_demo

# Run integration tests
go test ./cmd/integration_demo/...
go test ./cmd/complete_test/...
```

### Dependencies

```bash
# Download dependencies
go mod download

# Update dependencies
go mod tidy
```

## Code Organization

- **Root directory**: Core accounting types and services (*.go files)
- **proto/accounting/**: Protocol Buffer schema definitions
- **cmd/**: Demo applications and integration tests
- **Documentation**:
  - `readme.md`: Project overview and features
  - `IMPLEMENTATION.md`: Feature implementation details
  - `AML_IMPLEMENTATION_SUMMARY.md`: AML rules documentation
  - `AML_RULES_DOCUMENTATION.md`: Detailed AML rule specifications
  - `PROTOBUF_GUIDE.md`: Protocol Buffer usage guide
  - `PROTOBUF_MIGRATION_SUMMARY.md`: Migration notes

## Development Guidelines

### When Making Changes

1. **Preserve Immutability**: Never modify existing transactions. Create reversing entries or new versions.

2. **Maintain Balance**: Always ensure debits equal credits. Use the `PostingEngine` for transaction creation.

3. **Use Proper Types**:
   - Store money amounts as `int64` in cents (smallest currency unit)
   - Use `Amount` type with currency code
   - Use `uuid.UUID` for IDs
   - Use `time.Time` for timestamps

4. **Event Sourcing Pattern**:
   - Emit events for all state changes
   - Store events in EventStore
   - Make state derivable from events

5. **Testing**:
   - Use temporary database files: `/tmp/test_*.db`
   - Clean up with `defer os.Remove(dbFile)`
   - Test both success and error paths
   - Verify accounting equation remains balanced

6. **Error Handling**:
   - Return descriptive errors with context
   - Use `fmt.Errorf` with `%w` for error wrapping
   - Validate inputs before processing

### Code Style

- Follow standard Go conventions
- Use meaningful variable names
- Add comments for complex accounting logic
- Document public APIs with godoc comments
- Keep functions focused and testable

### Important Constraints

1. **No CRUD Updates**: Use append-only patterns. Create new records instead of updating.

2. **Always Validate Balance**: Before saving a transaction, verify total debits equal total credits.

3. **Respect Period Locks**: Hard-closed periods cannot accept new transactions. Create adjusting entries in open periods.

4. **Multi-Currency**: Store original currency amounts and exchange rates. Calculate base currency equivalents.

5. **Audit Trail**: Track user ID, timestamps, and reasons for all changes.

6. **Regulatory Compliance**: Ensure changes comply with GAAP/IFRS standards as configured.

### AML & Compliance

When working with AML features:
- Reference `AML_RULES_DOCUMENTATION.md` for rule specifications
- 27+ rule types covering CTR, SAR, structuring, PEP, sanctions, etc.
- Rules organized by: Transaction-based, Customer-based, Geographic, Advanced detection
- Use `AMLService` for automated detection
- Always maintain audit trails for compliance investigations

### Common Patterns

#### Creating a Transaction

```go
tx := &Transaction{
    Description: "Sale of goods",
    ValidTime:   time.Now(),
    Entries: []Entry{
        {AccountID: "cash", Type: Debit, Amount: Amount{Value: 100000, Currency: "USD"}},
        {AccountID: "revenue", Type: Credit, Amount: Amount{Value: 100000, Currency: "USD"}},
    },
}
// First create the transaction
err := engine.CreateTransaction(tx, userID)
if err != nil {
    // handle error
}
// Then post it
err = engine.PostTransaction(tx.ID, userID)
```

#### Querying Balances

```go
balance, err := engine.GetAccountBalance("cash", time.Now())
// balance is *BalanceResult with balance.Balance.Value (in cents) and balance.Balance.Currency
```

#### Running AML Checks

```go
amlService := engine.GetAMLService()
customerInfo := map[string]*AMLCustomer{
    "customer1": {CustomerID: "customer1", RiskLevel: "medium"},
}
alerts, err := amlService.MonitorTransaction(transaction, customerInfo)
```

## Need Help?

- Review test files for usage examples: `example_test.go`, `advanced_features_test.go`, `zbb_test.go`
- Check demo applications in `cmd/` directory
- Refer to documentation files in root directory
- Accounting principles must be preserved - when in doubt, ask for clarification

## API Design

This library is designed as an embedded package with no HTTP API. All interactions happen through the `AccountingEngine` Go API. Focus on type safety, error handling, and maintaining accounting invariants.
