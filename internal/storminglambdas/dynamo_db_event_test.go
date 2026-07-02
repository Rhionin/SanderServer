package storminglambdas

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Rhionin/SanderServer/internal/history"
	"github.com/Rhionin/SanderServer/internal/progress"
	"github.com/aws/aws-lambda-go/events"
)

func TestUnmarshalStreamImage(t *testing.T) {
	testCases := []struct {
		name     string
		record   string
		expected history.ProgressDynamoEntry
		expectErr bool
	}{
		{
			name: "returns ProgressDynamoEntry for valid stream image",
			record: `{
				"dynamodb": {
					"NewImage": {
						"ID": {"S": "latest_entry"},
						"TimestampUnixNano": {"N": "12345"},
						"WorksInProgress": {
							"L": [
								{"M": {"Progress": {"N": "100"}, "Title": {"S": "Moment Zero 2.0"}}},
								{"M": {"Progress": {"N": "100"}, "Title": {"S": "White Sand Prewriting (Prose Version)"}}},
								{"M": {"Progress": {"N": "28"}, "Title": {"S": "Words of Radiance Order Progress: Truthwatcher"}}},
								{"M": {"Progress": {"N": "23"}, "Title": {"S": "Words of Radiance BackerKit Progress: $325 Tier"}}}
							]
						}
					}
				}
			}`,
			expected: history.ProgressDynamoEntry{
				ID:                "latest_entry",
				TimestampUnixNano: 12345,
				WorksInProgress: []progress.WorkInProgress{
					{Title: "Moment Zero 2.0", Progress: 100},
					{Title: "White Sand Prewriting (Prose Version)", Progress: 100},
					{Title: "Words of Radiance Order Progress: Truthwatcher", Progress: 28},
					{Title: "Words of Radiance BackerKit Progress: $325 Tier", Progress: 23},
				},
			},
		},
		{
			name: "returns error when Progress is not numeric",
			record: `{
				"dynamodb": {
					"NewImage": {
						"ID": {"S": "latest_entry"},
						"TimestampUnixNano": {"N": "12345"},
						"WorksInProgress": {
							"L": [{
								"M": {
									"Progress": {"N": "not-a-number"},
									"Title": {"S": "Some Book"}
								}
							}]
						}
					}
				}
			}`,
			expectErr: true,
		},
		{
			name: "returns error when TimestampUnixNano is not numeric",
			record: `{
				"dynamodb": {
					"NewImage": {
						"ID": {"S": "latest_entry"},
						"TimestampUnixNano": {"N": "not-a-number"},
						"WorksInProgress": {"L": []}
					}
				}
			}`,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var record events.DynamoDBEventRecord
			if err := json.Unmarshal([]byte(tc.record), &record); err != nil {
				t.Fatalf("unmarshal test record: %s", err)
			}

			var actual history.ProgressDynamoEntry
			err := UnmarshalStreamImage(record.Change.NewImage, &actual)

			if tc.expectErr {
				if err == nil {
					t.Fatal("expected error, actual nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %+v, actual %+v", tc.expected, actual)
			}
		})
	}
}
