package progress

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type (
	// JSONFileReadWriter reads and writes progress to file
	JSONFileReadWriter struct {
		FilePath string
	}
)

var ErrEmptyProgressFile = fmt.Errorf("empty progress file")

// GetProgress gets progress from file
func (rw *JSONFileReadWriter) GetProgress() ([]WorkInProgress, error) {
	fileBytes, err := os.ReadFile(rw.FilePath)
	if err != nil {
		return nil, fmt.Errorf("read file %q: %w", rw.FilePath, err)
	}

	if len(fileBytes) == 0 {
		return []WorkInProgress{}, ErrEmptyProgressFile
	}

	var wips []WorkInProgress
	if err = json.Unmarshal(fileBytes, &wips); err != nil {
		return nil, fmt.Errorf("unmarshal file %q: %w", rw.FilePath, err)
	}

	return wips, nil
}

// WriteProgress writes progress to file
func (rw *JSONFileReadWriter) WriteProgress(wips []WorkInProgress) error {
	jsonStr, err := json.Marshal(wips)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := ioutil.WriteFile(rw.FilePath, jsonStr, 0644); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
