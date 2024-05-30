package eflag

import (
	"reflect"
	"testing"
)

// TestSplitWithComma tests the SplitWithComma function.
func TestSplitWithComma(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"a,b,c", []string{"a", "b", "c"}},
		{" a , b , c ", []string{"a", "b", "c"}},
		{"a,b , c", []string{"a", "b", "c"}},
		{"", []string{""}},
		{"a,", []string{"a", ""}},
	}

	for _, test := range tests {
		result := SplitWithComma(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("SplitWithComma(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}

// TestSplitWith tests the SplitWith function with various separators.
func TestSplitWith(t *testing.T) {
	tests := []struct {
		input    string
		sep      string
		expected []string
	}{
		{"a,b,c", ",", []string{"a", "b", "c"}},
		{" a | b | c ", "|", []string{"a", "b", "c"}},
		{"a;b ; c", ";", []string{"a", "b", "c"}},
		{"a:b:c", ":", []string{"a", "b", "c"}},
		{"", ",", []string{""}},
		{"a,", ",", []string{"a", ""}},
	}

	for _, test := range tests {
		result := SplitWith(test.input, test.sep)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("SplitWith(%q, %q) = %v; want %v", test.input, test.sep, result, test.expected)
		}
	}
}

// TestMixedCapsToScreamingSnake tests the MixedCapsToScreamingSnake function.
func TestMixedCapsToScreamingSnake(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"mixedCaps", "MIXED_CAPS"},
		{"MixedCaps", "MIXED_CAPS"},
		{"mixedCapsString", "MIXED_CAPS_STRING"},
		{"MixedCapsString", "MIXED_CAPS_STRING"},
		{"simpleTest", "SIMPLE_TEST"},
		{"SimpleTest", "SIMPLE_TEST"},
		{"already_screaming", "ALREADY_SCREAMING"},
		{"AlreadyScreaming", "ALREADY_SCREAMING"},
		{"simple", "SIMPLE"},
		{"Simple", "SIMPLE"},
	}

	for _, test := range tests {
		result := MixedCapsToScreamingSnake(test.input)
		if result != test.expected {
			t.Errorf("MixedCapsToScreamingSnake(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
