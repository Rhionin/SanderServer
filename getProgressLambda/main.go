package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/Rhionin/SanderServer/history"
	"github.com/Rhionin/SanderServer/progress"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
// secretName = "StormlightArchive"
)

type httpResponse struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

func main() {
	lambda.Start(GetProgress)
}

func GetProgress(ctx context.Context) (interface{}, error) {

	// // Create Secrets Manager client
	// svc := secretsmanager.NewFromConfig(config)

	// input := &secretsmanager.GetSecretValueInput{
	// 	SecretId:     aws.String(secretName),
	// 	VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	// }

	// result, err := svc.GetSecretValue(ctx, input)
	// if err != nil {
	// 	// For a list of exceptions thrown, see
	// 	// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
	// 	return nil, fmt.Errorf("get secrets: %w", err)
	// }

	// // Decrypts secret using the associated KMS key.
	// // The value is this string is the json-encoded object found at https://us-west-2.console.aws.amazon.com/secretsmanager/secret?name=StormlightArchive&region=us-west-2
	// var secretString string = *result.SecretString

	// var secrets secretStore
	// if err = json.Unmarshal([]byte(secretString), &secrets); err != nil {
	// 	return "", fmt.Errorf("unmarshal secrets: %w", err)
	// }

	checker := progress.WebProgressChecker{
		URL: "http://brandonsanderson.com",
	}

	latestProgress, err := checker.GetProgress()
	if err != nil {
		return "", fmt.Errorf("get progress: %w", err)
	}

	latestProgressSimplified := []progress.WorkInProgressSimple{}
	for _, p := range latestProgress {
		latestProgressSimplified = append(latestProgressSimplified, progress.WorkInProgressSimple{
			Title:    p.Title,
			Progress: p.Progress,
		})
	}
	historyClient, err := history.NewDynamoClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("new dynamo client: %w", err)
	}
	latestProgressFromHistory, err := historyClient.GetLatestProgressEntry(ctx)
	if err != nil && !errors.Is(err, history.ErrEmptyHistory) {
		return nil, fmt.Errorf("get latest progress entry from history: %w", err)
	}

	shouldAddHistoryEntry := errors.Is(err, history.ErrEmptyHistory) || !reflect.DeepEqual(latestProgressFromHistory.WorksInProgress, latestProgressSimplified)
	if shouldAddHistoryEntry {
		progressEntry := history.ProgressEntry{
			Timestamp:       time.Now(),
			WorksInProgress: latestProgressSimplified,
		}
		if errors.Is(err, history.ErrEmptyHistory) {
			fmt.Println("History does not have any entries yet. Adding new entry with timestamp", progressEntry.Timestamp)
		} else {
			fmt.Println("Current progress is different from previous history entry. Adding new entry with timestamp", progressEntry.Timestamp)
		}
		if err = historyClient.AddNewProgressEntry(ctx, progressEntry); err != nil {
			return nil, fmt.Errorf("add new history entry: %w", err)
		}
	} else {
		fmt.Println("No progress change.")
	}

	page, err := progress.CreateStatusPage(latestProgress)
	if err != nil {
		return "", fmt.Errorf("create status page: %w", err)
	}

	return httpResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		Body: string(page),
	}, nil
}
