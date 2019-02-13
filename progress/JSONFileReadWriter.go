package progress

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type (
	// JSONFileReadWriter reads and writes progress to file
	JSONFileReadWriter struct {
		FilePath string
	}
)

// GetProgress gets progress from file
func (rw *JSONFileReadWriter) GetProgress() []WorkInProgress {
	fileBytes, err := ioutil.ReadFile(rw.FilePath)
	if err != nil {
		fmt.Println(err)
		return []WorkInProgress{}
	}

	var wips []WorkInProgress
	if err = json.Unmarshal(fileBytes, &wips); err != nil {
		fmt.Println(err)
		return []WorkInProgress{}
	}

	return wips
}

// WriteProgress writes progress to file
func (rw *JSONFileReadWriter) WriteProgress(wips []WorkInProgress) error {
	jsonStr, err := json.Marshal(wips)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := ioutil.WriteFile(rw.FilePath, jsonStr, 0644); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
