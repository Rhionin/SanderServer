package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Rhionin/SanderServer/config"
	"github.com/Rhionin/SanderServer/firebase"
	"github.com/Rhionin/SanderServer/progress"
)

const (
	historyFile                   = getenvOrDefault("HISTORY_FILE", "./history.txt")
	firebaseCredentialsConfigPath = getenvRequired("FIREBASE_CONFIG")
)

func main() {
	ctx := context.Background()
	firebaseClient, err := firebase.NewMessagingClient(ctx, firebaseCredentialsConfigPath)
	if err != nil {
		fmt.Println(err)
		panic("Failed to initialize Firebase messaging client")
	}

	history := progress.JSONFileReadWriter{
		FilePath: historyFile,
	}
	monitor := progress.Monitor{
		LiveReader: progress.GetProgress,
		History:    history,
		Config:     config.GetConfig(),
	}

	monitor.ScheduleProgressCheckJob(ctx, firebaseClient)
}

func getenvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getenvRequired(key) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		panic("Must provide " + key)
	}
	return value
}
