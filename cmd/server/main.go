package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Rhionin/SanderServer/config"
	"github.com/Rhionin/SanderServer/firebase"
	"github.com/Rhionin/SanderServer/progress"
)

var (
	historyFile                   = getenvOrDefault("HISTORY_FILE", "./history.json")
	firebaseCredentialsConfigPath = getenvRequired("FIREBASE_CONFIG")
	configPath                    = getenvRequired("CONFIG")
)

func main() {
	ctx := context.Background()
	firebaseClient, err := firebase.NewMessagingClient(ctx, firebaseCredentialsConfigPath)
	if err != nil {
		log.Fatalf("initialize Firebase messaging client: %s", err)
	}

	historyFile, err := filepath.Abs(historyFile)
	if err != nil {
		log.Fatalf("get history file absolute path: %s", err)
	}
	log.Println("Writing history to", historyFile)

	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		_, err := os.Create(historyFile)
		if err != nil {
			log.Fatalf("create history file: %s", err)
		}
	}

	history := progress.JSONFileReadWriter{
		FilePath: historyFile,
	}
	monitor := progress.Monitor{
		LiveReader: progress.WebProgressChecker{
			URL: "http://brandonsanderson.com",
		},
		History: &history,
		Config:  config.GetConfig(configPath),
	}

	cancel := monitor.ScheduleProgressCheckJob(ctx, firebaseClient)
	defer cancel()

	waitForInterruptSignal()
	log.Println("exiting")
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

func waitForInterruptSignal() {
	// Go signal notification works by sending `os.Signal`
	// values on a channel. We'll create a channel to
	// receive these notifications (we'll also make one to
	// notify us when the program can exit).
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// This goroutine executes a blocking receive for
	// signals. When it gets one it'll print it out
	// and then notify the program that it can finish.
	go func() {
		sig := <-sigs
		log.Println()
		log.Println(sig)
		done <- true
	}()

	// The program will wait here until it gets the
	// expected signal (as indicated by the goroutine
	// above sending a value on `done`) and then exit.
	<-done
}
