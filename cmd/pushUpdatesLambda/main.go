package main

import (
	"github.com/Rhionin/SanderServer/internal/storminglambdas"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(storminglambdas.PushUpdates)
}
