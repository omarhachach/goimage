package util

import (
	"fmt"
	"os"
	"testing"
	"unicode/utf8"
)

var dir = "testing/"

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
			false,
		},
	}

	for i, test := range tests {
		errorPrefix := fmt.Sprintf("Test [%v]: ", i)
		actual := Contains(test.slice, test.string)
		if actual != test.expected {
			t.Errorf(errorPrefix+"Expected %v, got %v", test.expected, actual)
		}
	}
}

func TestCheckExists(t *testing.T) {
	err := os.Mkdir(dir, 0777)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile(dir+"test.txt", os.O_RDWR|os.O_CREATE, 0777)
	f.Close()
	if err != nil {
		t.Fatal(err)
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
		errorPrefix := fmt.Sprintf("Test [%v]: ", i)
		actual := CheckExists(test.filename, dir)
		if actual != test.expected {
			t.Errorf(errorPrefix+"Expected %v, got %v", test.expected, actual)
		}
	}
	err = os.RemoveAll(dir)
	if err != nil {
		t.Error(err)
	}
}

func TestGenerateName(t *testing.T) {
	tests := []struct {
		length int
	}{
		{
			4,
		},
		{
			6,
		},
		{
			9,
		},
	}

	for i, test := range tests {
		errorPrefix := fmt.Sprintf("Test [%v]: ", i)
		actual := GenerateName(test.length)
		if utf8.RuneCountInString(actual) != test.length {
			t.Errorf(errorPrefix+"Expected length of %v, got %v (%v)", test.length, utf8.RuneCountInString(actual), actual)
		}
	}
}

func TestGetFileExtFromDir(t *testing.T) {
	err := os.Mkdir(dir, 0777)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile(dir+"png.png", os.O_RDONLY|os.O_CREATE, 0777)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	f, err = os.OpenFile(dir+"jpeg.jpeg", os.O_RDONLY|os.O_CREATE, 0777)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		expected string
	}{
		{
			"png",
			".png",
		},
		{
			"jpeg",
			".jpeg",
		},
		{
			"gif",
			"",
		},
	}

	for i, test := range tests {
		errorPrefix := fmt.Sprintf("Test [%v]: ", i)
		actual := GetFileExtFromDir(test.name, dir)
		if actual != test.expected {
			t.Errorf(errorPrefix+"Expected %v, got %v", test.expected, actual)
		}
	}

	err = os.RemoveAll(dir)
	if err != nil {
		t.Error(err)
	}
}

func TestGetFileExt(t *testing.T) {
	tests := []struct {
		string   string
		expected string
	}{
		{
			"form-data; name=\"file\"; filename=\"bluelogo.png\"",
			".png",
		},
		{
			"form-data; name=\"file\"; filename=\"bluelogo.jpeg\"",
			".jpeg",
		},
	}

	for i, test := range tests {
		errorPrefix := fmt.Sprintf("Test [%v]: ", i)
		actual := GetFileExt(test.string)
		if actual != test.expected {
			t.Errorf(errorPrefix+"Expected %v, got %v", test.expected, actual)
		}
	}
}
