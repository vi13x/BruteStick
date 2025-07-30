package test

import (
	"bruteforce/internal/brute"
	"testing"
)

func TestNextPassword(t *testing.T) {
	alphabet := "abc"
	tests := []struct {
		current  string
		expected string
		done     bool
	}{
		{"a", "b", false},
		{"c", "aa", true}, // переход к следующей длине (завершение)
		{"ab", "ac", false},
		{"ac", "ba", false},
		{"cc", "aaa", true},
	}

	for _, tt := range tests {
		got, done := brute.nextPassword(tt.current, alphabet)
		if got != tt.expected || done != tt.done {
			t.Errorf("nextPassword(%q) = %q, %v; want %q, %v", tt.current, got, done, tt.expected, tt.done)
		}
	}
}
