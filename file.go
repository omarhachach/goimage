package goimage

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/omar-h/goimage/utils"
)

// File holds a file.
type File struct {
	File      multipart.File
	Header    *multipart.FileHeader
	Basename  string // Without extension
	Fullname  string // With extension
	Extension string
	MIMEType  string
	Size      int
}

// NewFile will create a new file from a multipart.FileHeader.
func NewFile(file multipart.File, fileHeader *multipart.FileHeader) *File {
	return &File{
		File:      file,
		Header:    fileHeader,
		Basename:  utils.GetFileBasename(fileHeader.Filename),
		Fullname:  fileHeader.Filename,
		Extension: filepath.Ext(fileHeader.Filename),
		MIMEType:  fileHeader.Header["Content-Type"][0],
		Size:      int(fileHeader.Size),
	}
}

// Place will move the file onto a specific location.
// Returns os package errors. (os.ErrFileExist and os.ErrPermission)
func (f *File) Place(location string) error {
	file, err := os.OpenFile(location+f.Fullname, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, f.File)
	if err != nil {
		return err
	}

	return nil
}

// GenerateName will generate a new name with a given length.
func (f *File) GenerateName(len int) *File {
	return f.GiveName(utils.GenerateName(len))
}

// GiveName will give the File a new name, and update the basename and fullname.
func (f *File) GiveName(name string) *File {
	f.Basename = name
	f.Fullname = name + f.Extension
	return f
}

// Close will properly close the file.
func (f *File) Close() error {
	return f.File.Close()
}
