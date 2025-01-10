package main

import (
	"context"
	"log"
	"os"

	"github.com/Rhionin/SanderServer/internal/progress"
	"github.com/Rhionin/SanderServer/internal/slack"
)

var (
	slackWebhookURL = os.Getenv("SLACK_WEBHOOK_URL")
)

func main() {
	channelOverride := "#cjc-slack-testing"
	updateClient := slack.NewUpdateClient(slackWebhookURL, channelOverride)

	updates := []progress.ProgressUpdate{
		{Title: "Book 1", Progress: 25},
		{Title: "Book 2 has a very long name copyedit and stuff", Progress: 50, PrevProgress: 30},
		{Title: "Book 3", Progress: 75},
		{Title: "Book 4", Progress: 100, PrevProgress: 80},
	}

	if err := updateClient.SendUpdate(context.Background(), updates); err != nil {
		log.Fatalf("Send slack update failed: %s", err)
	}
}
