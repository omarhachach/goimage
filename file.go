package goimage

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

// FileInfo holds info about a file.
type FileInfo struct {
	Basename  string
	Extension string
	Filename  string
}

// GetFileBasename returns the base file name of a given filename.
// Eg. the file name without the extension.
func GetFileBasename(filename string) (basename string) {
	index := strings.LastIndex(filename, ".")
	if index == -1 {
		return filename
	}

	return filename[:index]
}

// GetFileInfo returns the existing files info.
func GetFileInfo(dirname, filename string) (fileInfo *FileInfo, err error) {
	existingFiles, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	for _, existingFile := range existingFiles {
		name := existingFile.Name()
		basename := GetFileBasename(name)

		if basename == GetFileBasename(filename) {
			return &FileInfo{
				Basename:  basename,
				Extension: filepath.Ext(name),
				Filename:  name,
			}, nil
		}
	}

	return nil, nil
}

// GetFileExtension returns the given files extension.
func GetFileExtension(header *multipart.FileHeader) (extension string) {
	return filepath.Ext(header.Filename)[1:]
}

// GetFileMIMEType get an image mime type.
func GetFileMIMEType(header textproto.MIMEHeader) (MIMEType string) {
	return header["Content-Type"][0]
}

// MoveFile copies a file to the given filepath.
func MoveFile(file multipart.File, filepath string) (err error) {
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, file)
	if err != nil {
		newErr := os.Remove(filepath)
		if newErr != nil {
			return newErr
		}

		return err
	}

	return nil
}
