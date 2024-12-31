package main

import (
	"context"

	"fmt"

	"github.com/Rhionin/SanderServer/progress"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

const (
	secretName = "StormlightArchive"
	region     = "us-west-2"
)

func main() {
	lambda.Start(GetProgress)
}

func GetProgress(ctx context.Context) (interface{}, error) {

	config, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("load default config: %w", err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		return nil, fmt.Errorf("get secrets: %w", err)
	}

	// Decrypts secret using the associated KMS key.
	// The value is this string is the json-encoded object found at https://us-west-2.console.aws.amazon.com/secretsmanager/secret?name=StormlightArchive&region=us-west-2
	var secretString string = *result.SecretString
	fmt.Println(secretString[:5], secretString[len(secretString)-5:])

	checker := progress.WebProgressChecker{
		URL: "http://brandonsanderson.com",
	}

	latestProgress, err := checker.GetProgress()
	if err != nil {
		return "", fmt.Errorf("get progress: %w", err)
	}

	response := fmt.Sprintf("Latest progress from %s\n", checker.URL)
	for _, wip := range latestProgress {
		response += fmt.Sprintf("\t%s\n", wip.ToString())
	}

	return response, nil
}
