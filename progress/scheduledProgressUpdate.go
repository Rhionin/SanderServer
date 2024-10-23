package progress

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/Rhionin/SanderServer/config"
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
		GetProgress() ([]WorkInProgress, error)
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

//go:embed error-page.html
var ErrorPageContent []byte

// ScheduleProgressCheckJob schedules a job to repeatedly check progress
// on Brandon Sanderson's books
func (m *Monitor) ScheduleProgressCheckJob(ctx context.Context, firebaseClient *messaging.Client) context.CancelFunc {
	prevWips, err := m.History.GetProgress()
	if err != nil {
		if errors.Is(err, ErrEmptyProgressFile) {
			log.Println("WARNING: No progress found in history")
		} else {
			log.Fatal("get progress:", err)
		}
	}

	if m.Config.SlackWebhookURL != "" {
		log.Println("Slack notifications enabled")
	}

	statusPagePublishingEnabled := m.Config.GithubUsername != "" && m.Config.GithubApiKey != ""
	if statusPagePublishingEnabled {
		log.Println("Status page publishing enabled")
	}

	ctx, cancel := context.WithCancel(ctx)

	log.Println("Progress check interval:", m.Config.ProgressCheckInterval)
	log.Println("First check at", time.Now().Add(m.Config.ProgressCheckInterval))

	ticker := time.NewTicker(m.Config.ProgressCheckInterval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				currentWips, err := m.LiveReader.GetProgress()
				if err != nil {
					log.Println("Failed to get progress:", err)
				}
				if len(currentWips) > 0 {
					m.History.WriteProgress(currentWips)

					if len(prevWips) > 0 {
						areEqual := reflect.DeepEqual(currentWips, prevWips)

						if !areEqual {
							log.Println("Update found! Pushing notification. Next check at", time.Now().Add(m.Config.ProgressCheckInterval))

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
								log.Println("Failed to send FCM update:", err)
							}
							if m.Config.SlackWebhookURL != "" {
								if err := SendSlackUpdate(m.Config.SlackWebhookURL, wipsUpdate); err != nil {
									log.Println(err)
								}
							}

							if statusPagePublishingEnabled {
								log.Println("Publishing status page update...")
								statusPageContent, err := CreateStatusPage(currentWips)
								if err != nil {
									log.Println("Failed to create status page:", err)
								} else {
									if err = PublishStatusPage(m.Config.GithubUsername, m.Config.GithubApiKey, statusPageContent); err != nil {
										log.Println("Failed to publish status page:", err)
									} else {
										log.Println("Status page update complete!")
									}
								}
							}

						} else {
							log.Println("No update. Next check at", time.Now().Add(m.Config.ProgressCheckInterval))
						}
					}
				} else {
					log.Println("No works in progress detected.")
					if statusPagePublishingEnabled {
						if err := PublishStatusPage(m.Config.GithubUsername, m.Config.GithubApiKey, ErrorPageContent); err != nil {
							log.Println("Failed to publish error status page:", err)
						} else {
							log.Println("Error status page publish complete!")
						}
					}
				}
				prevWips = currentWips
			}
		}
	}()

	return cancel
}

// SendFCMUpdate pushes an update via FCM
func SendFCMUpdate(ctx context.Context, firebaseClient *messaging.Client, wips []WorkInProgress, topic string) (string, error) {

	log.Println("Sending FCM message to topic "+topic, wips)

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
