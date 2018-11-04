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

func TestGetFileBasename(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name         string
		args         args
		wantBasename string
	}{
		{
			name: "check filename with single period",
			args: args{
				filename: "hd739x.jpg",
			},
			wantBasename: "hd739x",
		},
		{
			name: "check filename with two periods",
			args: args{
				filename: "hd739x.file.jpg",
			},
			wantBasename: "hd739x.file",
		},
		{
			name: "check filename with no periods",
			args: args{
				filename: "hd739x",
			},
			wantBasename: "hd739x",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotBasename := GetFileBasename(tt.args.filename); gotBasename != tt.wantBasename {
				t.Errorf("GetFileBasename() = %v, want %v", gotBasename, tt.wantBasename)
			}
		})
	}
}
