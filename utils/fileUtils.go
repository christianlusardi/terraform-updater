package utils

import (
	"github.com/spf13/afero"
	"io"
	"os"
)

type FileService struct {
	Os afero.Fs
}

type IFileService interface {
	MoveFile(source string, destination string) error
	FileExists(file string) bool
}

func (fs *FileService) MoveFile(source string, destination string) error {
	src, err := fs.Os.Open(source)

	if err != nil {
		return err
	}

	defer src.Close()

	fi, err := src.Stat()

	if err != nil {
		return err
	}

	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	perm := fi.Mode() & os.ModePerm
	dst, err := fs.Os.OpenFile(destination, flag, perm)

	if err != nil {
		return err
	}

	defer dst.Close()
	_, err = io.Copy(dst, src)

	if err != nil {
		dst.Close()
		fs.Os.Remove(destination)
		return err
	}
	err = dst.Close()

	if err != nil {
		return err
	}
	err = src.Close()

	if err != nil {
		return err
	}
	err = fs.Os.Remove(source)

	if err != nil {
		return err
	}

	return nil
}



func (fs *FileService) FileExists(file string) bool {

	if _, err := fs.Os.Stat(file); err == nil {
		return true
	}

	return false

}
