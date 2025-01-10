package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Rhionin/SanderServer/internal/progress"
	"github.com/justinrixx/retryhttp"
)

var (
	ErrNoWebhookURL      = errors.New("no webhook url")
	ErrNoProgressUpdates = errors.New("no progress updates")
)

type (
	UpdateClient struct {
		WebhookURL      string
		ChannelOverride string
	}

	slackAttachment struct {
		Color  string       `json:"color"`
		Title  string       `json:"title"`
		Text   string       `json:"text"`
		Fields []slackField `json:"fields"`
	}

	slackField struct {
		Title string `json:"title"`
		Value string `json:"value"`
		Short bool   `json:"short"`
	}

	slackPost struct {
		Channel     string            `json:"channel,omitempty"`
		Text        string            `json:"text"`
		Attachments []slackAttachment `json:"attachments"`
	}
)

func NewUpdateClient(webhookURL, channelOverride string) *UpdateClient {
	return &UpdateClient{
		WebhookURL:      webhookURL,
		ChannelOverride: channelOverride,
	}
}

func (client *UpdateClient) GetName() string {
	return "slack"
}

// SendSlackUpdate sends an update to slack
func (client *UpdateClient) SendUpdate(ctx context.Context, progressUpdates []progress.ProgressUpdate) error {
	if client.WebhookURL == "" {
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
		Channel: client.ChannelOverride,
		Text:    "*Brandon Sanderson has posted a progress update:*",
		Attachments: []slackAttachment{
			{
				Color:  "#007500",
				Fields: fields,
			},
		},
	})
	req, err := http.NewRequest(http.MethodPost, client.WebhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: retryhttp.New(),
	}
	resp, err := httpClient.Do(req)
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
