package post

import (
	"testing"

	apptemplate "devlog-studio/backend/internal/template"
)

func TestValidateAlgorithmRequiresCoreFields(t *testing.T) {
	err := validateAlgorithm(TypeAlgorithm, apptemplate.AlgorithmInput{})
	validationErr, ok := IsValidationError(err)
	if !ok {
		t.Fatalf("expected validation error, got %T", err)
	}
	if validationErr.Code != "VALIDATION_ERROR" {
		t.Fatalf("Code = %q, want VALIDATION_ERROR", validationErr.Code)
	}
}

func TestValidateAlgorithmAcceptsValidInput(t *testing.T) {
	err := validateAlgorithm(TypeAlgorithm, apptemplate.AlgorithmInput{
		Platform:     "Programmers",
		ProblemTitle: "등굣길",
		ProblemURL:   "https://school.programmers.co.kr/example",
		Language:     "Java",
		Approach:     "DP",
		Code:         "class Solution {}",
	})
	if err != nil {
		t.Fatalf("validateAlgorithm() error = %v", err)
	}
}
