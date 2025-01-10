package storminglambdas

import (
	"encoding/json"
	"testing"

	"github.com/Rhionin/SanderServer/internal/history"
	"github.com/Rhionin/SanderServer/internal/progress"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalStreamImage(t *testing.T) {
	recordPayload := []byte(`{
    "dynamodb": {
        "NewImage": {
            "ID": {
                "S": "latest_entry"
            },
            "TimestampUnixNano": {
                "N": "12345"
            },
            "WorksInProgress": {
                "L": [
                    {
                        "M": {
                            "Progress": {
                                "N": "100"
                            },
                            "Title": {
                                "S": "Moment Zero 2.0"
                            }
                        }
                    },
                    {
                        "M": {
                            "Progress": {
                                "N": "100"
                            },
                            "Title": {
                                "S": "White Sand Prewriting (Prose Version)"
                            }
                        }
                    },
                    {
                        "M": {
                            "Progress": {
                                "N": "28"
                            },
                            "Title": {
                                "S": "Words of Radiance Order Progress: Truthwatcher"
                            }
                        }
                    },
                    {
                        "M": {
                            "Progress": {
                                "N": "23"
                            },
                            "Title": {
                                "S": "Words of Radiance BackerKit Progress: $325 Tier"
                            }
                        }
                    }
                ]
            }
        }
    }
}`)

	var record events.DynamoDBEventRecord
	if err := json.Unmarshal(recordPayload, &record); err != nil {
		t.Fatal(err)
	}

	var actualEntry history.ProgressDynamoEntry
	if err := UnmarshalStreamImage(record.Change.NewImage, &actualEntry); err != nil {
		t.Fatalf("unmarshal stream image: %s", err)
	}

	expectedEntry := history.ProgressDynamoEntry{
		ID:                "latest_entry",
		TimestampUnixNano: 12345,
		WorksInProgress: []progress.WorkInProgress{
			{Title: "Moment Zero 2.0", Progress: 100},
			{Title: "White Sand Prewriting (Prose Version)", Progress: 100},
			{Title: "Words of Radiance Order Progress: Truthwatcher", Progress: 28},
			{Title: "Words of Radiance BackerKit Progress: $325 Tier", Progress: 23},
		},
	}
	require.Equal(t, expectedEntry, actualEntry)
}
