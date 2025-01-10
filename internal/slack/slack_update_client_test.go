package slack_test

import (
	"testing"

	"github.com/Rhionin/SanderServer/internal/slack"
	"github.com/Rhionin/SanderServer/internal/storminglambdas"
)

// Verify that UpdateClient implements the PushTarget interface
func TestSlackUpdateClientImplementsPushTargetInterface(t *testing.T) {
	var _ storminglambdas.PushTarget = (*slack.UpdateClient)(nil)
}
