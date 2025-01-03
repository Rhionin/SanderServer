package progress

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var (
	ErrNoWebhookURL      = errors.New("no webhook url")
	ErrNoProgressUpdates = errors.New("no progress updates")
)

type slackAttachment struct {
	Color  string       `json:"color"`
	Title  string       `json:"title"`
	Text   string       `json:"text"`
	Fields []slackField `json:"fields"`
}

type slackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type slackPost struct {
	Channel     string            `json:"channel,omitempty"`
	Text        string            `json:"text"`
	Attachments []slackAttachment `json:"attachments"`
}

// SendSlackUpdate sends an update to slack
func SendSlackUpdate(webhookURL string, progressUpdates []ProgressUpdate, channelOverride string) error {
	if webhookURL == "" {
		return ErrNoWebhookURL
	}
	if len(progressUpdates) == 0 {
		return ErrNoProgressUpdates
	}

	fields := []slackField{}
	for _, update := range progressUpdates {
		fields = append(fields, slackField{
			Value: update.String(),
		})
	}

	slackBody, _ := json.Marshal(slackPost{
		Channel: channelOverride,
		Text:    "*Brandon Sanderson has posted a progress update:*",
		Attachments: []slackAttachment{
			{
				Color:  "#007500",
				Fields: fields,
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
