package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/Rhionin/SanderServer/internal/firebase"
	"github.com/Rhionin/SanderServer/internal/progress"
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

	someConstant, err := strconv.Atoi(os.Getenv("EXTRA_CONSTANT"))
	if err != nil {
		someConstant = 0
	}
	wips := []progress.ProgressUpdate{
		{Title: "Book 1", Progress: 25},
		{Title: "Book 2 has a very long name copyedit and stuff", Progress: 50 + someConstant, PrevProgress: 30},
		{Title: "Book 3", Progress: 75},
		{Title: "Book 4", Progress: 100, PrevProgress: 80},
	}

	response, err := firebase.SendFCMUpdate(ctx, firebaseClient, wips, "flutter_devprogress")
	if err != nil {
		fmt.Printf("Error sending flutter FCM update: %s\n", err)
	}
	fmt.Println(response)
}
