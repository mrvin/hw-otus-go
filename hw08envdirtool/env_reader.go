package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]string

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("ReadDir: %v", err)
	}

	envVar := make(Environment)
	for _, entry := range entries {
		if !entry.IsDir() {
			file, err := os.Open(filepath.Join(dir, entry.Name()))
			if err != nil {
				log.Printf("ReadDir: %v", err)
				continue
			}
			defer file.Close()

			input := bufio.NewScanner(file)
			if input.Scan() {
				valEnv := strings.TrimRight(input.Text(), " \t")
				envVar[entry.Name()] = string(bytes.ReplaceAll([]byte(valEnv), []byte("\x00"), []byte("\n")))
			} else {
				if err := input.Err(); err == nil {
					envVar[entry.Name()] = ""
				}
			}
			if err := input.Err(); err != nil {
				log.Printf("ReadDir: %v", err)
			}
		}
	}

	return envVar, nil
}
