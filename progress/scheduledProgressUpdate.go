package progress

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/Rhionin/SanderServer/config"
	"github.com/robfig/cron"
)

type (
	// Monitor monitors progress
	Monitor struct {
		LiveReader Reader
		History    ReadWriter
		Config     config.Config
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

	if m.Config.SlackWebhookURL != "" {
		fmt.Println("Slack notifications enabled")
	}

	c := cron.New()
	fmt.Println(m.Config.ProgressCheckInterval)
	c.AddFunc(m.Config.ProgressCheckInterval, func() {
		currentWips := m.LiveReader.GetProgress()
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

					if _, err := SendFCMUpdate(ctx, firebaseClient, wipsUpdate, m.Config.ProgressTopic); err != nil {
						fmt.Println(err)
					}
					if m.Config.SlackWebhookURL != "" {
						if err := SendSlackUpdate(m.Config.SlackWebhookURL, wipsUpdate); err != nil {
							fmt.Println(err)
						}
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

	fmt.Println("Sending FCM message to topic "+topic, wips);

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
				Title:       "Stormwatch",
				Body:        "Brandon Sanderson posted a progress update",
				ClickAction: "FLUTTER_NOTIFICATION_CLICK",
			},
			CollapseKey: "progress_update",
		},
	}

	return firebaseClient.Send(ctx, message)
}
