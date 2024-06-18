package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64, isQuiet bool) error {
	infFromFile, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("file to read from: %w", err)
	}
	if !infFromFile.Mode().IsRegular() {
		return fmt.Errorf("file to read from: %w", ErrUnsupportedFile)
	}

	if offset > infFromFile.Size() {
		return fmt.Errorf("file to read from: %w", ErrOffsetExceedsFileSize)
	}
	if limit == 0 || limit+offset > infFromFile.Size() {
		limit = infFromFile.Size() - offset
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("file to read from: %w", err)
	}
	defer fromFile.Close()

	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("file to read from: %w", err)
	}

	toFile, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, infFromFile.Mode())
	if err != nil {
		return fmt.Errorf("file to write to: %w", err)
	}

	var reader io.Reader = fromFile
	if !isQuiet {
		bar := pb.Full.Start64(limit)
		reader = bar.NewProxyReader(fromFile)
		defer bar.Finish()
	}

	_, err = io.CopyN(toFile, reader, limit)

	if closeErr := toFile.Close(); err == nil {
		err = closeErr
	}

	return err
}
