package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	appconfig "github.com/Rhionin/SanderServer/config"
	"github.com/Rhionin/SanderServer/history"

	"github.com/Rhionin/SanderServer/progress"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

const (
	secretName = "StormlightArchive"
)

type secretStore struct {
	SlackWebhookURL string `json:"SLACK_WEBHOOK_URL"`
}

func main() {
	lambda.Start(PushUpdates)
}

func PushUpdates(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(appconfig.AWSRegion))
	if err != nil {
		return fmt.Errorf("load default config: %w", err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		return fmt.Errorf("get secrets: %w", err)
	}

	// Decrypts secret using the associated KMS key.
	// The value is this string is the json-encoded object found at https://us-west-2.console.aws.amazon.com/secretsmanager/secret?name=StormlightArchive&region=us-west-2
	var secretString string = *result.SecretString

	var secrets secretStore
	if err = json.Unmarshal([]byte(secretString), &secrets); err != nil {
		return fmt.Errorf("unmarshal secrets: %w", err)
	}

	historyClient, err := history.NewDynamoClient(ctx)
	if err != nil {
		return fmt.Errorf("new dynamo client: %w", err)
	}
	progressEntries, err := historyClient.GetLatestProgressEntries(ctx, 2)
	if err != nil && !errors.Is(err, history.ErrEmptyHistory) {
		return fmt.Errorf("get latest 2 progress entries: %w", err)
	}

	var updates []progress.ProgressUpdate
	switch len(progressEntries) {
	case 0:
		return history.ErrEmptyHistory
	case 1:
		for _, entry := range progressEntries[0].WorksInProgress {
			updates = append(updates, progress.ProgressUpdate{
				Title:    entry.Title,
				Progress: entry.Progress,
			})
		}
	case 2:
		updates = progress.GetProgressUpdate(progressEntries[0].WorksInProgress, progressEntries[1].WorksInProgress)
	default:
		return fmt.Errorf("expected no more than 2 progress entries, got %d", len(progressEntries))
	}

	channelOverride := "#cjc-slack-testing"
	if err := progress.SendSlackUpdate(secrets.SlackWebhookURL, updates, channelOverride); err != nil {
		return fmt.Errorf("send slack update: %w", err)
	}

	return nil
}
