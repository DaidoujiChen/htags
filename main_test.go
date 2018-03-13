package main

import "testing"

func TestClearString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{"TestClearString - 1", "hello", "hello"},
		{"TestClearString - 2", " hello", "hello"},
		{"TestClearString - 3", "hello ", "hello"},
		{"TestClearString - 4", " hello ", "hello"},
		{"TestClearString - 5", "", ""},
		{"TestClearString - 6", " ", ""},
		{"TestClearString - 7", "  ", ""},
		{"TestClearString - 8", "   ", ""},
		{"TestClearString - 9", " he llo ", "he llo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if clearString(tt.input) != tt.output {
				t.Fail()
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output bool
	}{
		{"TestIsEmpty - 1", "   ", true},
		{"TestIsEmpty - 2", "hello", false},
		{"TestIsEmpty - 3", "hel lo ", false},
		{"TestIsEmpty - 4", "  f ", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if isEmpty(tt.input) != tt.output {
				t.Fail()
			}
		})
	}
}
