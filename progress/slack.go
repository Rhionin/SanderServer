package progress

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type slackRequestBody struct {
	Blocks []block `json:"blocks"`
}

type block struct {
	Type string `json:"type"`
	Text text   `json:"text"`
}

type text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// SendSlackUpdate sends an update to slack
func SendSlackUpdate(webhookURL string, wips []WorkInProgress) error {

	msg := "Brandon Sanderson has posted a progress update:"
	for _, wip := range wips {
		progressStr := fmt.Sprintf("%d%%", wip.Progress)
		if wip.PrevProgress > 0 {
			progressStr = fmt.Sprintf("%d%% => %d%%", wip.PrevProgress, wip.Progress)
		}
		msg += fmt.Sprintf("\n• %s (%s)", wip.Title, progressStr)
	}

	slackBody, _ := json.Marshal(slackRequestBody{
		Blocks: []block{
			{
				Type: "section",
				Text: text{
					Type: "mrkdwn",
					Text: msg,
				},
			},
		},
	})
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack: " + buf.String())
	}
	return nil
}
