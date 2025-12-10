package accounting

import (
	"time"

	pb "accounting/proto/accounting"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProtoBytes serializes any protobuf message to bytes
func ToProtoBytes(msg proto.Message) ([]byte, error) {
	return proto.Marshal(msg)
}

// FromProtoBytes deserializes bytes into a protobuf message
func FromProtoBytes(data []byte, msg proto.Message) error {
	return proto.Unmarshal(data, msg)
}

// Helper function to convert time.Time to protobuf Timestamp
func timeToProto(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

// Helper function to convert protobuf Timestamp to time.Time
func protoToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

// Helper function to convert optional time.Time to protobuf Timestamp
func optionalTimeToProto(t *time.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(*t)
}

// Helper function to convert protobuf Timestamp to optional time.Time
func protoToOptionalTime(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

// ====================================================================================
// Amount Conversions
// ====================================================================================

func (a *Amount) ToProto() *pb.Amount {
	if a == nil {
		return nil
	}
	return &pb.Amount{
		Value:            a.Value,
		Currency:         string(a.Currency),
		BaseValue:        a.BaseValue,
		BaseCurrency:     string(a.BaseCurrency),
		ExchangeRate:     a.ExchangeRate,
		ExchangeRateDate: optionalTimeToProto(a.ExchangeRateDate),
	}
}

func AmountFromProto(pbAmount *pb.Amount) *Amount {
	if pbAmount == nil {
		return nil
	}
	return &Amount{
		Value:            pbAmount.Value,
		Currency:         Currency(pbAmount.Currency),
		BaseValue:        pbAmount.BaseValue,
		BaseCurrency:     Currency(pbAmount.BaseCurrency),
		ExchangeRate:     pbAmount.ExchangeRate,
		ExchangeRateDate: protoToOptionalTime(pbAmount.ExchangeRateDate),
	}
}

// ====================================================================================
// Dimension Conversions
// ====================================================================================

func (d *Dimension) ToProto() *pb.Dimension {
	if d == nil {
		return nil
	}
	return &pb.Dimension{
		Key:   string(d.Key),
		Value: d.Value,
	}
}

func DimensionFromProto(pbDim *pb.Dimension) *Dimension {
	if pbDim == nil {
		return nil
	}
	return &Dimension{
		Key:   DimensionKey(pbDim.Key),
		Value: pbDim.Value,
	}
}

func DimensionsToProto(dims []Dimension) []*pb.Dimension {
	result := make([]*pb.Dimension, len(dims))
	for i, d := range dims {
		result[i] = d.ToProto()
	}
	return result
}

func DimensionsFromProto(pbDims []*pb.Dimension) []Dimension {
	result := make([]Dimension, len(pbDims))
	for i, pbDim := range pbDims {
		result[i] = *DimensionFromProto(pbDim)
	}
	return result
}

// ====================================================================================
// Account Conversions
// ====================================================================================

func (a *Account) ToProto() *pb.Account {
	if a == nil {
		return nil
	}
	
	var accountType pb.AccountType
	switch a.Type {
	case Asset:
		accountType = pb.AccountType_ACCOUNT_TYPE_ASSET
	case Liability:
		accountType = pb.AccountType_ACCOUNT_TYPE_LIABILITY
	case Equity:
		accountType = pb.AccountType_ACCOUNT_TYPE_EQUITY
	case Income:
		accountType = pb.AccountType_ACCOUNT_TYPE_INCOME
	case Expense:
		accountType = pb.AccountType_ACCOUNT_TYPE_EXPENSE
	default:
		accountType = pb.AccountType_ACCOUNT_TYPE_UNSPECIFIED
	}
	
	return &pb.Account{
		Id:         a.ID,
		ParentId:   a.ParentID,
		Code:       a.Code,
		Name:       a.Name,
		Type:       accountType,
		Dimensions: DimensionsToProto(a.Dimensions),
		Currency:   string(a.Currency),
		CreatedAt:  timeToProto(a.CreatedAt),
		ClosedAt:   optionalTimeToProto(a.ClosedAt),
	}
}

func AccountFromProto(pbAcc *pb.Account) *Account {
	if pbAcc == nil {
		return nil
	}
	
	var accountType AccountType
	switch pbAcc.GetType() {
	case pb.AccountType_ACCOUNT_TYPE_ASSET:
		accountType = Asset
	case pb.AccountType_ACCOUNT_TYPE_LIABILITY:
		accountType = Liability
	case pb.AccountType_ACCOUNT_TYPE_EQUITY:
		accountType = Equity
	case pb.AccountType_ACCOUNT_TYPE_INCOME:
		accountType = Income
	case pb.AccountType_ACCOUNT_TYPE_EXPENSE:
		accountType = Expense
	}
	
	return &Account{
		ID:         pbAcc.Id,
		ParentID:   pbAcc.ParentId,
		Code:       pbAcc.Code,
		Name:       pbAcc.Name,
		Type:       accountType,
		Dimensions: DimensionsFromProto(pbAcc.Dimensions),
		Currency:   Currency(pbAcc.Currency),
		CreatedAt:  protoToTime(pbAcc.CreatedAt),
		ClosedAt:   protoToOptionalTime(pbAcc.ClosedAt),
	}
}

// ====================================================================================
// Entry Conversions
// ====================================================================================

func (e *Entry) ToProto() *pb.Entry {
	if e == nil {
		return nil
	}
	
	var entryType pb.EntryType
	switch e.Type {
	case Debit:
		entryType = pb.EntryType_ENTRY_TYPE_DEBIT
	case Credit:
		entryType = pb.EntryType_ENTRY_TYPE_CREDIT
	default:
		entryType = pb.EntryType_ENTRY_TYPE_UNSPECIFIED
	}
	
	return &pb.Entry{
		Id:            e.ID,
		TransactionId: e.TransactionID,
		AccountId:     e.AccountID,
		Type:          entryType,
		Amount:        e.Amount.ToProto(),
		Dimensions:    DimensionsToProto(e.Dimensions),
	}
}

func EntryFromProto(pbEntry *pb.Entry) *Entry {
	if pbEntry == nil {
		return nil
	}
	
	var entryType EntryType
	switch pbEntry.GetType() {
	case pb.EntryType_ENTRY_TYPE_DEBIT:
		entryType = Debit
	case pb.EntryType_ENTRY_TYPE_CREDIT:
		entryType = Credit
	}
	
	amount := AmountFromProto(pbEntry.Amount)
	if amount == nil {
		amount = &Amount{} // Default to empty amount if nil
	}
	
	return &Entry{
		ID:            pbEntry.Id,
		TransactionID: pbEntry.TransactionId,
		AccountID:     pbEntry.AccountId,
		Type:          entryType,
		Amount:        *amount,
		Dimensions:    DimensionsFromProto(pbEntry.Dimensions),
	}
}

func EntriesToProto(entries []Entry) []*pb.Entry {
	result := make([]*pb.Entry, len(entries))
	for i, e := range entries {
		result[i] = e.ToProto()
	}
	return result
}

func EntriesFromProto(pbEntries []*pb.Entry) []Entry {
	result := make([]Entry, 0, len(pbEntries))
	for _, pbEntry := range pbEntries {
		entry := EntryFromProto(pbEntry)
		if entry != nil {
			result = append(result, *entry)
		}
	}
	return result
}

// ====================================================================================
// Transaction Conversions
// ====================================================================================

func (t *Transaction) ToProto() *pb.Transaction {
	if t == nil {
		return nil
	}
	
	var status pb.TransactionStatus
	switch t.Status {
	case Pending:
		status = pb.TransactionStatus_TRANSACTION_STATUS_PENDING
	case Posted:
		status = pb.TransactionStatus_TRANSACTION_STATUS_POSTED
	case Reversed:
		status = pb.TransactionStatus_TRANSACTION_STATUS_REVERSED
	case InBatch:
		status = pb.TransactionStatus_TRANSACTION_STATUS_IN_BATCH
	default:
		status = pb.TransactionStatus_TRANSACTION_STATUS_UNSPECIFIED
	}
	
	return &pb.Transaction{
		Id:              t.ID,
		Description:     t.Description,
		ValidTime:       timeToProto(t.ValidTime),
		TransactionTime: timeToProto(t.TransactionTime),
		Status:          status,
		Entries:         EntriesToProto(t.Entries),
		SourceRef:       t.SourceRef,
		UserId:          t.UserID,
		CreatedAt:       timeToProto(t.CreatedAt),
		UpdatedAt:       timeToProto(t.UpdatedAt),
	}
}

func TransactionFromProto(pbTxn *pb.Transaction) *Transaction {
	if pbTxn == nil {
		return nil
	}
	
	var status TransactionStatus
	switch pbTxn.GetStatus() {
	case pb.TransactionStatus_TRANSACTION_STATUS_PENDING:
		status = Pending
	case pb.TransactionStatus_TRANSACTION_STATUS_POSTED:
		status = Posted
	case pb.TransactionStatus_TRANSACTION_STATUS_REVERSED:
		status = Reversed
	case pb.TransactionStatus_TRANSACTION_STATUS_IN_BATCH:
		status = InBatch
	}
	
	return &Transaction{
		ID:              pbTxn.Id,
		Description:     pbTxn.Description,
		ValidTime:       protoToTime(pbTxn.ValidTime),
		TransactionTime: protoToTime(pbTxn.TransactionTime),
		Status:          status,
		Entries:         EntriesFromProto(pbTxn.Entries),
		SourceRef:       pbTxn.SourceRef,
		UserID:          pbTxn.UserId,
		CreatedAt:       protoToTime(pbTxn.CreatedAt),
		UpdatedAt:       protoToTime(pbTxn.UpdatedAt),
	}
}

// ToBytes serializes a Transaction to protobuf bytes
func (t *Transaction) ToBytes() ([]byte, error) {
	return ToProtoBytes(t.ToProto())
}

// TransactionFromBytes deserializes a Transaction from protobuf bytes
func TransactionFromBytes(data []byte) (*Transaction, error) {
	pbTxn := &pb.Transaction{}
	if err := FromProtoBytes(data, pbTxn); err != nil {
		return nil, err
	}
	return TransactionFromProto(pbTxn), nil
}

// ====================================================================================
// Period Conversions
// ====================================================================================

func (p *Period) ToProto() *pb.Period {
	if p == nil {
		return nil
	}
	return &pb.Period{
		Id:           p.ID,
		Name:         p.Name,
		Start:        timeToProto(p.Start),
		End:          timeToProto(p.End),
		SoftClosedAt: optionalTimeToProto(p.SoftClosedAt),
		HardClosedAt: optionalTimeToProto(p.HardClosedAt),
	}
}

func PeriodFromProto(pbPeriod *pb.Period) *Period {
	if pbPeriod == nil {
		return nil
	}
	return &Period{
		ID:           pbPeriod.Id,
		Name:         pbPeriod.Name,
		Start:        protoToTime(pbPeriod.Start),
		End:          protoToTime(pbPeriod.End),
		SoftClosedAt: protoToOptionalTime(pbPeriod.SoftClosedAt),
		HardClosedAt: protoToOptionalTime(pbPeriod.HardClosedAt),
	}
}

// ====================================================================================
// JournalEvent Conversions
// ====================================================================================

func (j *JournalEvent) ToProto() *pb.JournalEvent {
	if j == nil {
		return nil
	}
	return &pb.JournalEvent{
		Id:              j.ID,
		EventType:       j.EventType,
		Payload:         j.Payload,
		ValidTime:       timeToProto(j.ValidTime),
		TransactionTime: timeToProto(j.TransactionTime),
		UserId:          j.UserID,
	}
}

func JournalEventFromProto(pbEvent *pb.JournalEvent) *JournalEvent {
	if pbEvent == nil {
		return nil
	}
	return &JournalEvent{
		ID:              pbEvent.Id,
		EventType:       pbEvent.EventType,
		Payload:         pbEvent.Payload,
		ValidTime:       protoToTime(pbEvent.ValidTime),
		TransactionTime: protoToTime(pbEvent.TransactionTime),
		UserID:          pbEvent.UserId,
	}
}
