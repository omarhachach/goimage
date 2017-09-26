package utils

import (
	"fmt"
	"testing"
)

func TestContainsString(t *testing.T) {
	testSlice := []string{"hi", "this", "is", "a", "test"}
	tests := []struct {
		TestString string
		Expected   bool
	}{
		{
			TestString: "hi",
			Expected:   true,
		},
		{
			TestString: "no",
			Expected:   false,
		},
	}

	for i, test := range tests {
		errorPrefix := fmt.Sprintf("Error [%v]: ", i)
		actual := ContainsString(test.TestString, testSlice)
		if actual != test.Expected {
			t.Errorf(errorPrefix+"Expected %v, got %v.", test.Expected, actual)
		}
	}
}

func TestGenerateName(t *testing.T) {
	tests := []struct {
		Length int
	}{
		{
			Length: 4,
		},
		{
			Length: 8,
		},
		{
			Length: 12,
		},
	}

	for i, test := range tests {
		errorPrefix := fmt.Sprintf("Error [%v]: ", i)
		actual := len(GenerateName(test.Length))
		if actual != test.Length {
			t.Errorf(errorPrefix+"Expected len of %v, got %v.", test.Length, actual)
		}
	}
}
