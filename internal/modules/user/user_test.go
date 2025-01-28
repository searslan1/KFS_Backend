package user

import "testing"

func TestValidateEmail(t *testing.T) {
	email := "test@example.com"
	if err := ValidateEmail(email); err != nil {
		t.Errorf("Email validation failed for valid email: %s", email)
	}
}
