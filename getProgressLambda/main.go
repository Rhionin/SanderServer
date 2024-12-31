package main

import (
	"context"
	"os"

	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	username = os.Getenv("GIT_USERNAME")
	apiKey   = os.Getenv("GIT_API_KEY")
)

func main() {
	lambda.Start(GetProgress)
}

func GetProgress(ctx context.Context) (interface{}, error) {
	// checker := progress.WebProgressChecker{
	// 	URL: "http://brandonsanderson.com",
	// }

	// latestProgress, err := checker.GetProgress()
	// if err != nil {
	// 	return fmt.Errorf("get progress: %w", err)
	// }

	// fmt.Println("Latest progress from " + checker.URL)
	// for _, wip := range latestProgress {
	// 	fmt.Println("\t", wip.ToString())
	// }
	fmt.Println("Hello, StormWatch logs!!")

	return "Hello, StormWatch caller!!", nil
}
