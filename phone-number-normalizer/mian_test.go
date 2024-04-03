package main

import "testing"

func TestNormalize(t *testing.T) {
	tests := [...][2]string{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1234567892"},
		{"(123) 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"123-456-7890", "1234567890"},
		{"1234567892", "1234567892"},
		{"(123)456-7892", "1234567892"},
	}

	for _, args := range tests {
		if got := normalizePhoneNumber(args[0]); got != args[1] {
			t.Errorf("Expected %v got %v", args[1], got)
		}
	}
}
