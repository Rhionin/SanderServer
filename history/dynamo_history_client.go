package history

import (
	"context"
	"errors"
	"fmt"
	"time"

	appconfig "github.com/Rhionin/SanderServer/config"

	"github.com/Rhionin/SanderServer/progress"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	latestEntryID = "latest_entry"
)

var (
	ErrEmptyHistory = errors.New("history database is empty")
)

type (
	DynamoClient struct {
		client *dynamodb.Client
	}

	ProgressEntry struct {
		Timestamp       time.Time
		WorksInProgress []progress.WorkInProgressSimple
	}

	progressDynamoEntry struct {
		ID                string
		TimestampUnixNano int64
		WorksInProgress   []progress.WorkInProgressSimple
	}
)

func NewDynamoClient(ctx context.Context) (*DynamoClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(appconfig.AWSRegion))
	if err != nil {
		return nil, fmt.Errorf("load default config: %w", err)
	}

	c := dynamodb.NewFromConfig(cfg)

	return &DynamoClient{client: c}, nil
}

func (c *DynamoClient) GetLatestProgressEntry(ctx context.Context) (ProgressEntry, error) {
	latestProgressResult, err := c.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(appconfig.HistoryDynamoTableName),
		KeyConditionExpression: aws.String("ID = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: latestEntryID},
		},
		ScanIndexForward: aws.Bool(false), // Get the latest (descending order)
		Limit:            aws.Int32(1),
	})
	if err != nil {
		return ProgressEntry{}, fmt.Errorf("get latest progress from DynamoDB: %w", err)
	}

	if len(latestProgressResult.Items) == 0 {
		count, err := c.GetEntryCount(ctx)
		if err != nil {
			return ProgressEntry{}, fmt.Errorf("get entry count: %w", err)
		}
		if count > 0 {
			return ProgressEntry{}, fmt.Errorf("could not find latest progress entry; zero results returned from dynamoDB despite %d records existing", count)
		}
		return ProgressEntry{}, ErrEmptyHistory
	}

	var latestProgressFromHistory progressDynamoEntry
	err = attributevalue.UnmarshalMap(latestProgressResult.Items[0], &latestProgressFromHistory)
	if err != nil {
		return ProgressEntry{}, fmt.Errorf("unmarshal DynamoDB item: %w", err)
	}

	return latestProgressFromHistory.toProgressEntry(), nil
}

func (c *DynamoClient) AddNewProgressEntry(ctx context.Context, entry ProgressEntry) error {
	dynamoItem, err := attributevalue.MarshalMap(entry.toDynamoProgressEntry())
	if err != nil {
		return fmt.Errorf("marshal progress dymamo entry: %w", err)
	}
	if _, err = c.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(appconfig.HistoryDynamoTableName),
		Item:      dynamoItem,
	}); err != nil {
		return fmt.Errorf("put history entry into dynamoDB: %w", err)
	}

	return nil
}

func (c *DynamoClient) GetEntryCount(ctx context.Context) (int32, error) {
	result, err := c.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(appconfig.HistoryDynamoTableName),
		KeyConditionExpression: aws.String("ID = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: latestEntryID},
		},
		Select: types.SelectCount,
	})
	if err != nil {
		return 0, fmt.Errorf("dynamo: get entry count: %w", err)
	}

	return result.Count, nil
}

func (e progressDynamoEntry) toProgressEntry() ProgressEntry {
	return ProgressEntry{
		Timestamp:       time.Unix(0, e.TimestampUnixNano),
		WorksInProgress: e.WorksInProgress,
	}
}

func (e ProgressEntry) toDynamoProgressEntry() progressDynamoEntry {
	return progressDynamoEntry{
		ID:                latestEntryID,
		TimestampUnixNano: e.Timestamp.UnixNano(),
		WorksInProgress:   e.WorksInProgress,
	}
}
