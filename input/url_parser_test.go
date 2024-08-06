package input

import "testing"

func TestValidateUrlWithValidUrl(t *testing.T) {
	err := ValidateUrl(nil, []string{"https://m8n.dev"})
	if err != nil {
		t.Error("Error should be nil")
	}
}

func TestValidateUrlWithInvalidUrl(t *testing.T) {
	err := ValidateUrl(nil, []string{"foobar"})
	if err == nil {
		t.Error("Error should not be nil")
	}
}

func TestValidateUrlWithRelative(t *testing.T) {
	err := ValidateUrl(nil, []string{"/this/is/relative"})
	if err == nil {
		t.Error("Error should not be nil")
	}
}
