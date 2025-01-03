package main

import (
	"log"
	"os"

	"github.com/Rhionin/SanderServer/progress"
)

var (
	slackWebhookURL = os.Getenv("SLACK_WEBHOOK_URL")
)

func main() {
	wips := []progress.WorkInProgress{
		{Title: "Book 1", Progress: 25},
		{Title: "Book 2 has a very long name copyedit and stuff", Progress: 50, PrevProgress: 30},
		{Title: "Book 3", Progress: 75},
		{Title: "Book 4", Progress: 100, PrevProgress: 80},
	}
	channelOverride := "#cjc-slack-testing"

	if err := progress.SendSlackUpdate(slackWebhookURL, wips, channelOverride); err != nil {
		log.Fatalf("Send slack update failed: %s", err)
	}
}
