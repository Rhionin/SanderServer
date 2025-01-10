package firebase

import (
	"context"
	"encoding/json"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/Rhionin/SanderServer/internal/progress"
	"google.golang.org/api/option"
)

// NewMessagingClient returns a new Firebase messaging client
func NewMessagingClient(ctx context.Context, firebaseCredentialsConfigPath string) (*messaging.Client, error) {
	opt := option.WithCredentialsFile(firebaseCredentialsConfigPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	return app.Messaging(ctx)
}

// SendFCMUpdate pushes an update via FCM
func SendFCMUpdate(ctx context.Context, firebaseClient *messaging.Client, wips []progress.ProgressUpdate, topic string) (string, error) {

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
