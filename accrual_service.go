package accounting

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AccrualService handles accrual and deferral recognition
type AccrualService struct {
	storage       *Storage
	postingEngine *PostingEngine
	eventStore    *EventStore
}

// NewAccrualService creates a new accrual service
func NewAccrualService(storage *Storage, postingEngine *PostingEngine, eventStore *EventStore) *AccrualService {
	return &AccrualService{
		storage:       storage,
		postingEngine: postingEngine,
		eventStore:    eventStore,
	}
}

// AccrualType represents the type of accrual
type AccrualType string

const (
	AccrualRevenue  AccrualType = "REVENUE"
	AccrualExpense  AccrualType = "EXPENSE"
	DeferralRevenue AccrualType = "DEFERRED_REVENUE"
	DeferralExpense AccrualType = "DEFERRED_EXPENSE"
)

// RecognitionEntry represents a single recognition entry in a schedule
type RecognitionEntry struct {
	ID              string     `json:"id"`
	ScheduleID      string     `json:"schedule_id"`
	PeriodNumber    int        `json:"period_number"`
	RecognitionDate time.Time  `json:"recognition_date"`
	Amount          *Amount    `json:"amount"`
	Status          string     `json:"status"` // "PENDING", "PROCESSED", "FAILED"
	TransactionID   string     `json:"transaction_id,omitempty"`
	ProcessedAt     *time.Time `json:"processed_at,omitempty"`
}

// AccrualTemplate defines the accounts and rules for accrual recognition
type AccrualTemplate struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	AccrualType       AccrualType `json:"accrual_type"`
	AccrualAccountID  string      `json:"accrual_account_id"`
	RevenueAccountID  string      `json:"revenue_account_id,omitempty"`
	ExpenseAccountID  string      `json:"expense_account_id,omitempty"`
	DeferralAccountID string      `json:"deferral_account_id,omitempty"`
	Dimensions        []Dimension `json:"dimensions,omitempty"`
}

// CreateRecognitionSchedule creates a new recognition schedule
func (as *AccrualService) CreateRecognitionSchedule(
	txnID string,
	totalAmount *Amount,
	frequency ScheduleFrequency,
	occurrences int,
	startDate time.Time,
	template *AccrualTemplate,
	userID string,
) (*RecognitionSchedule, error) {

	schedule := &RecognitionSchedule{
		ID:            uuid.New().String(),
		TransactionID: txnID,
		Frequency:     frequency,
		Occurrences:   occurrences,
		StartTime:     startDate,
		CreatedAt:     time.Now(),
	}

	// Save the schedule
	if err := as.storage.SaveSchedule(schedule); err != nil {
		return nil, fmt.Errorf("failed to save schedule: %w", err)
	}

	// Generate recognition entries
	if err := as.generateRecognitionEntries(schedule, totalAmount, template); err != nil {
		return nil, fmt.Errorf("failed to generate recognition entries: %w", err)
	}

	return schedule, nil
}

// generateRecognitionEntries generates all recognition entries for a schedule
func (as *AccrualService) generateRecognitionEntries(
	schedule *RecognitionSchedule,
	totalAmount *Amount,
	template *AccrualTemplate,
) error {

	// Calculate amount per period
	amountPerPeriod := totalAmount.Value / int64(schedule.Occurrences)
	remainder := totalAmount.Value % int64(schedule.Occurrences)

	currentDate := schedule.StartTime

	for i := 0; i < schedule.Occurrences; i++ {
		// Add remainder to last entry to handle rounding
		amount := amountPerPeriod
		if i == schedule.Occurrences-1 {
			amount += remainder
		}

		// Create recognition entry (in a real system, you'd have a separate storage method)
		_ = &RecognitionEntry{
			ID:              uuid.New().String(),
			ScheduleID:      schedule.ID,
			PeriodNumber:    i + 1,
			RecognitionDate: currentDate,
			Amount: &Amount{
				Value:    amount,
				Currency: totalAmount.Currency,
			},
			Status: "PENDING",
		}

		// For now, we'll just store in memory or use a simple approach
		// In a real implementation, you'd save this entry to storage

		// Advance to next period
		currentDate = as.addPeriod(currentDate, schedule.Frequency)
	}

	return nil
}

// addPeriod adds a period to a date based on frequency
func (as *AccrualService) addPeriod(date time.Time, frequency ScheduleFrequency) time.Time {
	switch frequency {
	case Monthly:
		return date.AddDate(0, 1, 0)
	case Quarterly:
		return date.AddDate(0, 3, 0)
	case Yearly:
		return date.AddDate(1, 0, 0)
	default:
		return date.AddDate(0, 1, 0) // Default to monthly
	}
}

// ProcessPendingRecognitions processes all pending recognition entries up to a given date
func (as *AccrualService) ProcessPendingRecognitions(upToDate time.Time, userID string) error {
	schedules, err := as.storage.GetAllSchedules()
	if err != nil {
		return fmt.Errorf("failed to get schedules: %w", err)
	}

	for _, schedule := range schedules {
		if err := as.processSchedule(schedule, upToDate, userID); err != nil {
			// Log error but continue processing other schedules
			fmt.Printf("Error processing schedule %s: %v\n", schedule.ID, err)
		}
	}

	return nil
}

// processSchedule processes a single recognition schedule
func (as *AccrualService) processSchedule(schedule *RecognitionSchedule, upToDate time.Time, userID string) error {
	// Get all pending recognition entries for this schedule
	// In a real implementation, you'd have proper storage methods for recognition entries

	// For this example, we'll generate and process entries on the fly
	currentDate := schedule.StartTime

	for i := 0; i < schedule.Occurrences; i++ {
		if currentDate.After(upToDate) {
			break // Don't process future recognitions
		}

		// Check if this period has already been processed
		// In a real system, you'd check the database

		// Create recognition transaction
		if err := as.createRecognitionTransaction(schedule, currentDate, userID); err != nil {
			return fmt.Errorf("failed to create recognition transaction for period %d: %w", i+1, err)
		}

		currentDate = as.addPeriod(currentDate, schedule.Frequency)
	}

	return nil
}

// createRecognitionTransaction creates a journal entry for accrual/deferral recognition
func (as *AccrualService) createRecognitionTransaction(
	schedule *RecognitionSchedule,
	recognitionDate time.Time,
	userID string,
) error {

	// Get the original transaction to understand the context
	originalTxn, err := as.storage.GetTransaction(schedule.TransactionID)
	if err != nil {
		return fmt.Errorf("failed to get original transaction: %w", err)
	}

	// Create recognition transaction
	recognitionTxn := &Transaction{
		ID:              uuid.New().String(),
		Description:     fmt.Sprintf("Accrual recognition for %s", originalTxn.Description),
		ValidTime:       recognitionDate,
		TransactionTime: time.Now(),
		Status:          Pending,
		SourceRef:       fmt.Sprintf("ACCRUAL_%s", schedule.ID),
		UserID:          userID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Calculate recognition amount (simplified - in practice this would be more complex)
	totalAmount := int64(0)
	for _, entry := range originalTxn.Entries {
		totalAmount += entry.Amount.Value
	}
	recognitionAmount := totalAmount / int64(schedule.Occurrences)

	// Create recognition entries based on accrual type
	// This is a simplified example - real implementation would be more sophisticated

	// Example: Revenue recognition
	debitEntry := Entry{
		ID:            uuid.New().String(),
		TransactionID: recognitionTxn.ID,
		AccountID:     "unearned_revenue", // Deferred revenue account
		Type:          Debit,
		Amount: Amount{
			Value:    recognitionAmount,
			Currency: originalTxn.Entries[0].Amount.Currency,
		},
	}

	creditEntry := Entry{
		ID:            uuid.New().String(),
		TransactionID: recognitionTxn.ID,
		AccountID:     "revenue", // Revenue account
		Type:          Credit,
		Amount: Amount{
			Value:    recognitionAmount,
			Currency: originalTxn.Entries[0].Amount.Currency,
		},
	}

	recognitionTxn.Entries = []Entry{debitEntry, creditEntry}

	// Create transaction creation event
	_, err = as.eventStore.CreateEvent(
		EventCreateTransaction,
		TransactionCreatedEvent{Transaction: recognitionTxn},
		recognitionTxn.ValidTime,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction event: %w", err)
	}

	// Post the transaction
	if err := as.postingEngine.PostTransaction(recognitionTxn, userID); err != nil {
		return fmt.Errorf("failed to post recognition transaction: %w", err)
	}

	return nil
}

// GetScheduleStatus returns the status of a recognition schedule
func (as *AccrualService) GetScheduleStatus(scheduleID string) (*ScheduleStatus, error) {
	schedule, err := as.getScheduleByID(scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	status := &ScheduleStatus{
		ScheduleID:          scheduleID,
		TotalOccurrences:    schedule.Occurrences,
		ProcessedCount:      0,                    // Would be calculated from stored recognition entries
		RemainingCount:      schedule.Occurrences, // Would be calculated
		NextRecognitionDate: schedule.StartTime,   // Would be calculated based on processed entries
	}

	return status, nil
}

// getScheduleByID retrieves a schedule by ID (helper method)
func (as *AccrualService) getScheduleByID(scheduleID string) (*RecognitionSchedule, error) {
	// In a real implementation, you'd have a proper storage method
	schedules, err := as.storage.GetAllSchedules()
	if err != nil {
		return nil, err
	}

	for _, schedule := range schedules {
		if schedule.ID == scheduleID {
			return schedule, nil
		}
	}

	return nil, fmt.Errorf("schedule not found: %s", scheduleID)
}

// ScheduleStatus represents the current status of a recognition schedule
type ScheduleStatus struct {
	ScheduleID          string    `json:"schedule_id"`
	TotalOccurrences    int       `json:"total_occurrences"`
	ProcessedCount      int       `json:"processed_count"`
	RemainingCount      int       `json:"remaining_count"`
	NextRecognitionDate time.Time `json:"next_recognition_date"`
	CompletionRate      float64   `json:"completion_rate"`
}
