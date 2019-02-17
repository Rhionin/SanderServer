package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Rhionin/SanderServer/config"
	"github.com/Rhionin/SanderServer/firebase"
	"github.com/Rhionin/SanderServer/progress"
)

var (
	historyFile                   = getenvOrDefault("HISTORY_FILE", "./history.txt")
	firebaseCredentialsConfigPath = getenvRequired("FIREBASE_CONFIG")
	configPath                    = getenvRequired("CONFIG")
)

func main() {
	ctx := context.Background()
	firebaseClient, err := firebase.NewMessagingClient(ctx, firebaseCredentialsConfigPath)
	if err != nil {
		fmt.Println(err)
		panic("Failed to initialize Firebase messaging client")
	}

	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		if _, err := os.Create(historyFile); err != nil {
			panic(err)
		}
	}

	history := progress.JSONFileReadWriter{
		FilePath: historyFile,
	}
	monitor := progress.Monitor{
		LiveReader: progress.WebProgressChecker{},
		History:    &history,
		Config:     config.GetConfig(configPath),
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

func getenvRequired(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		panic("Must provide " + key)
	}
	return value
}
