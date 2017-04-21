package util

import (
	"testing"
	"fmt"
)

func TestSliceContains(t *testing.T) {
	tests := []struct {
		slice []string
		string string
		expected bool
	}{
		{
			[]string{
				"foo",
				"bar",
				"foobar",
				"barfoo",
			},
			"foo",
			true,
		},
		{
			[]string{
				"foo",
				"bar",
				"foobar",
				"barfoo",
			},
			"test",
			true,
		},
	}

	for i, test := range tests {
		errorPrefix := fmt.Sprintf("Test [%d]: ", i)
		actual := Contains(test.slice, test.string)
		if actual != test.expected {
			t.Errorf(errorPrefix + "Expected %d, got %d", test.expected, actual)
		}
	}
}
