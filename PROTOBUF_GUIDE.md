# Protocol Buffer Integration Guide

This project uses Protocol Buffers (protobuf) for efficient binary serialization of data types, replacing JSON serialization.

## Overview

All JSON-serializable types have been migrated to Protocol Buffer schema definitions. This provides:
- **Binary serialization**: More compact than JSON (typically 3-10x smaller)
- **Type safety**: Strongly typed schema definitions
- **Performance**: Faster serialization/deserialization
- **Cross-language support**: Can be used with other languages if needed
- **Schema evolution**: Built-in support for backward/forward compatibility

## Directory Structure

```
proto/
└── accounting/
    ├── accounting.proto         # Core accounting types
    ├── aml.proto               # Anti-Money Laundering types
    ├── zbb.proto               # Zero-Based Budgeting types
    ├── compliance.proto        # Compliance & tax types
    ├── multi_company.proto     # Multi-company types
    ├── forensic.proto          # Forensic accounting types
    ├── reporting.proto         # Financial reporting types
    ├── accrual.proto          # Accrual/deferral types
    ├── reconciliation.proto   # Reconciliation types
    └── *.pb.go                # Generated Go code
```

## Setup

### Prerequisites

1. **Install Protocol Buffer Compiler (protoc)**:
   ```bash
   # On Ubuntu/Debian
   sudo apt-get install protobuf-compiler
   
   # On macOS
   brew install protobuf
   
   # Verify installation
   protoc --version
   ```

2. **Install Go protobuf plugin**:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   ```

3. **Add Go bin to PATH**:
   ```bash
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

### Regenerating Proto Files

If you modify any `.proto` files, regenerate the Go code:

```bash
cd /path/to/fin
protoc --go_out=. --go_opt=paths=source_relative proto/accounting/*.proto
```

## Usage

### Converting Between Types

The project provides conversion functions between existing Go types and protobuf types in `proto_converters.go`.

#### Example: Transaction Serialization

```go
// Create a transaction
txn := &Transaction{
    ID:          "txn-001",
    Description: "Payment received",
    ValidTime:   time.Now(),
    Status:      Posted,
    Entries:     []Entry{...},
}

// Convert to protobuf bytes
data, err := txn.ToBytes()
if err != nil {
    log.Fatal(err)
}

// Save to storage (data is now compact binary format)
storage.Save("transactions/txn-001", data)

// Later, deserialize from bytes
txn2, err := TransactionFromBytes(data)
if err != nil {
    log.Fatal(err)
}
```

#### Example: Using Proto Directly

```go
import pb "accounting/proto/accounting"

// Create a proto message directly
pbTxn := &pb.Transaction{
    Id:          "txn-002",
    Description: "Invoice payment",
    Status:      pb.TransactionStatus_TRANSACTION_STATUS_POSTED,
    // ...
}

// Serialize to bytes
data, err := proto.Marshal(pbTxn)
if err != nil {
    log.Fatal(err)
}

// Deserialize from bytes
pbTxn2 := &pb.Transaction{}
err = proto.Unmarshal(data, pbTxn2)
```

### Available Conversion Methods

For each major type, there are `ToProto()` and `FromProto()` methods:

```go
// Amount
amount := &Amount{Value: 100000, Currency: "USD"}
pbAmount := amount.ToProto()
amount2 := AmountFromProto(pbAmount)

// Account
account := &Account{ID: "acc-001", Code: "1000", Type: Asset}
pbAccount := account.ToProto()
account2 := AccountFromProto(pbAccount)

// Transaction
pbTxn := txn.ToProto()
txn2 := TransactionFromProto(pbTxn)

// Entry
pbEntry := entry.ToProto()
entry2 := EntryFromProto(pbEntry)
```

### Serialization Helpers

```go
// Generic protobuf serialization
func ToProtoBytes(msg proto.Message) ([]byte, error)
func FromProtoBytes(data []byte, msg proto.Message) error

// Transaction-specific helpers
func (t *Transaction) ToBytes() ([]byte, error)
func TransactionFromBytes(data []byte) (*Transaction, error)
```

## Benefits

### Size Comparison

Typical size reduction compared to JSON:

| Type          | JSON Size | Protobuf Size | Reduction |
|---------------|-----------|---------------|-----------|
| Transaction   | ~2.5 KB   | ~800 bytes    | ~68%      |
| Account       | ~500 bytes| ~150 bytes    | ~70%      |
| AML Alert     | ~3 KB     | ~1 KB         | ~67%      |

### Performance

Protobuf serialization is typically 2-5x faster than JSON:

```go
// Benchmark results (example)
BenchmarkJSON-8         50000    24500 ns/op
BenchmarkProto-8       200000     6200 ns/op
```

## Schema Evolution

Protocol Buffers support schema evolution through field numbers:

1. **Never change field numbers** - This breaks compatibility
2. **Adding fields** - OK, old code ignores new fields
3. **Removing fields** - Mark as reserved
4. **Changing types** - Generally not safe

Example of adding a new field:

```protobuf
message Transaction {
  string id = 1;
  string description = 2;
  // ... existing fields ...
  string new_field = 11;  // New field with unique number
}
```

## Proto Schema Reference

### Core Types (accounting.proto)

- `Amount`: Monetary values with multi-currency support
- `Dimension`: Key-value tags for OLAP analytics
- `Account`: Chart of accounts node
- `Entry`: Single debit/credit line
- `Transaction`: Bi-temporal transaction
- `Period`: Accounting period
- `Ledger`: Ledger/sub-ledger
- `Reconciliation`: Reconciliation tracking
- `RecognitionSchedule`: Accrual/deferral schedule
- `JournalEvent`: Event sourcing record

### AML Types (aml.proto)

- `AMLAlert`: Anti-money laundering alert
- `AMLRule`: AML compliance rule
- `AMLCustomer`: Customer AML profile
- `AMLTransaction`: Transaction with AML metadata
- `AMLInvestigation`: Investigation record
- `AMLDashboard`: Compliance dashboard data

### ZBB Types (zbb.proto)

- `BudgetPeriod`: Budget cycle
- `BudgetRequest`: Zero-based budget request
- `BudgetLineItem`: Individual budget line
- `Justification`: Business justification
- `BudgetAllocation`: Approved allocation
- `BudgetTracking`: Budget vs actual tracking

## Testing

All existing tests continue to work with the protobuf integration. Run tests:

```bash
go test -v .
```

## Migration Notes

The existing code continues to use the original Go struct definitions. The protobuf types are used for:

1. **Serialization/Deserialization**: When saving to or loading from storage
2. **Wire Protocol**: If transmitting data over network (future)
3. **Event Payloads**: In event sourcing system

The conversion layer (`proto_converters.go`) bridges between the two representations transparently.

## Future Enhancements

Potential improvements:

1. **Direct Proto Usage**: Gradually migrate to using proto types directly in services
2. **gRPC Integration**: Add gRPC API using proto definitions
3. **Cross-Language Clients**: Generate clients in Python, Java, etc.
4. **Streaming**: Use proto for efficient event streaming
5. **Validation**: Add proto validation rules using proto-gen-validate

## Troubleshooting

### Proto Generation Fails

```bash
# Make sure protoc-gen-go is in PATH
which protoc-gen-go

# If not found, reinstall
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

### Import Errors

If you get import errors for `accounting/proto/accounting`:

```bash
# Run go mod tidy
go mod tidy

# Ensure the proto package builds
go build ./proto/accounting/...
```

### Type Conversion Issues

Check that you're using the correct conversion functions:
- Use `ToProto()` methods on Go structs
- Use `FromProto()` functions with `pb` prefix for parameters
- Use helper functions like `ToProtoBytes()` for serialization

## Resources

- [Protocol Buffers Documentation](https://protobuf.dev/)
- [Go Protocol Buffer Tutorial](https://protobuf.dev/getting-started/gotutorial/)
- [Protocol Buffer Style Guide](https://protobuf.dev/programming-guides/style/)
