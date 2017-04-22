package util

import (
	"fmt"
	"os"
	"testing"
)

func TestContains(t *testing.T) {
	tests := []struct {
		slice    []string
		string   string
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
			t.Errorf(errorPrefix+"Expected %d, got %d", test.expected, actual)
		}
	}
}

func TestCheckExists(t *testing.T) {
	const dir string = "testing/"
	err := os.Mkdir(dir, 0777)
	if err != nil {
		t.Fatalf("Error creating dir: \n%d", err)
	}

	_, err = os.OpenFile(dir+"test.txt", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		t.Fatalf("Error creating test file: \n%d", err)
	}

	tests := []struct {
		filename string
		expected bool
	}{
		{
			"test",
			true,
		},
		{
			"foo",
			false,
		},
	}

	for i, test := range tests {
		errorPrefix := fmt.Sprintf("Test [%d]: ", i)
		actual := CheckExists(test.filename, dir)
		if actual != test.expected {
			t.Errorf(errorPrefix+"Expected %d, got %d", test.expected, actual)
		}
	}

	err = os.Remove(dir)
	if err != nil {
		t.Fatalf("Error removing dir: \n%d", err)
	}
}