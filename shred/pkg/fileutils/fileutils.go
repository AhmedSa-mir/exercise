package fileutils

import (
	"os"
	"crypto/rand"
	"io/fs"
)

// CheckFileType checks whether the path file type matches the mode 
func CheckFileType(path string, mode fs.FileMode) (bool, error) {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, err
	}
	if fileInfo.Mode() & mode == 0 {
		return false, nil
	}
	return true, nil
}

// OpenFile implements os.OpenFile and returns filesize
func OpenFile(path string, flag int, mode fs.FileMode) (*os.File, int64, error) {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, 0, os.ErrNotExist
	}

	f, err := os.OpenFile(path, flag, mode)
    if err != nil {
		return nil, 0, err
	}

	return f, fileInfo.Size(), nil
}

// WriteRandomBytes writes bytes to a file at specific index
func WriteRandomBytes(f *os.File, size int64, index int64) error {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		return err
	}

	if _, err := f.WriteAt(buf, index); err != nil {
		return err
	}

	return nil
}