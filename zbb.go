package accounting

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ----------------------------------------------------------------------------
// Zero-Based Budgeting (ZBB) Structures
// ----------------------------------------------------------------------------

// BudgetPeriod represents a budget cycle (e.g., fiscal year, quarter)
type BudgetPeriod struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"` // e.g., "FY2025", "Q1-2025"
	StartDate time.Time          `json:"start_date"`
	EndDate   time.Time          `json:"end_date"`
	Status    BudgetPeriodStatus `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	CreatedBy string             `json:"created_by"`
}

type BudgetPeriodStatus string

const (
	BudgetPeriodDraft     BudgetPeriodStatus = "DRAFT"
	BudgetPeriodOpen      BudgetPeriodStatus = "OPEN"
	BudgetPeriodSubmitted BudgetPeriodStatus = "SUBMITTED"
	BudgetPeriodApproved  BudgetPeriodStatus = "APPROVED"
	BudgetPeriodLocked    BudgetPeriodStatus = "LOCKED"
)

// BudgetRequest represents a zero-based budget request for a cost center/department
type BudgetRequest struct {
	ID             string              `json:"id"`
	PeriodID       string              `json:"period_id"`
	RequestorID    string              `json:"requestor_id"`  // User who created the request
	DepartmentID   string              `json:"department_id"` // Department/cost center
	Title          string              `json:"title"`
	Description    string              `json:"description"`
	TotalAmount    *Amount             `json:"total_amount"`
	Status         BudgetRequestStatus `json:"status"`
	LineItems      []BudgetLineItem    `json:"line_items"`
	Justifications []Justification     `json:"justifications"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	SubmittedAt    *time.Time          `json:"submitted_at,omitempty"`
	ApprovedAt     *time.Time          `json:"approved_at,omitempty"`
	ApprovedBy     string              `json:"approved_by,omitempty"`
}

type BudgetRequestStatus string

const (
	BudgetRequestDraft            BudgetRequestStatus = "DRAFT"
	BudgetRequestSubmitted        BudgetRequestStatus = "SUBMITTED"
	BudgetRequestUnderReview      BudgetRequestStatus = "UNDER_REVIEW"
	BudgetRequestApproved         BudgetRequestStatus = "APPROVED"
	BudgetRequestRejected         BudgetRequestStatus = "REJECTED"
	BudgetRequestRevisionRequired BudgetRequestStatus = "REVISION_REQUIRED"
)

// BudgetLineItem represents individual expense items within a budget request
type BudgetLineItem struct {
	ID            string      `json:"id"`
	AccountID     string      `json:"account_id"` // GL account for this expense
	AccountName   string      `json:"account_name"`
	Amount        *Amount     `json:"amount"`
	Description   string      `json:"description"`
	Priority      Priority    `json:"priority"`            // CRITICAL, HIGH, MEDIUM, LOW
	Recurring     bool        `json:"recurring"`           // Is this a recurring expense?
	Frequency     string      `json:"frequency,omitempty"` // Monthly, Quarterly, etc.
	Vendor        string      `json:"vendor,omitempty"`
	Dimensions    []Dimension `json:"dimensions,omitempty"`
	Justification string      `json:"justification"` // Why this expense is needed
}

type Priority string

const (
	PriorityCritical Priority = "CRITICAL"
	PriorityHigh     Priority = "HIGH"
	PriorityMedium   Priority = "MEDIUM"
	PriorityLow      Priority = "LOW"
)

// Justification represents detailed business justification for budget items
type Justification struct {
	ID               string                `json:"id"`
	Category         JustificationCategory `json:"category"`
	Title            string                `json:"title"`
	Description      string                `json:"description"`
	BusinessCase     string                `json:"business_case"`       // Why this is needed
	ExpectedOutcome  string                `json:"expected_outcome"`    // What will be achieved
	RiskOfNotFunding string                `json:"risk_of_not_funding"` // What happens if not funded
	Alternatives     []Alternative         `json:"alternatives,omitempty"`
	Metrics          []JustificationMetric `json:"metrics,omitempty"`
	SupportingDocs   []string              `json:"supporting_docs,omitempty"`
	CreatedAt        time.Time             `json:"created_at"`
	CreatedBy        string                `json:"created_by"`
}

type JustificationCategory string

const (
	JustificationOperational JustificationCategory = "OPERATIONAL"
	JustificationStrategic   JustificationCategory = "STRATEGIC"
	JustificationCompliance  JustificationCategory = "COMPLIANCE"
	JustificationGrowth      JustificationCategory = "GROWTH"
	JustificationMaintenance JustificationCategory = "MAINTENANCE"
)

// Alternative represents alternative approaches or cost options
type Alternative struct {
	Description string   `json:"description"`
	Cost        *Amount  `json:"cost"`
	Pros        []string `json:"pros"`
	Cons        []string `json:"cons"`
}

// JustificationMetric represents measurable outcomes from budget spend
type JustificationMetric struct {
	Name        string  `json:"name"`        // e.g., "Cost Savings", "Revenue Increase"
	Value       float64 `json:"value"`       // Expected value
	Unit        string  `json:"unit"`        // e.g., "USD", "hours", "customers"
	Timeframe   string  `json:"timeframe"`   // When this will be achieved
	Measurement string  `json:"measurement"` // How this will be measured
}

// BudgetApproval represents approval workflow for budget requests
type BudgetApproval struct {
	ID             string         `json:"id"`
	RequestID      string         `json:"request_id"`
	ApproverID     string         `json:"approver_id"`
	ApproverLevel  int            `json:"approver_level"` // 1=Manager, 2=Director, 3=VP, etc.
	Status         ApprovalStatus `json:"status"`
	ApprovedAmount *Amount        `json:"approved_amount,omitempty"`
	Comments       string         `json:"comments,omitempty"`
	ApprovedAt     *time.Time     `json:"approved_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
}

type ApprovalStatus string

const (
	ApprovalPending  ApprovalStatus = "PENDING"
	ApprovalApproved ApprovalStatus = "APPROVED"
	ApprovalRejected ApprovalStatus = "REJECTED"
	ApprovalSkipped  ApprovalStatus = "SKIPPED"
)

// BudgetAllocation represents final approved budget allocation
type BudgetAllocation struct {
	ID           string      `json:"id"`
	PeriodID     string      `json:"period_id"`
	RequestID    string      `json:"request_id"`
	DepartmentID string      `json:"department_id"`
	AccountID    string      `json:"account_id"`
	Amount       *Amount     `json:"amount"`
	SpentAmount  *Amount     `json:"spent_amount"`
	Remaining    *Amount     `json:"remaining"`
	Description  string      `json:"description"`
	Dimensions   []Dimension `json:"dimensions,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// BudgetTracking represents budget vs actual tracking
type BudgetTracking struct {
	AllocationID    string    `json:"allocation_id"`
	TransactionID   string    `json:"transaction_id"`
	Amount          *Amount   `json:"amount"`
	Description     string    `json:"description"`
	TrackedAt       time.Time `json:"tracked_at"`
	RemainingBudget *Amount   `json:"remaining_budget"`
}

// ----------------------------------------------------------------------------
// Zero-Based Budgeting Service
// ----------------------------------------------------------------------------

type ZBBService struct {
	storage *Storage
}

func NewZBBService(storage *Storage) *ZBBService {
	return &ZBBService{
		storage: storage,
	}
}

// CreateBudgetPeriod creates a new budget period
func (zbb *ZBBService) CreateBudgetPeriod(period *BudgetPeriod, userID string) error {
	if period.ID == "" {
		period.ID = uuid.New().String()
	}
	period.CreatedAt = time.Now()
	period.CreatedBy = userID
	period.Status = BudgetPeriodDraft

	return zbb.storage.SaveBudgetPeriod(period)
}

// CreateBudgetRequest creates a new zero-based budget request
func (zbb *ZBBService) CreateBudgetRequest(request *BudgetRequest, userID string) error {
	if request.ID == "" {
		request.ID = uuid.New().String()
	}

	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()
	request.RequestorID = userID
	request.Status = BudgetRequestDraft

	// Calculate total amount from line items
	total := &Amount{Value: 0, Currency: "USD"}
	for _, item := range request.LineItems {
		if item.Amount != nil {
			total.Value += item.Amount.Value
		}
	}
	request.TotalAmount = total

	return zbb.storage.SaveBudgetRequest(request)
}

// AddJustification adds business justification to a budget request
func (zbb *ZBBService) AddJustification(requestID string, justification *Justification, userID string) error {
	request, err := zbb.storage.GetBudgetRequest(requestID)
	if err != nil {
		return fmt.Errorf("failed to get budget request: %w", err)
	}

	if justification.ID == "" {
		justification.ID = uuid.New().String()
	}
	justification.CreatedAt = time.Now()
	justification.CreatedBy = userID

	request.Justifications = append(request.Justifications, *justification)
	request.UpdatedAt = time.Now()

	return zbb.storage.SaveBudgetRequest(request)
}

// SubmitBudgetRequest submits request for approval
func (zbb *ZBBService) SubmitBudgetRequest(requestID string, userID string) error {
	request, err := zbb.storage.GetBudgetRequest(requestID)
	if err != nil {
		return fmt.Errorf("failed to get budget request: %w", err)
	}

	if request.Status != BudgetRequestDraft {
		return fmt.Errorf("can only submit draft requests")
	}

	// Validate that all line items have justifications
	for _, item := range request.LineItems {
		if item.Justification == "" {
			return fmt.Errorf("line item '%s' missing justification", item.Description)
		}
	}

	request.Status = BudgetRequestSubmitted
	now := time.Now()
	request.SubmittedAt = &now
	request.UpdatedAt = now

	return zbb.storage.SaveBudgetRequest(request)
}

// ApproveBudgetRequest approves a budget request
func (zbb *ZBBService) ApproveBudgetRequest(requestID string, approverID string, approvedAmount *Amount, comments string) error {
	request, err := zbb.storage.GetBudgetRequest(requestID)
	if err != nil {
		return fmt.Errorf("failed to get budget request: %w", err)
	}

	if request.Status != BudgetRequestSubmitted && request.Status != BudgetRequestUnderReview {
		return fmt.Errorf("can only approve submitted requests")
	}

	// Create approval record
	approval := &BudgetApproval{
		ID:             uuid.New().String(),
		RequestID:      requestID,
		ApproverID:     approverID,
		ApproverLevel:  1, // Simplified for demo
		Status:         ApprovalApproved,
		ApprovedAmount: approvedAmount,
		Comments:       comments,
		CreatedAt:      time.Now(),
		ApprovedAt:     &[]time.Time{time.Now()}[0],
	}

	err = zbb.storage.SaveBudgetApproval(approval)
	if err != nil {
		return fmt.Errorf("failed to save approval: %w", err)
	}

	// Update request status
	request.Status = BudgetRequestApproved
	now := time.Now()
	request.ApprovedAt = &now
	request.ApprovedBy = approverID
	request.UpdatedAt = now

	return zbb.storage.SaveBudgetRequest(request)
}

// CreateBudgetAllocation creates budget allocation from approved request
func (zbb *ZBBService) CreateBudgetAllocation(requestID string, userID string) error {
	request, err := zbb.storage.GetBudgetRequest(requestID)
	if err != nil {
		return fmt.Errorf("failed to get budget request: %w", err)
	}

	if request.Status != BudgetRequestApproved {
		return fmt.Errorf("can only allocate approved requests")
	}

	// Create allocations for each line item
	for _, item := range request.LineItems {
		allocation := &BudgetAllocation{
			ID:           uuid.New().String(),
			PeriodID:     request.PeriodID,
			RequestID:    requestID,
			DepartmentID: request.DepartmentID,
			AccountID:    item.AccountID,
			Amount:       item.Amount,
			SpentAmount:  &Amount{Value: 0, Currency: item.Amount.Currency},
			Remaining:    &Amount{Value: item.Amount.Value, Currency: item.Amount.Currency},
			Description:  item.Description,
			Dimensions:   item.Dimensions,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err = zbb.storage.SaveBudgetAllocation(allocation)
		if err != nil {
			return fmt.Errorf("failed to save allocation: %w", err)
		}
	}

	return nil
}

// TrackBudgetSpending tracks actual spending against budget allocations
func (zbb *ZBBService) TrackBudgetSpending(transactionID string, allocationID string) error {
	// Get transaction to validate amount
	txn, err := zbb.storage.GetTransaction(transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	allocation, err := zbb.storage.GetBudgetAllocation(allocationID)
	if err != nil {
		return fmt.Errorf("failed to get allocation: %w", err)
	}

	// Calculate spending amount from transaction entries
	var spendAmount int64
	for _, entry := range txn.Entries {
		if entry.AccountID == allocation.AccountID && entry.Type == Debit {
			spendAmount += entry.Amount.Value
		}
	}

	if spendAmount <= 0 {
		return fmt.Errorf("no valid spending found for account %s", allocation.AccountID)
	}

	// Update allocation
	allocation.SpentAmount.Value += spendAmount
	allocation.Remaining.Value -= spendAmount
	allocation.UpdatedAt = time.Now()

	err = zbb.storage.SaveBudgetAllocation(allocation)
	if err != nil {
		return fmt.Errorf("failed to update allocation: %w", err)
	}

	// Create tracking record
	tracking := &BudgetTracking{
		AllocationID:    allocationID,
		TransactionID:   transactionID,
		Amount:          &Amount{Value: spendAmount, Currency: allocation.Amount.Currency},
		Description:     txn.Description,
		TrackedAt:       time.Now(),
		RemainingBudget: &Amount{Value: allocation.Remaining.Value, Currency: allocation.Amount.Currency},
	}

	return zbb.storage.SaveBudgetTracking(tracking)
}

// GetBudgetVariance calculates variance between budget and actual
func (zbb *ZBBService) GetBudgetVariance(periodID string, departmentID string) (*BudgetVarianceReport, error) {
	allocations, err := zbb.storage.GetBudgetAllocationsByPeriodAndDept(periodID, departmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get allocations: %w", err)
	}

	report := &BudgetVarianceReport{
		PeriodID:     periodID,
		DepartmentID: departmentID,
		GeneratedAt:  time.Now(),
		Items:        make([]BudgetVarianceItem, 0),
	}

	var totalBudget, totalSpent, totalVariance int64

	for _, allocation := range allocations {
		variance := allocation.Amount.Value - allocation.SpentAmount.Value
		variancePercent := float64(variance) / float64(allocation.Amount.Value) * 100

		item := BudgetVarianceItem{
			AccountID:       allocation.AccountID,
			Description:     allocation.Description,
			BudgetAmount:    allocation.Amount,
			SpentAmount:     allocation.SpentAmount,
			Variance:        &Amount{Value: variance, Currency: allocation.Amount.Currency},
			VariancePercent: variancePercent,
		}

		report.Items = append(report.Items, item)

		totalBudget += allocation.Amount.Value
		totalSpent += allocation.SpentAmount.Value
		totalVariance += variance
	}

	report.TotalBudget = &Amount{Value: totalBudget, Currency: "USD"}
	report.TotalSpent = &Amount{Value: totalSpent, Currency: "USD"}
	report.TotalVariance = &Amount{Value: totalVariance, Currency: "USD"}
	if totalBudget > 0 {
		report.TotalVariancePercent = float64(totalVariance) / float64(totalBudget) * 100
	}

	return report, nil
}

// BudgetVarianceReport represents budget vs actual variance analysis
type BudgetVarianceReport struct {
	PeriodID             string               `json:"period_id"`
	DepartmentID         string               `json:"department_id"`
	TotalBudget          *Amount              `json:"total_budget"`
	TotalSpent           *Amount              `json:"total_spent"`
	TotalVariance        *Amount              `json:"total_variance"`
	TotalVariancePercent float64              `json:"total_variance_percent"`
	Items                []BudgetVarianceItem `json:"items"`
	GeneratedAt          time.Time            `json:"generated_at"`
}

type BudgetVarianceItem struct {
	AccountID       string  `json:"account_id"`
	Description     string  `json:"description"`
	BudgetAmount    *Amount `json:"budget_amount"`
	SpentAmount     *Amount `json:"spent_amount"`
	Variance        *Amount `json:"variance"`
	VariancePercent float64 `json:"variance_percent"`
}

// GetDepartmentBudgetSummary provides summary of all budget requests for a department
func (zbb *ZBBService) GetDepartmentBudgetSummary(periodID string, departmentID string) (*DepartmentBudgetSummary, error) {
	requests, err := zbb.storage.GetBudgetRequestsByPeriodAndDept(periodID, departmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget requests: %w", err)
	}

	summary := &DepartmentBudgetSummary{
		PeriodID:     periodID,
		DepartmentID: departmentID,
		GeneratedAt:  time.Now(),
		Requests:     make([]BudgetRequestSummary, 0),
	}

	var totalRequested, totalApproved int64
	statusCounts := make(map[BudgetRequestStatus]int)

	for _, request := range requests {
		reqSummary := BudgetRequestSummary{
			ID:                 request.ID,
			Title:              request.Title,
			Status:             request.Status,
			RequestedAmount:    request.TotalAmount,
			LineItemCount:      len(request.LineItems),
			JustificationCount: len(request.Justifications),
			SubmittedAt:        request.SubmittedAt,
			ApprovedAt:         request.ApprovedAt,
		}

		// Get approved amount from approvals
		if request.Status == BudgetRequestApproved {
			approvals, err := zbb.storage.GetBudgetApprovalsByRequest(request.ID)
			if err == nil && len(approvals) > 0 {
				reqSummary.ApprovedAmount = approvals[0].ApprovedAmount
				if reqSummary.ApprovedAmount != nil {
					totalApproved += reqSummary.ApprovedAmount.Value
				}
			}
		}

		summary.Requests = append(summary.Requests, reqSummary)

		totalRequested += request.TotalAmount.Value
		statusCounts[request.Status]++
	}

	summary.TotalRequested = &Amount{Value: totalRequested, Currency: "USD"}
	summary.TotalApproved = &Amount{Value: totalApproved, Currency: "USD"}
	summary.StatusCounts = statusCounts

	return summary, nil
}

type DepartmentBudgetSummary struct {
	PeriodID       string                      `json:"period_id"`
	DepartmentID   string                      `json:"department_id"`
	TotalRequested *Amount                     `json:"total_requested"`
	TotalApproved  *Amount                     `json:"total_approved"`
	StatusCounts   map[BudgetRequestStatus]int `json:"status_counts"`
	Requests       []BudgetRequestSummary      `json:"requests"`
	GeneratedAt    time.Time                   `json:"generated_at"`
}

type BudgetRequestSummary struct {
	ID                 string              `json:"id"`
	Title              string              `json:"title"`
	Status             BudgetRequestStatus `json:"status"`
	RequestedAmount    *Amount             `json:"requested_amount"`
	ApprovedAmount     *Amount             `json:"approved_amount,omitempty"`
	LineItemCount      int                 `json:"line_item_count"`
	JustificationCount int                 `json:"justification_count"`
	SubmittedAt        *time.Time          `json:"submitted_at,omitempty"`
	ApprovedAt         *time.Time          `json:"approved_at,omitempty"`
}
