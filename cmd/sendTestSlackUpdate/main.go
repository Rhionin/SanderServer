package main

import (
	"github.com/Rhionin/SanderServer/config"
	"github.com/Rhionin/SanderServer/progress"
)

func main() {
	cfg := config.GetConfig("./cmd/config.yaml")

	wips := []progress.WorkInProgress{
		{Title: "Book 1", Progress: 25},
		{Title: "Book 2 has a very long name copyedit and stuff", Progress: 50, PrevProgress: 30},
		{Title: "Book 3", Progress: 75},
		{Title: "Book 4", Progress: 100, PrevProgress: 80},
	}

	if err := progress.SendSlackUpdate(cfg.SlackWebhookURL, wips); err != nil {
		panic(err)
	}
}
