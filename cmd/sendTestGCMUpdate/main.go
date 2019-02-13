package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Rhionin/SanderServer/firebase"
	"github.com/Rhionin/SanderServer/progress"
)

func main() {
	firebaseCredentialsConfigPath := os.Getenv("FIREBASE_CONFIG")
	if len(firebaseCredentialsConfigPath) == 0 {
		panic("Must provide a FIREBASE_CONFIG path")
	}

	ctx := context.Background()
	firebaseClient, err := firebase.NewMessagingClient(ctx, firebaseCredentialsConfigPath)
	if err != nil {
		fmt.Println(err)
		panic("Failed to initialize Firebase messaging client")
	}

	wips := []progress.WorkInProgress{
		{Title: "Book 1", Progress: 25},
		{Title: "Book 2 has a very long name copyedit and stuff", Progress: 50, PrevProgress: 30},
		{Title: "Book 3", Progress: 75},
		{Title: "Book 4", Progress: 100, PrevProgress: 80},
	}

	response, err := progress.SendFCMUpdate(ctx, firebaseClient, wips, "/topics/devprogress")
	if err != nil {
		fmt.Printf("Error sending FCM update: %s\n", err)
	}

	fmt.Println(response)
}
