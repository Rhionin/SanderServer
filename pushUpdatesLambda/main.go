package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	appconfig "github.com/Rhionin/SanderServer/config"
	"github.com/Rhionin/SanderServer/history"
	"github.com/Rhionin/SanderServer/progress"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

func PushUpdates(ctx context.Context, event DynamoDBEvent) error {
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

	if len(event.Records) != 1 {
		return fmt.Errorf("expected 1 record in event trigger, but got %d", len(event.Records))
	}

	latestUpdate := event.Records[0]
	if latestUpdate.EventName != "INSERT" {
		fmt.Printf("Not processing %q operation\n", latestUpdate.EventName)
		return nil
	}
	var latestHistoryEntry history.ProgressDynamoEntry
	if err = dynamodbattribute.UnmarshalMap(latestUpdate.Change.NewImage, &latestHistoryEntry); err != nil {
		return fmt.Errorf("dynamodbattribute.UnmarshalMap for event record: %w", err)
	}

	penultimateUpdate, err := historyClient.GetLatestProgressEntryBeforeID(ctx, latestHistoryEntry)
	if err != nil && !errors.Is(err, history.ErrEmptyHistory) {
		return fmt.Errorf("get penultimate progress update entry: %w", err)
	} else if errors.Is(err, history.ErrEmptyHistory) {
		return history.ErrEmptyHistory
	} else if errors.Is(err, history.ErrNoEntryBeforeTarget) {
		fmt.Println("This appears to be the first history entry. No updates to push.")
		return nil
	}
	updates := progress.GetProgressUpdate(latestHistoryEntry.WorksInProgress, penultimateUpdate.WorksInProgress)

	channelOverride := "" // Post to the default channel
	if err := progress.SendSlackUpdate(secrets.SlackWebhookURL, updates, channelOverride); err != nil {
		return fmt.Errorf("send slack update: %w", err)
	}

	return nil
}

// Unmarshal the event in a way where we can actually parse the event records using dynamodbattribute.UnmarshalMap.
// See https://stackoverflow.com/a/50164289
type DynamoDBEvent struct {
	Records []DynamoDBEventRecord `json:"Records"`
}

type DynamoDBEventRecord struct {
	AWSRegion      string                       `json:"awsRegion"`
	Change         DynamoDBStreamRecord         `json:"dynamodb"`
	EventID        string                       `json:"eventID"`
	EventName      string                       `json:"eventName"`
	EventSource    string                       `json:"eventSource"`
	EventVersion   string                       `json:"eventVersion"`
	EventSourceArn string                       `json:"eventSourceARN"`
	UserIdentity   *events.DynamoDBUserIdentity `json:"userIdentity,omitempty"`
}

type DynamoDBStreamRecord struct {
	ApproximateCreationDateTime events.SecondsEpochTime `json:"ApproximateCreationDateTime,omitempty"`
	// changed to map[string]*dynamodb.AttributeValue
	Keys map[string]*dynamodb.AttributeValue `json:"Keys,omitempty"`
	// changed to map[string]*dynamodb.AttributeValue
	NewImage map[string]*dynamodb.AttributeValue `json:"NewImage,omitempty"`
	// changed to map[string]*dynamodb.AttributeValue
	OldImage       map[string]*dynamodb.AttributeValue `json:"OldImage,omitempty"`
	SequenceNumber string                              `json:"SequenceNumber"`
	SizeBytes      int64                               `json:"SizeBytes"`
	StreamViewType string                              `json:"StreamViewType"`
}
