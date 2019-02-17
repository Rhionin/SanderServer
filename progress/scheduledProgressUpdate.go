package progress

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"firebase.google.com/go/messaging"
)

type (
	// Monitor monitors progress
	Monitor struct {
		LiveReader Reader
		History    ReadWriter
	}

	// Reader reads progress
	Reader interface {
		GetProgress() []WorkInProgress
	}

	// Writer writes progress
	Writer interface {
		WriteProgress(wips []WorkInProgress) error
	}

	// ReadWriter read and writes progress
	ReadWriter interface {
		Reader
		Writer
	}
)

// ScheduleProgressCheckJob schedules a job to repeatedly check progress
// on Brandon Sanderson's books
func (m *Monitor) ScheduleProgressCheckJob(ctx context.Context, firebaseClient *messaging.Client) {
	prevWips := m.History.GetProgress()

	config := cfg.GetConfig()

	c := cron.New()
	fmt.Println(config.ProgressCheckInterval)
	c.AddFunc(config.ProgressCheckInterval, func() {
		currentWips := CheckProgress() // m.LiveReader.GetProgress()
		if len(currentWips) > 0 {
			m.History.WriteProgress(currentWips)

			if len(prevWips) > 0 {
				areEqual := reflect.DeepEqual(currentWips, prevWips)

				if !areEqual {
					fmt.Println("Update found! Pushing notification. Next check at", c.Entries()[0].Next)

					// Get previous progress for existing works in progress
					wipsUpdate := make([]WorkInProgress, len(currentWips))
					copy(wipsUpdate, currentWips)
					for i := 0; i < len(wipsUpdate); i++ {
						currentWip := &wipsUpdate[i]
						for j := 0; j < len(prevWips); j++ {
							prevWip := &prevWips[j]
							if currentWip.Title == prevWip.Title && currentWip.Progress != prevWip.Progress {
								currentWip.PrevProgress = prevWip.Progress
							}
						}
					}

					if _, err := SendFCMUpdate(ctx, firebaseClient, wipsUpdate, config.ProgressTopic); err != nil {
						fmt.Println(err)
					}
				} else {
					fmt.Println("No update. Next check at", c.Entries()[0].Next)
				}
			}
		}
		prevWips = currentWips
	})
	fmt.Println("First check at", c.Entries()[0].Next)
	c.Start()

}

// SendFCMUpdate pushes an update via FCM
func SendFCMUpdate(ctx context.Context, firebaseClient *messaging.Client, wips []WorkInProgress, topic string) (string, error) {

	wipsStr, err := json.Marshal(wips)
	if err != nil {
		return "", err
	}

	oneHour := time.Duration(1) * time.Hour
	message := &messaging.Message{
		Topic: topic,
		Data: map[string]string{
			"worksInProgress": string(wipsStr),
		},
		Android: &messaging.AndroidConfig{
			TTL:      &oneHour,
			Priority: "normal",
			Notification: &messaging.AndroidNotification{
				Title: "Stormwatch",
				Body:  "Brandon Sanderson posted a progress update",
				Icon:  "ic_stat_ic_notification",
				Color: "#4195f4",
			},
		},
	}

	return firebaseClient.Send(ctx, message)
}
