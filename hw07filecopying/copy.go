package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	infFromFile, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if offset > infFromFile.Size() {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 || limit+offset > infFromFile.Size() {
		limit = infFromFile.Size() - offset
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	toFile, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, infFromFile.Mode())
	if err != nil {
		return err
	}

	_, err = io.CopyN(toFile, fromFile, limit)

	if closeErr := toFile.Close(); err == nil {
		err = closeErr
	}

	return err
}
