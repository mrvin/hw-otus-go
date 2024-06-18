package main

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	fromPath := "testdata/input.txt"
	toPath := "out.txt"
	isQuiet := true
	var tests = []struct {
		offset int64
		limit  int64
		want   string
	}{
		{0, 0, "testdata/out_offset0_limit0.txt"},
		{0, 10, "testdata/out_offset0_limit10.txt"},
		{0, 1000, "testdata/out_offset0_limit1000.txt"},
		{0, 10000, "testdata/out_offset0_limit10000.txt"},
		{100, 1000, "testdata/out_offset100_limit1000.txt"},
		{6000, 1000, "testdata/out_offset6000_limit1000.txt"},
	}

	for _, test := range tests {
		if err := Copy(fromPath, toPath, test.offset, test.limit, isQuiet); err != nil {
			t.Errorf("err = %v, want %v", err, nil)
		}
		defer os.Remove(toPath)

		if ok, _ := cmpFiles(toPath, test.want); !ok {
			t.Errorf("files: %s, %s - not equel", toPath, test.want)
		}
	}
}

func TestCopyExceededOffset(t *testing.T) {
	fromPath := "testdata/input.txt"
	toPath := "out.txt"
	isQuiet := true

	infFile, _ := os.Stat(fromPath)

	if err := Copy(fromPath, toPath, infFile.Size()+100, 0, isQuiet); !errors.Is(err, ErrOffsetExceedsFileSize) {
		t.Errorf("err = %v, want %v", err, nil)
	}
}

func cmpFiles(filePath1, filePath2 string) (bool, error) {
	infFile1, err := os.Stat(filePath1)
	if err != nil {
		return false, err
	}

	infFile2, err := os.Stat(filePath2)
	if err != nil {
		return false, err
	}

	if infFile1.Size() != infFile2.Size() {
		return false, nil
	}

	file1, err := os.ReadFile(filePath1)
	if err != nil {
		return false, err
	}

	file2, err := os.ReadFile(filePath2)
	if err != nil {
		return false, err
	}

	return bytes.Equal(file1, file2), nil
}
