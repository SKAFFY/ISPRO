package main

import (
	"regexp"
	"testing"

	"github.com/dlclark/regexp2"
)

var (
	testSimplePattern   = regexp.MustCompile(`^[^:]+:[^:]+:[^:]+$`)
	testComplexPassword = regexp2.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*])[a-zA-Z\d!@#$%^&*]{8,}$`, 0)
)

func ValidateSimpleLinkTest(link string) bool {
	return testSimplePattern.MatchString(link)
}

func ValidateComplexPasswordTest(password string) bool {
	match, _ := testComplexPassword.MatchString(password)
	return match
}

func TestValidateSimpleLink(t *testing.T) {
	tests := []struct {
		name     string
		link     string
		expected bool
	}{
		{"valid three fields", "source:target:description", true},
		{"valid simple", "a:b:c", true},
		{"invalid only one field", "invalid", false},
		{"invalid four fields", "a:b:c:d", false},
		{"empty string", "", false},
		{"valid with spaces", "my source:my target:my desc", true},
		{"valid with numbers", "123:456:789", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateSimpleLinkTest(tt.link)
			if result != tt.expected {
				t.Errorf("ValidateSimpleLink(%q) = %v, want %v", tt.link, result, tt.expected)
			}
		})
	}
}

func TestValidateComplexPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"valid password", "Password1!", true},
		{"valid complex", "MyP@ssw0rd", true},
		{"too short", "Pass1!", false},
		{"no lowercase", "PASSWORD1!", false},
		{"no uppercase", "password1!", false},
		{"no digit", "Password!", false},
		{"no special char", "Password1", false},
		{"empty string", "", false},
		{"valid with all chars", "Aa1!bbbb", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateComplexPasswordTest(tt.password)
			if result != tt.expected {
				t.Errorf("ValidateComplexPassword(%q) = %v, want %v", tt.password, result, tt.expected)
			}
		})
	}
}

func BenchmarkValidateSimpleLink(b *testing.B) {
	link := "source:target:description"
	for i := 0; i < b.N; i++ {
		_ = testSimplePattern.MatchString(link)
	}
}

func BenchmarkValidateComplexPassword(b *testing.B) {
	password := "Password1!"
	for i := 0; i < b.N; i++ {
		_, _ = testComplexPassword.MatchString(password)
	}
}
