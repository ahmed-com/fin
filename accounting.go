package accounting

// Primitive data structures for a minimally‑viable, yet extensible, accounting subsystem.
// Focus: immutable event log, double‑entry ledger, bi‑temporal bookkeeping, multi‑currency,
// multidimensional analytics, period closing, reconciliation, and compliance contexts.
// No business logic or calculations are included – those will be layered on in services.

import "time"

// ----------------------------------------------------------------------------
// ⏳ Time & Currency ----------------------------------------------------------------
// ----------------------------------------------------------------------------

// Currency is an ISO‑4217 code (e.g., "USD", "EGP").
type Currency string

// ----------------------------------------------------------------------------
// 📐 Multidimensional Tagging -------------------------------------------------------
// ----------------------------------------------------------------------------

type DimensionKey string

const (
    DimDepartment DimensionKey = "department"
    DimProduct    DimensionKey = "product"
    DimProject    DimensionKey = "project"
    DimRegion     DimensionKey = "region"
    DimCostCenter DimensionKey = "cost_center"
)

// Dimension is an arbitrary key/value tag that can be attached to any business fact
// (account, entry, transaction, schedule …) for OLAP‑style analytics.
type Dimension struct {
    Key   DimensionKey `json:"key"`
    Value string       `json:"value"`
}

// ----------------------------------------------------------------------------
// 💲 Monetary Amounts (Multi‑Currency) ----------------------------------------
// ----------------------------------------------------------------------------

type Amount struct {
    // Value in the smallest unit (cents, piasters, satang …).
    Value int64 `json:"value"`
    // The original currency of the transaction.
    Currency Currency `json:"currency"`

    // Optional base‑currency projection (filled when posting into a base ledger).
    BaseValue       int64     `json:"base_value,omitempty"`
    BaseCurrency    Currency  `json:"base_currency,omitempty"`
    ExchangeRate    float64   `json:"exchange_rate,omitempty"`
    ExchangeRateDate *time.Time `json:"exchange_rate_date,omitempty"`
}

// ----------------------------------------------------------------------------
// 🗃️  Chart of Accounts -----------------------------------------------------------
// ----------------------------------------------------------------------------

type AccountType string

const (
    Asset     AccountType = "ASSET"
    Liability AccountType = "LIABILITY"
    Equity    AccountType = "EQUITY"
    Income    AccountType = "INCOME"
    Expense   AccountType = "EXPENSE"
)

// Account is a node in a tree‑structured Chart of Accounts.
type Account struct {
    ID         string      `json:"id"`
    ParentID   string      `json:"parent_id,omitempty"`
    Code       string      `json:"code"` // e.g., "1001"
    Name       string      `json:"name"`
    Type       AccountType `json:"type"`

    // Default analytical tags that cascade to entries.
    Dimensions []Dimension `json:"dimensions,omitempty"`

    Currency   Currency   `json:"currency,omitempty"`
    CreatedAt  time.Time  `json:"created_at"`
    ClosedAt   *time.Time `json:"closed_at,omitempty"`
}

// ----------------------------------------------------------------------------
// 📜 Double‑Entry Journal ---------------------------------------------------------
// ----------------------------------------------------------------------------

type EntryType string

const (
    Debit  EntryType = "DEBIT"
    Credit EntryType = "CREDIT"
)

// Entry is a single debit or credit line.
type Entry struct {
    ID            string      `json:"id"`
    TransactionID string      `json:"transaction_id"`
    AccountID     string      `json:"account_id"`
    Type          EntryType   `json:"type"`
    Amount        Amount      `json:"amount"`
    Dimensions    []Dimension `json:"dimensions,omitempty"`
}

// ----------------------------------------------------------------------------
// ⏳⏳ Bi‑Temporal Transaction -----------------------------------------------
// ----------------------------------------------------------------------------

type TransactionStatus string

const (
    Pending  TransactionStatus = "PENDING"  // not yet posted
    Posted   TransactionStatus = "POSTED"   // posted to ledger
    Reversed TransactionStatus = "REVERSED" // reversed by a contra txn
    InBatch  TransactionStatus = "IN_BATCH" // waiting for nightly batch
)

type Transaction struct {
    ID              string            `json:"id"`
    Description     string            `json:"description,omitempty"`

    // Bi‑temporal coordinates.
    ValidTime       time.Time         `json:"valid_time"`       // effective/business date
    TransactionTime time.Time         `json:"transaction_time"` // when recorded in the system

    Status          TransactionStatus `json:"status"`
    Entries         []Entry           `json:"entries"`

    // Metadata -------------------------------------------------------------------
    SourceRef string    `json:"source_ref,omitempty"` // e.g., invoice‑ID, external UUID
    UserID    string    `json:"user_id,omitempty"`   // who created/modified

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// ----------------------------------------------------------------------------
// 📅 Accounting Periods & Closing --------------------------------------------------
// ----------------------------------------------------------------------------

type Period struct {
    ID   string    `json:"id"`
    Name string    `json:"name"` // e.g., "2025‑Q2"
    Start time.Time `json:"start"`
    End   time.Time `json:"end"`

    SoftClosedAt *time.Time `json:"soft_closed_at,omitempty"` // still editable
    HardClosedAt *time.Time `json:"hard_closed_at,omitempty"` // immutable
}

// ----------------------------------------------------------------------------
// 📚 Ledgers & Sub‑Ledgers --------------------------------------------------------
// ----------------------------------------------------------------------------

type LedgerType string

const (
    GeneralLedger      LedgerType = "GL"
    AccountsReceivable LedgerType = "AR"
    AccountsPayable    LedgerType = "AP"
    InventoryLedger    LedgerType = "INV"
)

type Ledger struct {
    ID       string     `json:"id"`
    Name     string     `json:"name"`
    Type     LedgerType `json:"type"`
    Currency Currency   `json:"currency,omitempty"`
}

// ----------------------------------------------------------------------------
// 🔍 Reconciliation --------------------------------------------------------------
// ----------------------------------------------------------------------------

type ReconciliationStatus string

const (
    Unreconciled ReconciliationStatus = "UNRECONCILED"
    Reconciled   ReconciliationStatus = "RECONCILED"
    Partial      ReconciliationStatus = "PARTIAL"
)

type Reconciliation struct {
    ID          string               `json:"id"`
    ExternalRef string               `json:"external_ref"` // bank‑statement line, etc.
    EntryIDs    []string             `json:"entry_ids"`
    Status      ReconciliationStatus `json:"status"`
    CreatedAt   time.Time            `json:"created_at"`
    CompletedAt *time.Time           `json:"completed_at,omitempty"`
}

// ----------------------------------------------------------------------------
// 🧾 Accrual / Deferral Recognition Schedules ------------------------------------
// ----------------------------------------------------------------------------

type ScheduleFrequency string

const (
    Monthly   ScheduleFrequency = "MONTHLY"
    Quarterly ScheduleFrequency = "QUARTERLY"
    Yearly    ScheduleFrequency = "YEARLY"
)

type RecognitionSchedule struct {
    ID            string            `json:"id"`
    TransactionID string            `json:"transaction_id"`
    Frequency     ScheduleFrequency `json:"frequency"`
    Occurrences   int               `json:"occurrences"` // e.g., 12 months
    StartTime     time.Time         `json:"start_time"`
    CreatedAt     time.Time         `json:"created_at"`
}

// ----------------------------------------------------------------------------
// 📑 Compliance / Reporting Contexts --------------------------------------------
// ----------------------------------------------------------------------------

type ReportingStandard string

const (
    GAAP ReportingStandard = "GAAP"
    IFRS ReportingStandard = "IFRS"
)

type ReportingContext struct {
    ID       string            `json:"id"`
    Standard ReportingStandard `json:"standard"`
    Currency Currency          `json:"currency"`
    PeriodID string            `json:"period_id"`
}

// ----------------------------------------------------------------------------
// 📝 Event Sourcing --------------------------------------------------------------
// ----------------------------------------------------------------------------

// JournalEvent is the atomic, append‑only log record used to reconstruct state.
type JournalEvent struct {
    ID              string    `json:"id"`
    EventType       string    `json:"event_type"` // e.g., "CREATE_TXN"
    Payload         []byte    `json:"payload"`    // serialized protobuf/JSON/etc.
    ValidTime       time.Time `json:"valid_time"`
    TransactionTime time.Time `json:"transaction_time"`
    UserID          string    `json:"user_id,omitempty"`
}
