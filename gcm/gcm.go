package gcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	cfg "github.com/Rhionin/SanderServer/config"
)

// Message a GCM message
type Message struct {
	To   string                 `json:"to"`
	Data map[string]interface{} `json:"data"`
}

// Send pushes an update via GCM
func Send(message Message) {
	config := cfg.GetConfig()

	url := "https://android.googleapis.com/gcm/send"

	fmt.Println("Sending GCM message", message)
	jsonStr, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "key="+config.GoogleAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("GCM response Body:", string(body))
}
