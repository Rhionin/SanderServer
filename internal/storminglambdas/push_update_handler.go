package storminglambdas

import (
	"context"
	"errors"
	"fmt"

	appconfig "github.com/Rhionin/SanderServer/config"
	"github.com/Rhionin/SanderServer/internal/history"
	"github.com/Rhionin/SanderServer/internal/progress"
	"github.com/Rhionin/SanderServer/internal/slack"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

const (
	secretName = "StormlightArchive"
)

type (
	PushUpdateHandler struct {
		History     historyClient
		PushTargets []PushTarget
	}

	historyClient interface {
		GetLatestProgressEntryBeforeID(ctx context.Context, targetEntry history.ProgressDynamoEntry) (history.ProgressEntry, error)
	}

	// PushTarget the interface for sending push notifications
	PushTarget interface {
		GetName() string
		SendUpdate(ctx context.Context, updates []progress.ProgressUpdate) error
	}
)

// PushUpdates sends notifications when a progress update occurs
func PushUpdates(ctx context.Context, event events.DynamoDBEvent) error {
	handler, err := NewPushUpdateHandlerFromContext(ctx)
	if err != nil {
		return fmt.Errorf("new push update handler from context: %w", err)
	}

	return handler.PushUpdates(ctx, event)
}

// NewPushUpdateHandlerFromContext creates a new push update handler by initializing dependencies from ctx
func NewPushUpdateHandlerFromContext(ctx context.Context) (*PushUpdateHandler, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(appconfig.AWSRegion))
	if err != nil {
		return nil, fmt.Errorf("load default config: %w", err)
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)
	historyClient, err := history.NewDynamoClient(dynamoClient)
	if err != nil {
		return nil, fmt.Errorf("new history dynamo client: %w", err)
	}

	secretsManager := secretsmanager.NewFromConfig(cfg)
	stormlightArchiveClient := NewStormlightArchiveClient(secretsManager)
	config, err := stormlightArchiveClient.GetSecrets(ctx)
	if err != nil {
		return nil, fmt.Errorf("get stormlight archive: %w", err)
	}

	slackChannelOverride := "" // Post to the default channel
	pushTargets := []PushTarget{
		slack.NewUpdateClient(config.SlackWebhookURL, slackChannelOverride),
	}

	return &PushUpdateHandler{
		History:     historyClient,
		PushTargets: pushTargets,
	}, nil
}

// PushUpdates sends notifications when a progress update occurs
func (handler *PushUpdateHandler) PushUpdates(ctx context.Context, event events.DynamoDBEvent) error {
	if len(event.Records) != 1 {
		return fmt.Errorf("expected 1 record in event trigger, but got %d", len(event.Records))
	}

	latestUpdate := event.Records[0]
	if latestUpdate.EventName != "INSERT" {
		fmt.Printf("Not processing %q operation\n", latestUpdate.EventName)
		return nil
	}

	var latestHistoryEntry history.ProgressDynamoEntry
	if err := UnmarshalStreamImage(latestUpdate.Change.NewImage, &latestHistoryEntry); err != nil {
		return fmt.Errorf("unmarshal stream image: %w", err)
	}

	penultimateUpdate, err := handler.History.GetLatestProgressEntryBeforeID(ctx, latestHistoryEntry)
	if err != nil && !errors.Is(err, history.ErrEmptyHistory) {
		return fmt.Errorf("get penultimate progress update entry: %w", err)
	} else if errors.Is(err, history.ErrEmptyHistory) {
		return history.ErrEmptyHistory
	} else if errors.Is(err, history.ErrNoEntryBeforeTarget) {
		fmt.Println("This appears to be the first history entry. No updates to push.")
		return nil
	}
	updates := progress.GetProgressUpdate(latestHistoryEntry.WorksInProgress, penultimateUpdate.WorksInProgress)

	for _, target := range handler.PushTargets {
		targetName := target.GetName()
		if err := target.SendUpdate(ctx, updates); err != nil {
			return fmt.Errorf("(%s) send update: %w", targetName, err)
		}
		fmt.Println("Update sent via", targetName)
	}

	return nil
}
