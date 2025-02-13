package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/Rhionin/SanderServer/internal/history"
	"github.com/Rhionin/SanderServer/internal/progress"
	"github.com/aws/aws-lambda-go/lambda"
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
	checker := progress.WebProgressChecker{
		URL: "http://brandonsanderson.com",
	}

	latestProgress, err := checker.GetProgress()
	if err != nil {
		return "", fmt.Errorf("get progress: %w", err)
	}

	historyClient, err := history.NewDynamoClientFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("new dynamo client: %w", err)
	}
	latestProgressFromHistory, err := historyClient.GetLatestProgressEntry(ctx)
	if err != nil && !errors.Is(err, history.ErrEmptyHistory) {
		return nil, fmt.Errorf("get latest progress entry from history: %w", err)
	}

	shouldAddHistoryEntry := errors.Is(err, history.ErrEmptyHistory) || !reflect.DeepEqual(latestProgressFromHistory.WorksInProgress, latestProgress)
	if shouldAddHistoryEntry {
		progressEntry := history.ProgressEntry{
			Timestamp:       time.Now(),
			WorksInProgress: latestProgress,
		}
		if errors.Is(err, history.ErrEmptyHistory) {
			fmt.Println("History does not have any entries yet. Adding new entry with timestamp", progressEntry.Timestamp)
		} else {
			fmt.Println("Current progress is different from previous history entry. Adding new entry with timestamp", progressEntry.Timestamp)
			latestBytes, _ := json.Marshal(latestProgress)
			fmt.Println("Current progress:", string(latestBytes))

			previousBytes, _ := json.Marshal(latestProgressFromHistory)
			fmt.Println("Previous progress:", string(previousBytes))
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
