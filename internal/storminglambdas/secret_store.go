package storminglambdas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type (
	StormlightArchiveClient struct {
		SecretsManager awsSecretsManager
	}

	StormlightArchive struct {
		SlackWebhookURL string `json:"SLACK_WEBHOOK_URL"`
	}

	awsSecretsManager interface {
		GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
	}
)

func NewStormlightArchiveClient(secretsManager awsSecretsManager) *StormlightArchiveClient {

	return &StormlightArchiveClient{
		SecretsManager: secretsManager,
	}
}

func (client *StormlightArchiveClient) GetSecrets(ctx context.Context) (StormlightArchive, error) {
	result, err := client.SecretsManager.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	})
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		return StormlightArchive{}, fmt.Errorf("get secrets: %w", err)
	}

	// Decrypts secret using the associated KMS key.
	// The value is this string is the json-encoded object found at https://us-west-2.console.aws.amazon.com/secretsmanager/secret?name=StormlightArchive&region=us-west-2
	var secretString string = *result.SecretString

	var secrets StormlightArchive
	if err = json.Unmarshal([]byte(secretString), &secrets); err != nil {
		return StormlightArchive{}, fmt.Errorf("unmarshal secrets: %w", err)
	}

	return secrets, nil
}
