package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"unicode/utf8"
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
			false,
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
	content := []byte("temporary file's content")
	dir, err := ioutil.TempDir("", "testing")
	if err != nil {
		t.Fatalf("Error creating tempdir: \n%d", err)
	}

	tmpfile, err := ioutil.TempFile(dir, "test")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
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
		errorPrefix := fmt.Sprintf("Test [%d]: ", i)
		actual := GenerateName(test.length)
		if utf8.RuneCountInString(actual) != test.length {
			t.Errorf(errorPrefix+"Expected length of %d, got %d (%d)", test.length, utf8.RuneCountInString(actual), actual)
		}
	}
}

func TestGetFileExtFromDir(t *testing.T) {
	err := os.Mkdir("testing", 0777)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.OpenFile("testing/png.png", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.OpenFile("testing/jpeg.jpeg", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		dir      string
		expected string
	}{
		{
			"png",
			"testing/",
			".png",
		},
		{
			"jpeg",
			"testing/",
			".jpeg",
		},
		{
			"gif",
			"testing/",
			"",
		},
	}

	for i, test := range tests {
		errorPrefix := fmt.Sprintf("Test [%d]: ", i)
		actual := GetFileExtFromDir(test.name, test.dir)
		if actual != test.expected {
			t.Errorf(errorPrefix+"Expected %d, got %d", test.expected, actual)
		}
	}

	os.RemoveAll("testing/png.png")
	os.RemoveAll("testing/jpeg.jpeg")
	err = os.RemoveAll("testing")
	if err != nil {
		t.Fatal(err)
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
		errorPrefix := fmt.Sprintf("Test [%d]: ", i)
		actual := GetFileExt(test.string)
		if actual != test.expected {
			t.Errorf(errorPrefix+"Expected %d, got %d", test.expected, actual)
		}
	}
}
