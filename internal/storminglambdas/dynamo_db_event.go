package storminglambdas

import (
	"fmt"

	"github.com/Rhionin/SanderServer/internal/history"
	"github.com/Rhionin/SanderServer/internal/progress"
	"github.com/aws/aws-lambda-go/events"
)

// UnmarshalStreamImage reads a DynamoDB stream NewImage into a ProgressDynamoEntry.
func UnmarshalStreamImage(image map[string]events.DynamoDBAttributeValue, out *history.ProgressDynamoEntry) error {
	timestamp, err := image["TimestampUnixNano"].Integer()
	if err != nil {
		return fmt.Errorf("parse TimestampUnixNano: %w", err)
	}

	works := image["WorksInProgress"].List()
	worksInProgress := make([]progress.WorkInProgress, len(works))
	for i, work := range works {
		fields := work.Map()
		p, err := fields["Progress"].Integer()
		if err != nil {
			return fmt.Errorf("parse WorksInProgress[%d].Progress: %w", i, err)
		}
		worksInProgress[i] = progress.WorkInProgress{
			Title:    fields["Title"].String(),
			Progress: int(p),
		}
	}

	out.ID = image["ID"].String()
	out.TimestampUnixNano = timestamp
	out.WorksInProgress = worksInProgress
	return nil
}
