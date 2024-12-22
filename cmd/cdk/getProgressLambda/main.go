package main

import (
	"context"
	"os"

	"fmt"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/Rhionin/SanderServer/progress"
)

var (
	username = os.Getenv("GIT_USERNAME")
	apiKey   = os.Getenv("GIT_API_KEY")
)

func main() {
	lambda.Start(GetProgress)
}

func GetProgress(ctx context.Context) error {
	checker := progress.WebProgressChecker{
		URL: "http://brandonsanderson.com",
	}

	latestProgress, err := checker.GetProgress()
	if err != nil {
		return fmt.Errorf("get progress: %w", err)
	}

	fmt.Println("Latest progress from " + checker.URL)
	for _, wip := range latestProgress {
		fmt.Println("\t", wip.ToString())
	}

	return nil
}
