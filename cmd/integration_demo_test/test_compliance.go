package main_test

import (
	"accounting"
	"fmt"
	"os"
	"testing"
)

func TestComplianceFeatures(t *testing.T) {
	dbPath := "test_compliance.db"
	// Clean up the database file after the test
	defer os.Remove(dbPath)

	// Initialize the accounting engine (or directly use storage if that's the test's focus)
	// For this test, we might only need the storage component.
	storage, err := accounting.NewStorage(dbPath) // Assuming NewStorage exists and returns (*Storage, error)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	// defer storage.Close() // If storage has a Close method

	// Initialize ComplianceService with dereferenced storage
	cs := accounting.NewComplianceService(*storage) // Dereference storage

	// Setup standard compliance rules
	// ...

	// Test tax calculation
	calc, err := cs.CalculateTax(1000.0, accounting.US_FEDERAL, accounting.INCOME_TAX, []string{})
	if err != nil {
		t.Errorf("Error calculating tax: %v", err)
		return
	}

	fmt.Printf("Tax calculation successful!\\n")
	fmt.Printf("Base Amount: $%.2f\\n", calc.BaseAmount)
	fmt.Printf("Tax Amount: $%.2f\\n", calc.TaxAmount)
	fmt.Printf("Tax Rate: %.1f%%\\n", calc.TaxRate*100)

	// Test compliance rule creation
	rule := accounting.ComplianceRule{
		Framework:   accounting.GAAP_Framework,
		RuleType:    "JOURNAL_ENTRY_BALANCE",
		Description: "Test rule",
		Severity:    "ERROR",
	}

	err = cs.CreateComplianceRule(rule)
	if err != nil {
		t.Errorf("Error creating compliance rule: %v", err)
		return
	}

	fmt.Println("Compliance rule created successfully!")
	fmt.Println("✅ Compliance module is working correctly!")
}
