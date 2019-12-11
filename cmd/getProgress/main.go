package main

import (
	"fmt"

	"github.com/Rhionin/SanderServer/progress"
)

func main() {
	checker := progress.WebProgressChecker{
		URL: "http://brandonsanderson.com",
	}

	latestProgress := checker.GetProgress()

	fmt.Println("Latest progress from " + checker.URL)
	for _, wip := range latestProgress {
		fmt.Println("\t", wip.ToString())
	}
}
