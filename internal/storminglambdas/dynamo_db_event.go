package storminglambdas

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// UnmarshalStreamImage converts events.DynamoDBAttributeValue to struct
// See https://stackoverflow.com/a/50017398
func UnmarshalStreamImage[V any](attribute map[string]events.DynamoDBAttributeValue, out *V) error {
	dbAttrMap := make(map[string]*dynamodb.AttributeValue)
	for k, v := range attribute {

		var dbAttr dynamodb.AttributeValue

		bytes, marshalErr := v.MarshalJSON()
		if marshalErr != nil {
			return fmt.Errorf("marshal event key %q: %s", k, marshalErr)
		}

		if err := json.Unmarshal(bytes, &dbAttr); err != nil {
			return fmt.Errorf("unmarshal event key %q: %s", k, err)
		}
		dbAttrMap[k] = &dbAttr
	}

	return dynamodbattribute.UnmarshalMap(dbAttrMap, out)

}
