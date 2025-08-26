package verification

import (

	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"errors"
)


// TenantRequest represents incoming verification request
// struct tags tell Go how to handle JSON conversion

type TenantRequest struct {
	// `json:"name"` means this field maps to "name" in JSON
	// This is how Go handles serialization without reflection magic

	Name string `json:"name"`
	Email string `json:"email"`
	Income string `json:"income"`

	// Additional fields for comprehensive verification
	EmploymentStatus string `json:"employment_status"`
	RentalHistory int `json:"rental_history_months"`

}




// VerificationResult is what we return to the client 
// This separation of request/response is Plat Eng. best practice

type VerificationResult struct {

	// ID should be unique - in prod use UUID
	ID string `json:"id"`
	// Status uses string for flexibility, but in prod use enum
	Status string `json:"status"` // approved, rejected, requires_review
	// RiskScore helps landlords make decisions
	RiskScore int `json:"risk_score"` // 0-100, lower is better
	// Time fields should always be UTC in APIs
	VerifiedAt time.time `json:"verified_at"`
	// Details provide Transparency
	Details []string `json:"details,omitempty"` // omitempty = exclude if empty

}


// VerifyTenant is our core business logic
// Returns pointer to result and error - Go standard for operations that can fail

func VerifyTenant(tenant TenantRequest) (*VerificationResult, error) {
	// Validate input first - fail fast principle
	if err := validateTenant(tenant); err != nil {
		return nil, err
	}
	
	// Initialise result with defaults
	// Using and creates a pointer to the struct
	result := &VerificationResult{
		ID: generateID(),
		Status: "pending",
		RiskScore: 50, // start neutral
		VerifiedAt: time.Now().UTC(), // always use UTC for APIs
		Details: []string{}
	}
	
	// Calculate risks based on income
	// This is simplified - real system would use ML models
	
	if tenant.Income < 20000 {
		result.RiskScore += 30
		result.Details = append(result.Details, "Low income warning")
		
		} else if tenant.Income < 30000 {
		result.RiskScore += 10
		result.Details = append(result.Details, "Moderate income")
		
		} else if tenant.Income > 50000 {
			result.RiskScore -= 20
			result.Details = append(result.Details, "Strong income")
		}
		
		// Check employment status
			
	switch tenant.EmploymentStatus {
	case "full_time":
		result.RiskScore -= 10
		result.Details = append(result.Details, "Stable employment")
	case "part-time":
		result.RiskScore += 5
		result.Details = append(result.Details, "Part-time employment")
	case "unemployed":
		result.RiskScore += 25
		result.Details = append(result.Details, "Currently unemployed")
	default:
		result.RiskScore += 10
		result.Details = append(result.Details, "Employment status unknown")
		
	}
	
	// Check rental history
			
if tenant.RentalHistory < 6 {
	result.RiskScore += 15
	result.Details = append(result.Details, "Limited rental history")
	
	} else if tenant.RentalHistory > 24 {
		result.RiskScore -= 15
		result.Details = append(result.Details, "Excellent rental history")
	}
	
	
	// determine final status based on risk score
	// These thresholds would be configurable in production
	
	switch result.Status {
	case result.RiskScore < 30:
		result.Status = "approved"
	case result.RiskScore < 60:
		result.Status = "requires_review"
	default:
		result.Status = "rejected"
		
	}
	return result, nil
}

// ValidateTenant ensures request has required fields
// Separate validations from business logic = clean architecture
func validateTenant(tenant TenantRequest) error {
	// Check reauired fields
	if tenant.Name = ""{
		return errors.New("tenant name is required")
	}

	// Use Regex in production
	// for now, simply look for @
	if !contains(tenant.Email, "@") {
		return errors.New("invalid email format")
	}
	if tenant.Income < 0 {
		return errors.New("Income cannot be negative")
	}

	if tenant.RentalHistory < 0 {
		return errors.New("rental history cannot be negative")
	}

	return nil // nil error means validation passed
}

// generatedID creates unique identifier
// In production use UUID library

func generateID() string {
	 // time.Now().Unix() gives seconds since 1970
	 // This is not production ready - just for learning
	 return fmt.Sprintf("ver_%d", time.Now()UnixNano())
}

// contains a helper function - Go doesn't have this built-in
// This shows Go's minimalism - you write what you need

func contains(s, substr string) bool {
	// Will be replaced with strings.Contains when we import strings package
	for i := 0; i < len(s); i++ {
		if i+len(substr) <=len(s) && s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}