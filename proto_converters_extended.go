package accounting

import (
	pb "accounting/proto/accounting"
	"google.golang.org/protobuf/proto"
)

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
