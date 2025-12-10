package accounting

import (
	pb "accounting/proto/accounting"
	"google.golang.org/protobuf/proto"
)

// ====================================================================================
// Ledger Conversions
// ====================================================================================

func (l *Ledger) ToProto() *pb.Ledger {
	if l == nil {
		return nil
	}
	
	var ledgerType pb.LedgerType
	switch l.Type {
	case GeneralLedger:
		ledgerType = pb.LedgerType_LEDGER_TYPE_GL
	case AccountsReceivable:
		ledgerType = pb.LedgerType_LEDGER_TYPE_AR
	case AccountsPayable:
		ledgerType = pb.LedgerType_LEDGER_TYPE_AP
	default:
		ledgerType = pb.LedgerType_LEDGER_TYPE_UNSPECIFIED
	}
	
	return &pb.Ledger{
		Id:       l.ID,
		Name:     l.Name,
		Type:     ledgerType,
		Currency: string(l.Currency),
	}
}

func LedgerFromProto(pbLedger *pb.Ledger) *Ledger {
	if pbLedger == nil {
		return nil
	}
	
	var ledgerType LedgerType
	switch pbLedger.GetType() {
	case pb.LedgerType_LEDGER_TYPE_GL:
		ledgerType = GeneralLedger
	case pb.LedgerType_LEDGER_TYPE_AR:
		ledgerType = AccountsReceivable
	case pb.LedgerType_LEDGER_TYPE_AP:
		ledgerType = AccountsPayable
	}
	
	return &Ledger{
		ID:       pbLedger.Id,
		Name:     pbLedger.Name,
		Type:     ledgerType,
		Currency: Currency(pbLedger.Currency),
	}
}

// ====================================================================================
// Reconciliation Conversions
// ====================================================================================

func (r *Reconciliation) ToProto() *pb.Reconciliation {
	if r == nil {
		return nil
	}
	
	var status pb.ReconciliationStatus
	switch r.Status {
	case Unreconciled:
		status = pb.ReconciliationStatus_RECONCILIATION_STATUS_UNRECONCILED
	case Reconciled:
		status = pb.ReconciliationStatus_RECONCILIATION_STATUS_RECONCILED
	case Partial:
		status = pb.ReconciliationStatus_RECONCILIATION_STATUS_PARTIAL
	default:
		status = pb.ReconciliationStatus_RECONCILIATION_STATUS_UNSPECIFIED
	}
	
	return &pb.Reconciliation{
		Id:          r.ID,
		ExternalRef: r.ExternalRef,
		EntryIds:    r.EntryIDs,
		Status:      status,
		CreatedAt:   timeToProto(r.CreatedAt),
		CompletedAt: optionalTimeToProto(r.CompletedAt),
	}
}

func ReconciliationFromProto(pbRecon *pb.Reconciliation) *Reconciliation {
	if pbRecon == nil {
		return nil
	}
	
	var status ReconciliationStatus
	switch pbRecon.GetStatus() {
	case pb.ReconciliationStatus_RECONCILIATION_STATUS_UNRECONCILED:
		status = Unreconciled
	case pb.ReconciliationStatus_RECONCILIATION_STATUS_RECONCILED:
		status = Reconciled
	case pb.ReconciliationStatus_RECONCILIATION_STATUS_PARTIAL:
		status = Partial
	}
	
	return &Reconciliation{
		ID:          pbRecon.Id,
		ExternalRef: pbRecon.ExternalRef,
		EntryIDs:    pbRecon.EntryIds,
		Status:      status,
		CreatedAt:   protoToTime(pbRecon.CreatedAt),
		CompletedAt: protoToOptionalTime(pbRecon.CompletedAt),
	}
}

// ====================================================================================
// RecognitionSchedule Conversions
// ====================================================================================

func (r *RecognitionSchedule) ToProto() *pb.RecognitionSchedule {
	if r == nil {
		return nil
	}
	
	var freq pb.ScheduleFrequency
	switch r.Frequency {
	case Monthly:
		freq = pb.ScheduleFrequency_SCHEDULE_FREQUENCY_MONTHLY
	case Quarterly:
		freq = pb.ScheduleFrequency_SCHEDULE_FREQUENCY_QUARTERLY
	case Yearly:
		freq = pb.ScheduleFrequency_SCHEDULE_FREQUENCY_YEARLY
	default:
		freq = pb.ScheduleFrequency_SCHEDULE_FREQUENCY_UNSPECIFIED
	}
	
	return &pb.RecognitionSchedule{
		Id:            r.ID,
		TransactionId: r.TransactionID,
		Frequency:     freq,
		Occurrences:   int32(r.Occurrences),
		StartTime:     timeToProto(r.StartTime),
		CreatedAt:     timeToProto(r.CreatedAt),
	}
}

func RecognitionScheduleFromProto(pbSched *pb.RecognitionSchedule) *RecognitionSchedule {
	if pbSched == nil {
		return nil
	}
	
	var freq ScheduleFrequency
	switch pbSched.GetFrequency() {
	case pb.ScheduleFrequency_SCHEDULE_FREQUENCY_MONTHLY:
		freq = Monthly
	case pb.ScheduleFrequency_SCHEDULE_FREQUENCY_QUARTERLY:
		freq = Quarterly
	case pb.ScheduleFrequency_SCHEDULE_FREQUENCY_YEARLY:
		freq = Yearly
	}
	
	return &RecognitionSchedule{
		ID:            pbSched.Id,
		TransactionID: pbSched.TransactionId,
		Frequency:     freq,
		Occurrences:   int(pbSched.Occurrences),
		StartTime:     protoToTime(pbSched.StartTime),
		CreatedAt:     protoToTime(pbSched.CreatedAt),
	}
}

// ToBytes serializes to protobuf bytes
func (r *RecognitionSchedule) ToBytes() ([]byte, error) {
	return proto.Marshal(r.ToProto())
}

// RecognitionScheduleFromBytes deserializes from protobuf bytes
func RecognitionScheduleFromBytes(data []byte) (*RecognitionSchedule, error) {
	pbSched := &pb.RecognitionSchedule{}
	if err := proto.Unmarshal(data, pbSched); err != nil {
		return nil, err
	}
	return RecognitionScheduleFromProto(pbSched), nil
}

// ====================================================================================
// Company and Multi-Company Conversions
// ====================================================================================

func (c *Company) ToProto() *pb.Company {
if c == nil {
return nil
}

var status pb.CompanyStatus
switch c.Status {
case CompanyActive:
status = pb.CompanyStatus_COMPANY_STATUS_ACTIVE
case CompanyInactive:
status = pb.CompanyStatus_COMPANY_STATUS_INACTIVE
case CompanySuspended:
status = pb.CompanyStatus_COMPANY_STATUS_SUSPENDED
case CompanyMerged:
status = pb.CompanyStatus_COMPANY_STATUS_MERGED
default:
status = pb.CompanyStatus_COMPANY_STATUS_UNSPECIFIED
}

// Convert metadata map to string map
metadata := make(map[string]string)
for k, v := range c.Metadata {
if str, ok := v.(string); ok {
metadata[k] = str
}
}

var address *pb.Address
if c.Address != nil {
address = &pb.Address{
Street1:    c.Address.Street1,
Street2:    c.Address.Street2,
City:       c.Address.City,
State:      c.Address.State,
PostalCode: c.Address.PostalCode,
Country:    c.Address.Country,
}
}

var settings *pb.CompanySettings
if c.Settings != nil {
var rules []*pb.AutoPostingRule
for _, rule := range c.Settings.AutoPostingRules {
if rule != nil {
var actions []*pb.PostingAction
for _, action := range rule.Actions {
if action != nil {
actions = append(actions, &pb.PostingAction{
Type:       action.Type,
Parameters: action.Parameters,
})
}
}
rules = append(rules, &pb.AutoPostingRule{
Id:        rule.ID,
Name:      rule.Name,
Condition: rule.Condition,
Actions:   actions,
IsActive:  rule.IsActive,
CreatedAt: timeToProto(rule.CreatedAt),
})
}
}

settings = &pb.CompanySettings{
DefaultChartOfAccounts:      c.Settings.DefaultChartOfAccounts,
AllowIntercompanyTransactions: c.Settings.AllowIntercompanyTxn,
RequireApprovalOver:         c.Settings.RequireApprovalOver.ToProto(),
AutoPostingRules:            rules,
PeriodLockingPolicy:         c.Settings.PeriodLockingPolicy,
ReportingCurrency:           c.Settings.ReportingCurrency,
}
}

return &pb.Company{
Id:              c.ID,
Name:            c.Name,
LegalName:       c.LegalName,
TaxId:           c.TaxID,
BaseCurrency:    c.BaseCurrency,
FiscalYearEnd:   timeToProto(c.FiscalYearEnd),
Address:         address,
Settings:        settings,
CreatedAt:       timeToProto(c.CreatedAt),
CreatedBy:       c.CreatedBy,
Status:          status,
ParentCompanyId: c.ParentCompanyID,
Metadata:        metadata,
}
}

func CompanyFromProto(pbCompany *pb.Company) *Company {
if pbCompany == nil {
return nil
}

var status CompanyStatus
switch pbCompany.GetStatus() {
case pb.CompanyStatus_COMPANY_STATUS_ACTIVE:
status = CompanyActive
case pb.CompanyStatus_COMPANY_STATUS_INACTIVE:
status = CompanyInactive
case pb.CompanyStatus_COMPANY_STATUS_SUSPENDED:
status = CompanySuspended
case pb.CompanyStatus_COMPANY_STATUS_MERGED:
status = CompanyMerged
}

// Convert string map to interface map
metadata := make(map[string]interface{})
for k, v := range pbCompany.Metadata {
metadata[k] = v
}

var address *Address
if pbCompany.Address != nil {
address = &Address{
Street1:    pbCompany.Address.Street1,
Street2:    pbCompany.Address.Street2,
City:       pbCompany.Address.City,
State:      pbCompany.Address.State,
PostalCode: pbCompany.Address.PostalCode,
Country:    pbCompany.Address.Country,
}
}

var settings *CompanySettings
if pbCompany.Settings != nil {
var rules []*AutoPostingRule
for _, pbRule := range pbCompany.Settings.AutoPostingRules {
if pbRule != nil {
var actions []*PostingAction
for _, pbAction := range pbRule.Actions {
if pbAction != nil {
actions = append(actions, &PostingAction{
Type:       pbAction.Type,
Parameters: pbAction.Parameters,
})
}
}
rules = append(rules, &AutoPostingRule{
ID:        pbRule.Id,
Name:      pbRule.Name,
Condition: pbRule.Condition,
Actions:   actions,
IsActive:  pbRule.IsActive,
CreatedAt: protoToTime(pbRule.CreatedAt),
})
}
}

settings = &CompanySettings{
DefaultChartOfAccounts: pbCompany.Settings.DefaultChartOfAccounts,
AllowIntercompanyTxn:   pbCompany.Settings.AllowIntercompanyTransactions,
RequireApprovalOver:    AmountFromProto(pbCompany.Settings.RequireApprovalOver),
AutoPostingRules:       rules,
PeriodLockingPolicy:    pbCompany.Settings.PeriodLockingPolicy,
ReportingCurrency:      pbCompany.Settings.ReportingCurrency,
}
}

return &Company{
ID:              pbCompany.Id,
Name:            pbCompany.Name,
LegalName:       pbCompany.LegalName,
TaxID:           pbCompany.TaxId,
BaseCurrency:    pbCompany.BaseCurrency,
FiscalYearEnd:   protoToTime(pbCompany.FiscalYearEnd),
Address:         address,
Settings:        settings,
CreatedAt:       protoToTime(pbCompany.CreatedAt),
CreatedBy:       pbCompany.CreatedBy,
Status:          status,
ParentCompanyID: pbCompany.ParentCompanyId,
Metadata:        metadata,
}
}

func (i *IntercompanyTransaction) ToProto() *pb.IntercompanyTransaction {
if i == nil {
return nil
}

var status pb.IntercompanyStatus
switch i.MatchingStatus {
case IntercompanyPending:
status = pb.IntercompanyStatus_INTERCOMPANY_STATUS_PENDING
case IntercompanyMatched:
status = pb.IntercompanyStatus_INTERCOMPANY_STATUS_MATCHED
case IntercompanyReconciled:
status = pb.IntercompanyStatus_INTERCOMPANY_STATUS_RECONCILED
case IntercompanyDispute:
status = pb.IntercompanyStatus_INTERCOMPANY_STATUS_DISPUTE
default:
status = pb.IntercompanyStatus_INTERCOMPANY_STATUS_UNSPECIFIED
}

return &pb.IntercompanyTransaction{
Id:                  i.ID,
Description:         i.Description,
SourceCompanyId:     i.SourceCompanyID,
TargetCompanyId:     i.TargetCompanyID,
SourceTransactionId: i.SourceTransactionID,
TargetTransactionId: i.TargetTransactionID,
Amount:              i.Amount.ToProto(),
ExchangeRate:        i.ExchangeRate,
MatchingStatus:      status,
CreatedAt:           timeToProto(i.CreatedAt),
CreatedBy:           i.CreatedBy,
ReconciledAt:        optionalTimeToProto(i.ReconciledAt),
ReconciledBy:        i.ReconciledBy,
}
}

func IntercompanyTransactionFromProto(pbTxn *pb.IntercompanyTransaction) *IntercompanyTransaction {
if pbTxn == nil {
return nil
}

var status IntercompanyStatus
switch pbTxn.GetMatchingStatus() {
case pb.IntercompanyStatus_INTERCOMPANY_STATUS_PENDING:
status = IntercompanyPending
case pb.IntercompanyStatus_INTERCOMPANY_STATUS_MATCHED:
status = IntercompanyMatched
case pb.IntercompanyStatus_INTERCOMPANY_STATUS_RECONCILED:
status = IntercompanyReconciled
case pb.IntercompanyStatus_INTERCOMPANY_STATUS_DISPUTE:
status = IntercompanyDispute
}

return &IntercompanyTransaction{
ID:                  pbTxn.Id,
Description:         pbTxn.Description,
SourceCompanyID:     pbTxn.SourceCompanyId,
TargetCompanyID:     pbTxn.TargetCompanyId,
SourceTransactionID: pbTxn.SourceTransactionId,
TargetTransactionID: pbTxn.TargetTransactionId,
Amount:              AmountFromProto(pbTxn.Amount),
ExchangeRate:        pbTxn.ExchangeRate,
MatchingStatus:      status,
CreatedAt:           protoToTime(pbTxn.CreatedAt),
CreatedBy:           pbTxn.CreatedBy,
ReconciledAt:        protoToOptionalTime(pbTxn.ReconciledAt),
ReconciledBy:        pbTxn.ReconciledBy,
}
}

func (c *ConsolidationGroup) ToProto() *pb.ConsolidationGroup {
if c == nil {
return nil
}

return &pb.ConsolidationGroup{
Id:                    c.ID,
Name:                  c.Name,
ParentCompanyId:       c.ParentCompany,
SubsidiaryIds:         c.ChildCompanies,
ConsolidationCurrency: "", // Currency not directly in Go struct
CreatedAt:             timeToProto(c.CreatedAt),
CreatedBy:             c.CreatedBy,
}
}

func ConsolidationGroupFromProto(pbGroup *pb.ConsolidationGroup) *ConsolidationGroup {
if pbGroup == nil {
return nil
}

return &ConsolidationGroup{
ID:                  pbGroup.Id,
Name:                pbGroup.Name,
ParentCompany:       pbGroup.ParentCompanyId,
ChildCompanies:      pbGroup.SubsidiaryIds,
ConsolidationMethod: "FULL", // Default method
CreatedAt:           protoToTime(pbGroup.CreatedAt),
CreatedBy:           pbGroup.CreatedBy,
}
}


// ====================================================================================
// Zero-Based Budgeting Conversions
// ====================================================================================

func (b *BudgetPeriod) ToProto() *pb.BudgetPeriod {
if b == nil {
return nil
}
var status pb.BudgetPeriodStatus
switch b.Status {
case BudgetPeriodDraft:
status = pb.BudgetPeriodStatus_BUDGET_PERIOD_STATUS_DRAFT
case BudgetPeriodOpen:
status = pb.BudgetPeriodStatus_BUDGET_PERIOD_STATUS_OPEN
case BudgetPeriodSubmitted:
status = pb.BudgetPeriodStatus_BUDGET_PERIOD_STATUS_SUBMITTED
case BudgetPeriodApproved:
status = pb.BudgetPeriodStatus_BUDGET_PERIOD_STATUS_APPROVED
case BudgetPeriodLocked:
status = pb.BudgetPeriodStatus_BUDGET_PERIOD_STATUS_LOCKED
}
return &pb.BudgetPeriod{
Id:        b.ID,
Name:      b.Name,
StartDate: timeToProto(b.StartDate),
EndDate:   timeToProto(b.EndDate),
Status:    status,
CreatedAt: timeToProto(b.CreatedAt),
CreatedBy: b.CreatedBy,
}
}

func BudgetPeriodFromProto(pbPeriod *pb.BudgetPeriod) *BudgetPeriod {
if pbPeriod == nil {
return nil
}
var status BudgetPeriodStatus
switch pbPeriod.GetStatus() {
case pb.BudgetPeriodStatus_BUDGET_PERIOD_STATUS_DRAFT:
status = BudgetPeriodDraft
case pb.BudgetPeriodStatus_BUDGET_PERIOD_STATUS_OPEN:
status = BudgetPeriodOpen
case pb.BudgetPeriodStatus_BUDGET_PERIOD_STATUS_SUBMITTED:
status = BudgetPeriodSubmitted
case pb.BudgetPeriodStatus_BUDGET_PERIOD_STATUS_APPROVED:
status = BudgetPeriodApproved
case pb.BudgetPeriodStatus_BUDGET_PERIOD_STATUS_LOCKED:
status = BudgetPeriodLocked
}
return &BudgetPeriod{
ID:        pbPeriod.Id,
Name:      pbPeriod.Name,
StartDate: protoToTime(pbPeriod.StartDate),
EndDate:   protoToTime(pbPeriod.EndDate),
Status:    status,
CreatedAt: protoToTime(pbPeriod.CreatedAt),
CreatedBy: pbPeriod.CreatedBy,
}
}

func (b *BudgetRequest) ToProto() *pb.BudgetRequest {
if b == nil {
return nil
}
var status pb.BudgetRequestStatus
switch b.Status {
case BudgetRequestDraft:
status = pb.BudgetRequestStatus_BUDGET_REQUEST_STATUS_DRAFT
case BudgetRequestSubmitted:
status = pb.BudgetRequestStatus_BUDGET_REQUEST_STATUS_SUBMITTED
case BudgetRequestUnderReview:
status = pb.BudgetRequestStatus_BUDGET_REQUEST_STATUS_UNDER_REVIEW
case BudgetRequestApproved:
status = pb.BudgetRequestStatus_BUDGET_REQUEST_STATUS_APPROVED
case BudgetRequestRejected:
status = pb.BudgetRequestStatus_BUDGET_REQUEST_STATUS_REJECTED
case BudgetRequestRevisionRequired:
status = pb.BudgetRequestStatus_BUDGET_REQUEST_STATUS_REVISION_REQUIRED
}
var lineItems []*pb.BudgetLineItem
for _, item := range b.LineItems {
var priority pb.Priority
switch item.Priority {
case PriorityCritical:
priority = pb.Priority_PRIORITY_CRITICAL
case PriorityHigh:
priority = pb.Priority_PRIORITY_HIGH
case PriorityMedium:
priority = pb.Priority_PRIORITY_MEDIUM
case PriorityLow:
priority = pb.Priority_PRIORITY_LOW
}
lineItems = append(lineItems, &pb.BudgetLineItem{
Id:            item.ID,
AccountId:     item.AccountID,
AccountName:   item.AccountName,
Amount:        item.Amount.ToProto(),
Description:   item.Description,
Priority:      priority,
Recurring:     item.Recurring,
Frequency:     item.Frequency,
Vendor:        item.Vendor,
Dimensions:    DimensionsToProto(item.Dimensions),
Justification: item.Justification,
})
}
var justifications []*pb.Justification
for _, j := range b.Justifications {
var category pb.JustificationCategory
switch j.Category {
case JustificationOperational:
category = pb.JustificationCategory_JUSTIFICATION_CATEGORY_OPERATIONAL
case JustificationStrategic:
category = pb.JustificationCategory_JUSTIFICATION_CATEGORY_STRATEGIC
case JustificationCompliance:
category = pb.JustificationCategory_JUSTIFICATION_CATEGORY_COMPLIANCE
case JustificationGrowth:
category = pb.JustificationCategory_JUSTIFICATION_CATEGORY_GROWTH
case JustificationMaintenance:
category = pb.JustificationCategory_JUSTIFICATION_CATEGORY_MAINTENANCE
}
justifications = append(justifications, &pb.Justification{
Id:                j.ID,
Category:          category,
Title:             j.Title,
Description:       j.Description,
BusinessCase:      j.BusinessCase,
ExpectedOutcome:   j.ExpectedOutcome,
RiskOfNotFunding:  j.RiskOfNotFunding,
SupportingDocs:    j.SupportingDocs,
CreatedAt:         timeToProto(j.CreatedAt),
CreatedBy:         j.CreatedBy,
})
}
return &pb.BudgetRequest{
Id:             b.ID,
PeriodId:       b.PeriodID,
RequestorId:    b.RequestorID,
DepartmentId:   b.DepartmentID,
Title:          b.Title,
Description:    b.Description,
TotalAmount:    b.TotalAmount.ToProto(),
Status:         status,
LineItems:      lineItems,
Justifications: justifications,
CreatedAt:      timeToProto(b.CreatedAt),
UpdatedAt:      timeToProto(b.UpdatedAt),
SubmittedAt:    optionalTimeToProto(b.SubmittedAt),
ApprovedAt:     optionalTimeToProto(b.ApprovedAt),
ApprovedBy:     b.ApprovedBy,
}
}

func BudgetRequestFromProto(pbRequest *pb.BudgetRequest) *BudgetRequest {
if pbRequest == nil {
return nil
}
var status BudgetRequestStatus
switch pbRequest.GetStatus() {
case pb.BudgetRequestStatus_BUDGET_REQUEST_STATUS_DRAFT:
status = BudgetRequestDraft
case pb.BudgetRequestStatus_BUDGET_REQUEST_STATUS_SUBMITTED:
status = BudgetRequestSubmitted
case pb.BudgetRequestStatus_BUDGET_REQUEST_STATUS_UNDER_REVIEW:
status = BudgetRequestUnderReview
case pb.BudgetRequestStatus_BUDGET_REQUEST_STATUS_APPROVED:
status = BudgetRequestApproved
case pb.BudgetRequestStatus_BUDGET_REQUEST_STATUS_REJECTED:
status = BudgetRequestRejected
case pb.BudgetRequestStatus_BUDGET_REQUEST_STATUS_REVISION_REQUIRED:
status = BudgetRequestRevisionRequired
}
var lineItems []BudgetLineItem
for _, pbItem := range pbRequest.LineItems {
var priority Priority
switch pbItem.Priority {
case pb.Priority_PRIORITY_CRITICAL:
priority = PriorityCritical
case pb.Priority_PRIORITY_HIGH:
priority = PriorityHigh
case pb.Priority_PRIORITY_MEDIUM:
priority = PriorityMedium
case pb.Priority_PRIORITY_LOW:
priority = PriorityLow
}
lineItems = append(lineItems, BudgetLineItem{
ID:            pbItem.Id,
AccountID:     pbItem.AccountId,
AccountName:   pbItem.AccountName,
Amount:        AmountFromProto(pbItem.Amount),
Description:   pbItem.Description,
Priority:      priority,
Recurring:     pbItem.Recurring,
Frequency:     pbItem.Frequency,
Vendor:        pbItem.Vendor,
Dimensions:    DimensionsFromProto(pbItem.Dimensions),
Justification: pbItem.Justification,
})
}
var justifications []Justification
for _, pbJ := range pbRequest.Justifications {
var category JustificationCategory
switch pbJ.Category {
case pb.JustificationCategory_JUSTIFICATION_CATEGORY_OPERATIONAL:
category = JustificationOperational
case pb.JustificationCategory_JUSTIFICATION_CATEGORY_STRATEGIC:
category = JustificationStrategic
case pb.JustificationCategory_JUSTIFICATION_CATEGORY_COMPLIANCE:
category = JustificationCompliance
case pb.JustificationCategory_JUSTIFICATION_CATEGORY_GROWTH:
category = JustificationGrowth
case pb.JustificationCategory_JUSTIFICATION_CATEGORY_MAINTENANCE:
category = JustificationMaintenance
}
justifications = append(justifications, Justification{
ID:               pbJ.Id,
Category:         category,
Title:            pbJ.Title,
Description:      pbJ.Description,
BusinessCase:     pbJ.BusinessCase,
ExpectedOutcome:  pbJ.ExpectedOutcome,
RiskOfNotFunding: pbJ.RiskOfNotFunding,
SupportingDocs:   pbJ.SupportingDocs,
CreatedAt:        protoToTime(pbJ.CreatedAt),
CreatedBy:        pbJ.CreatedBy,
})
}
return &BudgetRequest{
ID:             pbRequest.Id,
PeriodID:       pbRequest.PeriodId,
RequestorID:    pbRequest.RequestorId,
DepartmentID:   pbRequest.DepartmentId,
Title:          pbRequest.Title,
Description:    pbRequest.Description,
TotalAmount:    AmountFromProto(pbRequest.TotalAmount),
Status:         status,
LineItems:      lineItems,
Justifications: justifications,
CreatedAt:      protoToTime(pbRequest.CreatedAt),
UpdatedAt:      protoToTime(pbRequest.UpdatedAt),
SubmittedAt:    protoToOptionalTime(pbRequest.SubmittedAt),
ApprovedAt:     protoToOptionalTime(pbRequest.ApprovedAt),
ApprovedBy:     pbRequest.ApprovedBy,
}
}

func (b *BudgetApproval) ToProto() *pb.BudgetApproval {
if b == nil {
return nil
}
var status pb.ApprovalStatus
switch b.Status {
case ApprovalPending:
status = pb.ApprovalStatus_APPROVAL_STATUS_PENDING
case ApprovalApproved:
status = pb.ApprovalStatus_APPROVAL_STATUS_APPROVED
case ApprovalRejected:
status = pb.ApprovalStatus_APPROVAL_STATUS_REJECTED
case ApprovalSkipped:
status = pb.ApprovalStatus_APPROVAL_STATUS_SKIPPED
}
return &pb.BudgetApproval{
Id:             b.ID,
RequestId:      b.RequestID,
ApproverId:     b.ApproverID,
ApproverLevel:  int32(b.ApproverLevel),
Status:         status,
ApprovedAmount: b.ApprovedAmount.ToProto(),
Comments:       b.Comments,
ApprovedAt:     optionalTimeToProto(b.ApprovedAt),
CreatedAt:      timeToProto(b.CreatedAt),
}
}

func BudgetApprovalFromProto(pbApproval *pb.BudgetApproval) *BudgetApproval {
if pbApproval == nil {
return nil
}
var status ApprovalStatus
switch pbApproval.GetStatus() {
case pb.ApprovalStatus_APPROVAL_STATUS_PENDING:
status = ApprovalPending
case pb.ApprovalStatus_APPROVAL_STATUS_APPROVED:
status = ApprovalApproved
case pb.ApprovalStatus_APPROVAL_STATUS_REJECTED:
status = ApprovalRejected
case pb.ApprovalStatus_APPROVAL_STATUS_SKIPPED:
status = ApprovalSkipped
}
return &BudgetApproval{
ID:             pbApproval.Id,
RequestID:      pbApproval.RequestId,
ApproverID:     pbApproval.ApproverId,
ApproverLevel:  int(pbApproval.ApproverLevel),
Status:         status,
ApprovedAmount: AmountFromProto(pbApproval.ApprovedAmount),
Comments:       pbApproval.Comments,
ApprovedAt:     protoToOptionalTime(pbApproval.ApprovedAt),
CreatedAt:      protoToTime(pbApproval.CreatedAt),
}
}

func (b *BudgetAllocation) ToProto() *pb.BudgetAllocation {
if b == nil {
return nil
}
return &pb.BudgetAllocation{
Id:           b.ID,
PeriodId:     b.PeriodID,
RequestId:    b.RequestID,
DepartmentId: b.DepartmentID,
AccountId:    b.AccountID,
Amount:       b.Amount.ToProto(),
SpentAmount:  b.SpentAmount.ToProto(),
Remaining:    b.Remaining.ToProto(),
Description:  b.Description,
Dimensions:   DimensionsToProto(b.Dimensions),
CreatedAt:    timeToProto(b.CreatedAt),
UpdatedAt:    timeToProto(b.UpdatedAt),
}
}

func BudgetAllocationFromProto(pbAllocation *pb.BudgetAllocation) *BudgetAllocation {
if pbAllocation == nil {
return nil
}
return &BudgetAllocation{
ID:           pbAllocation.Id,
PeriodID:     pbAllocation.PeriodId,
RequestID:    pbAllocation.RequestId,
DepartmentID: pbAllocation.DepartmentId,
AccountID:    pbAllocation.AccountId,
Amount:       AmountFromProto(pbAllocation.Amount),
SpentAmount:  AmountFromProto(pbAllocation.SpentAmount),
Remaining:    AmountFromProto(pbAllocation.Remaining),
Description:  pbAllocation.Description,
Dimensions:   DimensionsFromProto(pbAllocation.Dimensions),
CreatedAt:    protoToTime(pbAllocation.CreatedAt),
UpdatedAt:    protoToTime(pbAllocation.UpdatedAt),
}
}

func (b *BudgetTracking) ToProto() *pb.BudgetTracking {
if b == nil {
return nil
}
return &pb.BudgetTracking{
AllocationId:    b.AllocationID,
TransactionId:   b.TransactionID,
Amount:          b.Amount.ToProto(),
Description:     b.Description,
TrackedAt:       timeToProto(b.TrackedAt),
RemainingBudget: b.RemainingBudget.ToProto(),
}
}

func BudgetTrackingFromProto(pbTracking *pb.BudgetTracking) *BudgetTracking {
if pbTracking == nil {
return nil
}
return &BudgetTracking{
AllocationID:    pbTracking.AllocationId,
TransactionID:   pbTracking.TransactionId,
Amount:          AmountFromProto(pbTracking.Amount),
Description:     pbTracking.Description,
TrackedAt:       protoToTime(pbTracking.TrackedAt),
RemainingBudget: AmountFromProto(pbTracking.RemainingBudget),
}
}


// ====================================================================================
// AML Conversions
// ====================================================================================

func (a *AMLRule) ToProto() *pb.AMLRule {
if a == nil {
return nil
}
var ruleType pb.AMLRuleType
switch a.Type {
case RuleCTR:
ruleType = pb.AMLRuleType_AML_RULE_TYPE_CTR
case RuleSAR:
ruleType = pb.AMLRuleType_AML_RULE_TYPE_SAR
default:
ruleType = pb.AMLRuleType_AML_RULE_TYPE_UNSPECIFIED
}
var framework pb.AMLFramework
switch a.Framework {
case BSA_Framework:
framework = pb.AMLFramework_AML_FRAMEWORK_BSA
case AMLD_Framework:
framework = pb.AMLFramework_AML_FRAMEWORK_AMLD
case FATF_Framework:
framework = pb.AMLFramework_AML_FRAMEWORK_FATF
}
return &pb.AMLRule{
Id:          a.ID,
Name:        a.Name,
Type:        ruleType,
Framework:   framework,
Description: a.Description,
Enabled:     a.Enabled,
CreatedAt:   timeToProto(a.CreatedAt),
UpdatedAt:   timeToProto(a.UpdatedAt),
}
}

func AMLRuleFromProto(pbRule *pb.AMLRule) *AMLRule {
if pbRule == nil {
return nil
}
var ruleType AMLRuleType
switch pbRule.GetType() {
case pb.AMLRuleType_AML_RULE_TYPE_CTR:
ruleType = RuleCTR
case pb.AMLRuleType_AML_RULE_TYPE_SAR:
ruleType = RuleSAR
}
var framework AMLFramework
switch pbRule.GetFramework() {
case pb.AMLFramework_AML_FRAMEWORK_BSA:
framework = BSA_Framework
case pb.AMLFramework_AML_FRAMEWORK_AMLD:
framework = AMLD_Framework
case pb.AMLFramework_AML_FRAMEWORK_FATF:
framework = FATF_Framework
}
return &AMLRule{
ID:          pbRule.Id,
Name:        pbRule.Name,
Type:        ruleType,
Framework:   framework,
Description: pbRule.Description,
Enabled:     pbRule.Enabled,
CreatedAt:   protoToTime(pbRule.CreatedAt),
UpdatedAt:   protoToTime(pbRule.UpdatedAt),
}
}

func (a *AMLAlert) ToProto() *pb.AMLAlert {
if a == nil {
return nil
}
var ruleType pb.AMLRuleType
switch a.RuleType {
case RuleCTR:
ruleType = pb.AMLRuleType_AML_RULE_TYPE_CTR
case RuleSAR:
ruleType = pb.AMLRuleType_AML_RULE_TYPE_SAR
}
var framework pb.AMLFramework
switch a.Framework {
case BSA_Framework:
framework = pb.AMLFramework_AML_FRAMEWORK_BSA
case AMLD_Framework:
framework = pb.AMLFramework_AML_FRAMEWORK_AMLD
case FATF_Framework:
framework = pb.AMLFramework_AML_FRAMEWORK_FATF
}
var riskLevel pb.AMLRiskLevel
switch a.RiskLevel {
case RiskLow:
riskLevel = pb.AMLRiskLevel_AML_RISK_LEVEL_LOW
case RiskMedium:
riskLevel = pb.AMLRiskLevel_AML_RISK_LEVEL_MEDIUM
case RiskHigh:
riskLevel = pb.AMLRiskLevel_AML_RISK_LEVEL_HIGH
case RiskCritical:
riskLevel = pb.AMLRiskLevel_AML_RISK_LEVEL_CRITICAL
}
return &pb.AMLAlert{
Id:             a.ID,
RuleType:       ruleType,
Framework:      framework,
RiskLevel:      riskLevel,
Title:          a.Title,
Description:    a.Description,
EntityId:       a.EntityID,
EntityType:     a.EntityType,
TransactionIds: a.TransactionIDs,
AccountIds:     a.AccountIDs,
Amount:         a.Amount.ToProto(),
Currency:       a.Currency,
DetectedAt:     timeToProto(a.DetectedAt),
Status:         a.Status,
AssignedTo:     a.AssignedTo,
CreatedAt:      timeToProto(a.CreatedAt),
UpdatedAt:      timeToProto(a.UpdatedAt),
}
}

func AMLAlertFromProto(pbAlert *pb.AMLAlert) *AMLAlert {
if pbAlert == nil {
return nil
}
var ruleType AMLRuleType
switch pbAlert.GetRuleType() {
case pb.AMLRuleType_AML_RULE_TYPE_CTR:
ruleType = RuleCTR
case pb.AMLRuleType_AML_RULE_TYPE_SAR:
ruleType = RuleSAR
}
var framework AMLFramework
switch pbAlert.GetFramework() {
case pb.AMLFramework_AML_FRAMEWORK_BSA:
framework = BSA_Framework
case pb.AMLFramework_AML_FRAMEWORK_AMLD:
framework = AMLD_Framework
case pb.AMLFramework_AML_FRAMEWORK_FATF:
framework = FATF_Framework
}
var riskLevel AMLRiskLevel
switch pbAlert.GetRiskLevel() {
case pb.AMLRiskLevel_AML_RISK_LEVEL_LOW:
riskLevel = RiskLow
case pb.AMLRiskLevel_AML_RISK_LEVEL_MEDIUM:
riskLevel = RiskMedium
case pb.AMLRiskLevel_AML_RISK_LEVEL_HIGH:
riskLevel = RiskHigh
case pb.AMLRiskLevel_AML_RISK_LEVEL_CRITICAL:
riskLevel = RiskCritical
}
return &AMLAlert{
ID:             pbAlert.Id,
RuleType:       ruleType,
Framework:      framework,
RiskLevel:      riskLevel,
Title:          pbAlert.Title,
Description:    pbAlert.Description,
EntityID:       pbAlert.EntityId,
EntityType:     pbAlert.EntityType,
TransactionIDs: pbAlert.TransactionIds,
AccountIDs:     pbAlert.AccountIds,
Amount:         AmountFromProto(pbAlert.Amount),
Currency:       pbAlert.Currency,
DetectedAt:     protoToTime(pbAlert.DetectedAt),
AssignedTo:     pbAlert.AssignedTo,
CreatedAt:      protoToTime(pbAlert.CreatedAt),
UpdatedAt:      protoToTime(pbAlert.UpdatedAt),
}
}

func (a *AMLCustomer) ToProto() *pb.AMLCustomer {
if a == nil {
return nil
}
var riskLevel pb.AMLRiskLevel
switch a.RiskLevel {
case RiskLow:
riskLevel = pb.AMLRiskLevel_AML_RISK_LEVEL_LOW
case RiskMedium:
riskLevel = pb.AMLRiskLevel_AML_RISK_LEVEL_MEDIUM
case RiskHigh:
riskLevel = pb.AMLRiskLevel_AML_RISK_LEVEL_HIGH
case RiskCritical:
riskLevel = pb.AMLRiskLevel_AML_RISK_LEVEL_CRITICAL
}
return &pb.AMLCustomer{
Id:               a.ID,
CustomerId:       a.CustomerID,
Name:             a.Name,
Type:             a.Type,
RiskLevel:        riskLevel,
Country:          a.Country,
IsPep:            a.IsPEP,
IsHighRisk:       a.IsHighRisk,
SanctionsMatch:   a.SanctionsMatch,
LastKycDate:      optionalTimeToProto(a.LastKYCDate),
LastCddDate:      optionalTimeToProto(a.LastCDDDate),
NextReviewDate:   optionalTimeToProto(a.NextReviewDate),
OnboardingDate:   timeToProto(a.OnboardingDate),
ExpectedActivity: a.ExpectedActivity,
BusinessPurpose:  a.BusinessPurpose,
CreatedAt:        timeToProto(a.CreatedAt),
UpdatedAt:        timeToProto(a.UpdatedAt),
}
}

func AMLCustomerFromProto(pbCustomer *pb.AMLCustomer) *AMLCustomer {
if pbCustomer == nil {
return nil
}
var riskLevel AMLRiskLevel
switch pbCustomer.GetRiskLevel() {
case pb.AMLRiskLevel_AML_RISK_LEVEL_LOW:
riskLevel = RiskLow
case pb.AMLRiskLevel_AML_RISK_LEVEL_MEDIUM:
riskLevel = RiskMedium
case pb.AMLRiskLevel_AML_RISK_LEVEL_HIGH:
riskLevel = RiskHigh
case pb.AMLRiskLevel_AML_RISK_LEVEL_CRITICAL:
riskLevel = RiskCritical
}
return &AMLCustomer{
ID:               pbCustomer.Id,
CustomerID:       pbCustomer.CustomerId,
Name:             pbCustomer.Name,
Type:             pbCustomer.Type,
RiskLevel:        riskLevel,
Country:          pbCustomer.Country,
IsPEP:            pbCustomer.IsPep,
IsHighRisk:       pbCustomer.IsHighRisk,
SanctionsMatch:   pbCustomer.SanctionsMatch,
LastKYCDate:      protoToOptionalTime(pbCustomer.LastKycDate),
LastCDDDate:      protoToOptionalTime(pbCustomer.LastCddDate),
NextReviewDate:   protoToOptionalTime(pbCustomer.NextReviewDate),
OnboardingDate:   protoToTime(pbCustomer.OnboardingDate),
ExpectedActivity: pbCustomer.ExpectedActivity,
BusinessPurpose:  pbCustomer.BusinessPurpose,
CreatedAt:        protoToTime(pbCustomer.CreatedAt),
UpdatedAt:        protoToTime(pbCustomer.UpdatedAt),
}
}

// ====================================================================================
// Compliance Conversions
// ====================================================================================

func (c *ComplianceRule) ToProto() *pb.ComplianceRule {
if c == nil {
return nil
}
var framework pb.ComplianceFramework
switch c.Framework {
case GAAP_Framework:
framework = pb.ComplianceFramework_COMPLIANCE_FRAMEWORK_GAAP
case IFRS_Framework:
framework = pb.ComplianceFramework_COMPLIANCE_FRAMEWORK_IFRS
case SOX_Framework:
framework = pb.ComplianceFramework_COMPLIANCE_FRAMEWORK_SOX
}
var accountType pb.AccountType
switch c.AccountType {
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
}
return &pb.ComplianceRule{
Id:          c.ID,
Framework:   framework,
RuleType:    c.RuleType,
Description: c.Description,
AccountType: accountType,
Conditions:  c.Conditions,
Actions:     c.Actions,
Severity:    c.Severity,
Active:      c.Active,
CreatedAt:   timeToProto(c.CreatedAt),
}
}

func ComplianceRuleFromProto(pbRule *pb.ComplianceRule) *ComplianceRule {
if pbRule == nil {
return nil
}
var framework ComplianceFramework
switch pbRule.GetFramework() {
case pb.ComplianceFramework_COMPLIANCE_FRAMEWORK_GAAP:
framework = GAAP_Framework
case pb.ComplianceFramework_COMPLIANCE_FRAMEWORK_IFRS:
framework = IFRS_Framework
case pb.ComplianceFramework_COMPLIANCE_FRAMEWORK_SOX:
framework = SOX_Framework
}
var accountType AccountType
switch pbRule.GetAccountType() {
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
return &ComplianceRule{
ID:          pbRule.Id,
Framework:   framework,
RuleType:    pbRule.RuleType,
Description: pbRule.Description,
AccountType: accountType,
Conditions:  pbRule.Conditions,
Actions:     pbRule.Actions,
Severity:    pbRule.Severity,
Active:      pbRule.Active,
CreatedAt:   protoToTime(pbRule.CreatedAt),
}
}

func (t *TaxRule) ToProto() *pb.TaxRule {
if t == nil {
return nil
}
var jurisdiction pb.TaxJurisdiction
switch t.Jurisdiction {
case US_FEDERAL:
jurisdiction = pb.TaxJurisdiction_TAX_JURISDICTION_US_FEDERAL
case US_STATE:
jurisdiction = pb.TaxJurisdiction_TAX_JURISDICTION_US_STATE
case EU_VAT:
jurisdiction = pb.TaxJurisdiction_TAX_JURISDICTION_EU_VAT
}
var taxType pb.TaxType
switch t.TaxType {
case INCOME_TAX:
taxType = pb.TaxType_TAX_TYPE_INCOME_TAX
case SALES_TAX:
taxType = pb.TaxType_TAX_TYPE_SALES_TAX
case VAT:
taxType = pb.TaxType_TAX_TYPE_VAT
}
return &pb.TaxRule{
Id:            t.ID,
Jurisdiction:  jurisdiction,
TaxType:       taxType,
Name:          t.Name,
Rate:          t.Rate,
MinAmount:     t.MinAmount,
MaxAmount:     t.MaxAmount,
Exemptions:    t.Exemptions,
EffectiveFrom: timeToProto(t.EffectiveFrom),
EffectiveTo:   optionalTimeToProto(t.EffectiveTo),
Active:        t.Active,
}
}

func TaxRuleFromProto(pbRule *pb.TaxRule) *TaxRule {
if pbRule == nil {
return nil
}
var jurisdiction TaxJurisdiction
switch pbRule.GetJurisdiction() {
case pb.TaxJurisdiction_TAX_JURISDICTION_US_FEDERAL:
jurisdiction = US_FEDERAL
case pb.TaxJurisdiction_TAX_JURISDICTION_US_STATE:
jurisdiction = US_STATE
case pb.TaxJurisdiction_TAX_JURISDICTION_EU_VAT:
jurisdiction = EU_VAT
}
var taxType TaxType
switch pbRule.GetTaxType() {
case pb.TaxType_TAX_TYPE_INCOME_TAX:
taxType = INCOME_TAX
case pb.TaxType_TAX_TYPE_SALES_TAX:
taxType = SALES_TAX
case pb.TaxType_TAX_TYPE_VAT:
taxType = VAT
}
return &TaxRule{
ID:            pbRule.Id,
Jurisdiction:  jurisdiction,
TaxType:       taxType,
Name:          pbRule.Name,
Rate:          pbRule.Rate,
MinAmount:     pbRule.MinAmount,
MaxAmount:     pbRule.MaxAmount,
Exemptions:    pbRule.Exemptions,
EffectiveFrom: protoToTime(pbRule.EffectiveFrom),
EffectiveTo:   protoToOptionalTime(pbRule.EffectiveTo),
Active:        pbRule.Active,
}
}

func (c *ComplianceViolation) ToProto() *pb.ComplianceViolation {
if c == nil {
return nil
}
return &pb.ComplianceViolation{
Id:            c.ID,
RuleId:        c.RuleID,
TransactionId: c.TransactionID,
AccountId:     c.AccountID,
Description:   c.Description,
Severity:      c.Severity,
Status:        c.Status,
DetectedAt:    timeToProto(c.DetectedAt),
ResolvedAt:    optionalTimeToProto(c.ResolvedAt),
Notes:         c.Notes,
}
}

func ComplianceViolationFromProto(pbViolation *pb.ComplianceViolation) *ComplianceViolation {
if pbViolation == nil {
return nil
}
return &ComplianceViolation{
ID:            pbViolation.Id,
RuleID:        pbViolation.RuleId,
TransactionID: pbViolation.TransactionId,
AccountID:     pbViolation.AccountId,
Description:   pbViolation.Description,
Severity:      pbViolation.Severity,
Status:        pbViolation.Status,
DetectedAt:    protoToTime(pbViolation.DetectedAt),
ResolvedAt:    protoToOptionalTime(pbViolation.ResolvedAt),
Notes:         pbViolation.Notes,
}
}

func (t *TaxReturn) ToProto() *pb.TaxReturn {
if t == nil {
return nil
}
var jurisdiction pb.TaxJurisdiction
switch t.Jurisdiction {
case US_FEDERAL:
jurisdiction = pb.TaxJurisdiction_TAX_JURISDICTION_US_FEDERAL
case US_STATE:
jurisdiction = pb.TaxJurisdiction_TAX_JURISDICTION_US_STATE
case EU_VAT:
jurisdiction = pb.TaxJurisdiction_TAX_JURISDICTION_EU_VAT
}
var taxType pb.TaxType
switch t.TaxType {
case INCOME_TAX:
taxType = pb.TaxType_TAX_TYPE_INCOME_TAX
case SALES_TAX:
taxType = pb.TaxType_TAX_TYPE_SALES_TAX
case VAT:
taxType = pb.TaxType_TAX_TYPE_VAT
}
return &pb.TaxReturn{
Id:             t.ID,
Jurisdiction:   jurisdiction,
TaxType:        taxType,
DueDate:        timeToProto(t.DueDate),
TotalTax:       t.TotalTax,
CreatedAt:      timeToProto(t.CreatedAt),
UpdatedAt:      timeToProto(t.UpdatedAt),
}
}

func TaxReturnFromProto(pbReturn *pb.TaxReturn) *TaxReturn {
if pbReturn == nil {
return nil
}
var jurisdiction TaxJurisdiction
switch pbReturn.GetJurisdiction() {
case pb.TaxJurisdiction_TAX_JURISDICTION_US_FEDERAL:
jurisdiction = US_FEDERAL
case pb.TaxJurisdiction_TAX_JURISDICTION_US_STATE:
jurisdiction = US_STATE
case pb.TaxJurisdiction_TAX_JURISDICTION_EU_VAT:
jurisdiction = EU_VAT
}
var taxType TaxType
switch pbReturn.GetTaxType() {
case pb.TaxType_TAX_TYPE_INCOME_TAX:
taxType = INCOME_TAX
case pb.TaxType_TAX_TYPE_SALES_TAX:
taxType = SALES_TAX
case pb.TaxType_TAX_TYPE_VAT:
taxType = VAT
}
return &TaxReturn{
ID:             pbReturn.Id,
Jurisdiction:   jurisdiction,
TaxType:        taxType,
DueDate:        protoToTime(pbReturn.DueDate),
TotalTax:       pbReturn.TotalTax,
CreatedAt:      protoToTime(pbReturn.CreatedAt),
UpdatedAt:      protoToTime(pbReturn.UpdatedAt),
}
}

