package accounting

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZeroBasedBudgeting(t *testing.T) {
	// Setup
	dbFile := "test_zbb.db"
	defer os.Remove(dbFile)

	engine, err := NewAccountingEngine(dbFile)
	require.NoError(t, err)
	defer engine.Close()

	userID := "budget_manager"
	departmentID := "marketing"

	// Create standard accounts
	err = engine.CreateStandardAccounts(userID)
	require.NoError(t, err)

	t.Run("Budget Period Management", func(t *testing.T) {
		// Create budget period
		period := &BudgetPeriod{
			Name:      "FY2025 Budget",
			StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
		}

		err = engine.CreateBudgetPeriod(period, userID)
		require.NoError(t, err)
		assert.NotEmpty(t, period.ID)
		assert.Equal(t, BudgetPeriodDraft, period.Status)
		assert.Equal(t, userID, period.CreatedBy)
	})

	t.Run("Budget Request Creation and Justification", func(t *testing.T) {
		// Create budget period first
		period := &BudgetPeriod{
			Name:      "Q1-2025 Budget",
			StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 3, 31, 23, 59, 59, 0, time.UTC),
		}
		err = engine.CreateBudgetPeriod(period, userID)
		require.NoError(t, err)

		// Create budget request with line items
		request := &BudgetRequest{
			PeriodID:     period.ID,
			DepartmentID: departmentID,
			Title:        "Marketing Campaign Budget",
			Description:  "Zero-based budget request for digital marketing campaigns",
			LineItems: []BudgetLineItem{
				{
					ID:            "line1",
					AccountID:     "expenses",
					AccountName:   "Marketing Expenses",
					Amount:        &Amount{Value: 50000, Currency: "USD"}, // $500
					Description:   "Google Ads Campaign",
					Priority:      PriorityHigh,
					Recurring:     true,
					Frequency:     "Monthly",
					Vendor:        "Google",
					Justification: "Increase brand awareness and drive lead generation",
				},
				{
					ID:            "line2",
					AccountID:     "expenses",
					AccountName:   "Marketing Expenses",
					Amount:        &Amount{Value: 30000, Currency: "USD"}, // $300
					Description:   "Content Creation Tools",
					Priority:      PriorityMedium,
					Recurring:     false,
					Vendor:        "Adobe",
					Justification: "Professional content creation for marketing materials",
				},
			},
		}

		err = engine.CreateBudgetRequest(request, userID)
		require.NoError(t, err)
		assert.NotEmpty(t, request.ID)
		assert.Equal(t, BudgetRequestDraft, request.Status)
		assert.Equal(t, userID, request.RequestorID)
		assert.Equal(t, int64(80000), request.TotalAmount.Value) // $800 total

		// Add detailed justification
		justification := &Justification{
			Category:         JustificationGrowth,
			Title:            "Digital Marketing Expansion",
			Description:      "Comprehensive digital marketing strategy to capture Q1 growth opportunities",
			BusinessCase:     "Market research shows 40% increase in digital engagement. We need to capitalize on this trend.",
			ExpectedOutcome:  "25% increase in qualified leads, 15% increase in conversion rate",
			RiskOfNotFunding: "Loss of market share to competitors who are investing in digital marketing",
			Alternatives: []Alternative{
				{
					Description: "Reduce scope to Google Ads only",
					Cost:        &Amount{Value: 50000, Currency: "USD"},
					Pros:        []string{"Lower cost", "Focused approach"},
					Cons:        []string{"Limited reach", "No content creation capability"},
				},
				{
					Description: "Outsource to marketing agency",
					Cost:        &Amount{Value: 120000, Currency: "USD"},
					Pros:        []string{"Professional expertise", "Full service"},
					Cons:        []string{"Higher cost", "Less control"},
				},
			},
			Metrics: []JustificationMetric{
				{
					Name:        "Lead Generation",
					Value:       500,
					Unit:        "leads",
					Timeframe:   "Q1 2025",
					Measurement: "CRM tracking and attribution",
				},
				{
					Name:        "Cost per Acquisition",
					Value:       160,
					Unit:        "USD",
					Timeframe:   "Q1 2025",
					Measurement: "Marketing analytics dashboard",
				},
			},
		}

		err = engine.AddJustification(request.ID, justification, userID)
		require.NoError(t, err)

		// Submit budget request
		err = engine.SubmitBudgetRequest(request.ID, userID)
		require.NoError(t, err)
	})

	t.Run("Budget Approval Workflow", func(t *testing.T) {
		// Create budget period
		period := &BudgetPeriod{
			Name:      "Q2-2025 Budget",
			StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 6, 30, 23, 59, 59, 0, time.UTC),
		}
		err = engine.CreateBudgetPeriod(period, userID)
		require.NoError(t, err)

		// Create and submit budget request
		request := &BudgetRequest{
			PeriodID:     period.ID,
			DepartmentID: departmentID,
			Title:        "Sales Tools Budget",
			Description:  "Essential sales enablement tools",
			LineItems: []BudgetLineItem{
				{
					ID:            "salestools1",
					AccountID:     "expenses",
					AccountName:   "Software Expenses",
					Amount:        &Amount{Value: 25000, Currency: "USD"}, // $250
					Description:   "CRM Software License",
					Priority:      PriorityCritical,
					Recurring:     true,
					Frequency:     "Monthly",
					Justification: "Critical for sales pipeline management and customer tracking",
				},
			},
		}

		err = engine.CreateBudgetRequest(request, userID)
		require.NoError(t, err)

		err = engine.SubmitBudgetRequest(request.ID, userID)
		require.NoError(t, err)

		// Approve budget request
		approverID := "finance_director"
		approvedAmount := &Amount{Value: 25000, Currency: "USD"}
		comments := "Approved as requested. Critical for sales operations."

		err = engine.ApproveBudgetRequest(request.ID, approverID, approvedAmount, comments)
		require.NoError(t, err)

		// Create budget allocation
		err = engine.CreateBudgetAllocation(request.ID, userID)
		require.NoError(t, err)
	})

	t.Run("Budget Tracking and Variance Analysis", func(t *testing.T) {
		// Create budget period
		period := &BudgetPeriod{
			Name:      "Q3-2025 Budget",
			StartDate: time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2025, 9, 30, 23, 59, 59, 0, time.UTC),
		}
		err = engine.CreateBudgetPeriod(period, userID)
		require.NoError(t, err)

		// Create, submit, and approve budget request
		request := &BudgetRequest{
			PeriodID:     period.ID,
			DepartmentID: departmentID,
			Title:        "IT Infrastructure Budget",
			Description:  "Server and network infrastructure upgrades",
			LineItems: []BudgetLineItem{
				{
					ID:            "it1",
					AccountID:     "expenses",
					AccountName:   "IT Expenses",
					Amount:        &Amount{Value: 100000, Currency: "USD"}, // $1000
					Description:   "Server Hardware",
					Priority:      PriorityHigh,
					Recurring:     false,
					Justification: "Replace aging servers to improve performance and reliability",
				},
			},
		}

		err = engine.CreateBudgetRequest(request, userID)
		require.NoError(t, err)

		err = engine.SubmitBudgetRequest(request.ID, userID)
		require.NoError(t, err)

		err = engine.ApproveBudgetRequest(request.ID, "cto", request.TotalAmount, "Approved for infrastructure upgrade")
		require.NoError(t, err)

		err = engine.CreateBudgetAllocation(request.ID, userID)
		require.NoError(t, err)

		// Create actual spending transaction
		transaction := &Transaction{
			Description: "Server Purchase",
			ValidTime:   time.Now(),
			Entries: []Entry{
				{
					AccountID: "expenses",
					Type:      Debit,
					Amount:    Amount{Value: 80000, Currency: "USD"}, // $800 actual spend
				},
				{
					AccountID: "cash",
					Type:      Credit,
					Amount:    Amount{Value: 80000, Currency: "USD"},
				},
			},
		}

		err = engine.CreateTransaction(transaction, userID)
		require.NoError(t, err)

		err = engine.PostTransaction(transaction.ID, userID)
		require.NoError(t, err)

		// Note: In a real implementation, we'd need to link the allocation ID
		// For this test, we'll create a mock allocation ID and track spending
		// This would normally be done automatically or through a UI

		// Generate variance report
		varianceReport, err := engine.GetBudgetVariance(period.ID, departmentID)
		require.NoError(t, err)
		assert.Equal(t, period.ID, varianceReport.PeriodID)
		assert.Equal(t, departmentID, varianceReport.DepartmentID)

		// Get department budget summary
		summary, err := engine.GetDepartmentBudgetSummary(period.ID, departmentID)
		require.NoError(t, err)
		assert.Equal(t, period.ID, summary.PeriodID)
		assert.Equal(t, departmentID, summary.DepartmentID)
		assert.NotEmpty(t, summary.Requests)
		assert.True(t, summary.TotalRequested.Value > 0)
	})

	t.Run("Comprehensive ZBB Workflow", func(t *testing.T) {
		// Create annual budget period
		period := &BudgetPeriod{
			Name:      "FY2026 Annual Budget",
			StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC),
		}
		err = engine.CreateBudgetPeriod(period, userID)
		require.NoError(t, err)

		// Create comprehensive budget request with multiple priorities
		request := &BudgetRequest{
			PeriodID:     period.ID,
			DepartmentID: "operations",
			Title:        "Operations Annual Budget - Zero Based",
			Description:  "Complete operational budget built from zero base with full justification",
			LineItems: []BudgetLineItem{
				{
					ID:            "ops1",
					AccountID:     "expenses",
					AccountName:   "Operational Expenses",
					Amount:        &Amount{Value: 500000, Currency: "USD"}, // $5000
					Description:   "Staff Training and Development",
					Priority:      PriorityCritical,
					Recurring:     true,
					Frequency:     "Quarterly",
					Justification: "Mandatory compliance training and skill development to maintain operational efficiency",
					Dimensions: []Dimension{
						{Key: DimDepartment, Value: "operations"},
						{Key: DimCostCenter, Value: "ops_training"},
					},
				},
				{
					ID:            "ops2",
					AccountID:     "expenses",
					AccountName:   "Operational Expenses",
					Amount:        &Amount{Value: 300000, Currency: "USD"}, // $3000
					Description:   "Process Automation Tools",
					Priority:      PriorityHigh,
					Recurring:     false,
					Justification: "Automation tools to reduce manual work and improve accuracy",
					Dimensions: []Dimension{
						{Key: DimDepartment, Value: "operations"},
						{Key: DimProject, Value: "automation_initiative"},
					},
				},
				{
					ID:            "ops3",
					AccountID:     "expenses",
					AccountName:   "Operational Expenses",
					Amount:        &Amount{Value: 150000, Currency: "USD"}, // $1500
					Description:   "Team Building Activities",
					Priority:      PriorityLow,
					Recurring:     true,
					Frequency:     "Semi-annually",
					Justification: "Improve team morale and collaboration",
					Dimensions: []Dimension{
						{Key: DimDepartment, Value: "operations"},
					},
				},
			},
		}

		err = engine.CreateBudgetRequest(request, userID)
		require.NoError(t, err)
		assert.Equal(t, int64(950000), request.TotalAmount.Value) // $9500 total

		// Add comprehensive justification with strategic alignment
		strategicJustification := &Justification{
			Category:    JustificationStrategic,
			Title:       "Strategic Operations Enhancement",
			Description: "Comprehensive operational improvements aligned with company growth strategy",
			BusinessCase: "Operations must scale efficiently to support 50% growth target. " +
				"Training ensures compliance and quality. Automation reduces errors and frees staff for strategic work. " +
				"Team building maintains culture during rapid growth.",
			ExpectedOutcome: "30% increase in operational efficiency, 50% reduction in manual errors, " +
				"95% compliance rate, 20% improvement in employee satisfaction",
			RiskOfNotFunding: "Inability to scale operations, compliance violations, " +
				"increased error rates, staff burnout, loss of key personnel",
			Alternatives: []Alternative{
				{
					Description: "Hire additional staff instead of automation",
					Cost:        &Amount{Value: 800000, Currency: "USD"},
					Pros:        []string{"Immediate capacity", "No technology risk"},
					Cons:        []string{"Higher long-term cost", "Scalability issues", "Human error risk"},
				},
				{
					Description: "Outsource operations",
					Cost:        &Amount{Value: 1200000, Currency: "USD"},
					Pros:        []string{"Expert management", "Scalable"},
					Cons:        []string{"Loss of control", "Higher cost", "Security concerns"},
				},
			},
			Metrics: []JustificationMetric{
				{
					Name:        "Operational Efficiency",
					Value:       30,
					Unit:        "percent improvement",
					Timeframe:   "FY2026",
					Measurement: "Process time tracking and KPI dashboard",
				},
				{
					Name:        "Error Rate Reduction",
					Value:       50,
					Unit:        "percent reduction",
					Timeframe:   "Q2 2026",
					Measurement: "Quality assurance metrics",
				},
				{
					Name:        "Employee Satisfaction",
					Value:       4.2,
					Unit:        "score out of 5",
					Timeframe:   "Annual survey",
					Measurement: "Anonymous employee satisfaction survey",
				},
			},
		}

		err = engine.AddJustification(request.ID, strategicJustification, userID)
		require.NoError(t, err)

		// Submit for approval
		err = engine.SubmitBudgetRequest(request.ID, userID)
		require.NoError(t, err)

		// Approve with partial amount (demonstrating budget negotiation)
		approvedAmount := &Amount{Value: 800000, Currency: "USD"} // $8000 (approved 84%)
		err = engine.ApproveBudgetRequest(request.ID, "coo", approvedAmount,
			"Approved training and automation. Team building deferred to Q3.")
		require.NoError(t, err)

		// Create allocations
		err = engine.CreateBudgetAllocation(request.ID, userID)
		require.NoError(t, err)

		// Get final summary
		summary, err := engine.GetDepartmentBudgetSummary(period.ID, "operations")
		require.NoError(t, err)

		assert.Equal(t, 1, len(summary.Requests))
		assert.Equal(t, BudgetRequestApproved, summary.Requests[0].Status)
		assert.Equal(t, int64(950000), summary.TotalRequested.Value)
		assert.NotNil(t, summary.StatusCounts)
		assert.Equal(t, 1, summary.StatusCounts[BudgetRequestApproved])
	})
}

func TestZBBValidationAndErrorHandling(t *testing.T) {
	// Setup
	dbFile := "test_zbb_validation.db"
	defer os.Remove(dbFile)

	engine, err := NewAccountingEngine(dbFile)
	require.NoError(t, err)
	defer engine.Close()

	userID := "test_user"

	t.Run("Submit Request Without Justifications", func(t *testing.T) {
		// Create period
		period := &BudgetPeriod{
			Name:      "Test Period",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 3, 0),
		}
		err = engine.CreateBudgetPeriod(period, userID)
		require.NoError(t, err)

		// Create request without line item justifications
		request := &BudgetRequest{
			PeriodID:     period.ID,
			DepartmentID: "test_dept",
			Title:        "Test Request",
			LineItems: []BudgetLineItem{
				{
					ID:          "test1",
					AccountID:   "expenses",
					Amount:      &Amount{Value: 10000, Currency: "USD"},
					Description: "Test Item",
					// No justification provided
				},
			},
		}

		err = engine.CreateBudgetRequest(request, userID)
		require.NoError(t, err)

		// Attempt to submit should fail due to missing justification
		err = engine.SubmitBudgetRequest(request.ID, userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing justification")
	})

	t.Run("Approve Non-Submitted Request", func(t *testing.T) {
		// Create period
		period := &BudgetPeriod{
			Name:      "Test Period 2",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 3, 0),
		}
		err = engine.CreateBudgetPeriod(period, userID)
		require.NoError(t, err)

		// Create request but don't submit
		request := &BudgetRequest{
			PeriodID:     period.ID,
			DepartmentID: "test_dept",
			Title:        "Test Request 2",
			LineItems: []BudgetLineItem{
				{
					ID:            "test2",
					AccountID:     "expenses",
					Amount:        &Amount{Value: 10000, Currency: "USD"},
					Description:   "Test Item",
					Justification: "Test justification",
				},
			},
		}

		err = engine.CreateBudgetRequest(request, userID)
		require.NoError(t, err)

		// Attempt to approve draft request should fail
		err = engine.ApproveBudgetRequest(request.ID, "approver",
			&Amount{Value: 10000, Currency: "USD"}, "Test approval")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "can only approve submitted requests")
	})

	t.Run("Allocate Non-Approved Request", func(t *testing.T) {
		// Create period
		period := &BudgetPeriod{
			Name:      "Test Period 3",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 3, 0),
		}
		err = engine.CreateBudgetPeriod(period, userID)
		require.NoError(t, err)

		// Create and submit request but don't approve
		request := &BudgetRequest{
			PeriodID:     period.ID,
			DepartmentID: "test_dept",
			Title:        "Test Request 3",
			LineItems: []BudgetLineItem{
				{
					ID:            "test3",
					AccountID:     "expenses",
					Amount:        &Amount{Value: 10000, Currency: "USD"},
					Description:   "Test Item",
					Justification: "Test justification",
				},
			},
		}

		err = engine.CreateBudgetRequest(request, userID)
		require.NoError(t, err)

		err = engine.SubmitBudgetRequest(request.ID, userID)
		require.NoError(t, err)

		// Attempt to allocate without approval should fail
		err = engine.CreateBudgetAllocation(request.ID, userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "can only allocate approved requests")
	})
}
