package goimage

import (
	"fmt"
	"mime/multipart"
	"net/textproto"
	"os"
	"testing"
)

func TestGetFileBasename(t *testing.T) {
	tests := []struct {
		Filename string
		Expected string
	}{
		{
			Filename: "testname.png",
			Expected: "testname",
		},
		{
			Filename: "noname",
			Expected: "noname",
		},
	}

	for i, test := range tests {
		errorPrefix := fmt.Sprintf("Error [%v]: ", i)
		actual := GetFileBasename(test.Filename)
		if actual != test.Expected {
			t.Errorf(errorPrefix+"Expected %v, got %v.", test.Expected, actual)
		}
	}
}

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		Header   *multipart.FileHeader
		Expected string
	}{
		{
			Header: &multipart.FileHeader{
				Filename: "testname.png",
			},
			Expected: "png",
		},
		{
			Header: &multipart.FileHeader{
				Filename: "testname.jpeg",
			},
			Expected: "jpeg",
		},
	}

	for i, test := range tests {
		errorPrefix := fmt.Sprintf("Error [%v]: ", i)
		actual := GetFileExtension(test.Header)
		if actual != test.Expected {
			t.Errorf(errorPrefix+"Expected %v, got %v.", test.Expected, actual)
		}
	}
}

func TestGetFileMIMEType(t *testing.T) {
	tests := []struct {
		Header   textproto.MIMEHeader
		Expected string
	}{
		{
			Header: textproto.MIMEHeader{
				"Content-Type": []string{"image/png"},
			},
			Expected: "image/png",
		},
		{
			Header: textproto.MIMEHeader{
				"Content-Type": []string{"image/jpeg"},
			},
			Expected: "image/jpeg",
		},
	}

	for i, test := range tests {
		errorPrefix := fmt.Sprintf("Error [%v]: ", i)
		actual := GetFileMIMEType(test.Header)
		if actual != test.Expected {
			t.Errorf(errorPrefix+"Expected %v, got %v.", test.Expected, actual)
		}
	}
}

func TestMoveFile(t *testing.T) {
	tempDir := "goimagetest/"

	err := os.Mkdir(tempDir, 0644)
	if err != nil {
		t.Error("Error creating temp directory: " + err.Error())
		return
	}

	file, err := os.OpenFile(tempDir+"test.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Error("Error creating test file: " + err.Error())
		return
	}

	MoveFile(file, tempDir+"newTest.txt")

	file.Close()

	newFile, err := os.Open(tempDir + "newTest.txt")
	if err != nil {
		t.Error("Error on opening moved file: " + err.Error())
		return
	}

	newFile.Close()

	err = os.RemoveAll(tempDir)
	if err != nil {
		t.Error("Error deleting temp directory: " + err.Error())
		return
	}
}
