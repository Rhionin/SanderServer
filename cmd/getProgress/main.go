package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Rhionin/SanderServer/internal/progress"
)

func main() {
	checker := progress.WebProgressChecker{
		URL: "http://brandonsanderson.com",
	}

	latestProgress, err := checker.GetProgress()
	if err != nil {
		log.Fatalf("get progress: %s", err)
	}

	fmt.Println("Latest progress from " + checker.URL)
	for _, wip := range latestProgress {
		fmt.Println("\t", wip.String())
	}

	var page []byte
	if len(latestProgress) == 0 {
		fmt.Println("\tNo works in progress detected...")
		page = progress.ErrorPageContent
	} else {
		page, err = progress.CreateStatusPage(latestProgress)
		if err != nil {
			log.Fatal(err)
		}
	}

	// create and open the status-page.html file to write into it
	filename := "status-page.html"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}

	// write into the file
	_, err = file.Write(page)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}

	// close the file
	file.Close()

	// get the absolute path of the written file
	absPath, err := filepath.Abs(filename)
	if err != nil {
		fmt.Println("Error getting absolute file path:", err)
		os.Exit(1)
	}

	// print the absolute path of the written file
	fmt.Println("Status page created successfully at:", absPath)
}
