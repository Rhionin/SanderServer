package firebase

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
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
