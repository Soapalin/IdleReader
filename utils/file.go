package utils

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

func IsFileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func IsDir(path string) bool {
	return true
}

// {KeepFromEnd} is a helper function to keep number of byte of the file
// {end} is the number of byte to keep
func KeepFromEnd(path string, end int) {
	fin, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	fi, err := fin.Stat()
	if err != nil {
		panic(err)
	}

	if fi.Size() <= int64(end) {
		return
	}

	dir, _ := filepath.Split(path)
	fout, err := os.Create(filepath.Join(dir, "tmp.log"))
	if err != nil {
		panic(err)
	}

	_, err = fin.Seek(fi.Size()-int64(end), io.SeekStart)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(fout, fin)
	if err != nil {
		panic(err)
	}

	fin.Close()
	fout.Close()

	if err := os.Remove(path); err != nil {
		panic(err)
	}

	if err := os.Rename(filepath.Join(dir, "tmp.log"), path); err != nil {
		panic(err)
	}

}
