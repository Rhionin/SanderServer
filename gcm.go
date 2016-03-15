package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"rhionin.com/Rhionin/SanderServer/progress"
)

type gcmMessage struct {
	To   string                 `json:"to"`
	Data map[string]interface{} `json:"data"`
}

// SendGCMUpdate pushes an update via GCM
func SendGCMUpdate(wips []progress.WorkInProgress) {
	config := GetConfig()

	url := "https://android.googleapis.com/gcm/send"

	jsonMsg := gcmMessage{
		To: "/topics/progress",
		Data: map[string]interface{}{
			"worksInProgress": wips,
		},
	}
	fmt.Println(jsonMsg)

	jsonStr, _ := json.Marshal(jsonMsg)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

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
	fmt.Println("response Body:", string(body))
}
