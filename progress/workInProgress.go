package progress

import (
	"fmt"
)

// WorkInProgress represents each work and its progress
type WorkInProgress struct {
	Title        string `json:"title"`
	Progress     int    `json:"progress"`
	PrevProgress int    `json:"prevProgress"`
}

// ToString prints to string
func (wip *WorkInProgress) ToString() string {
	return fmt.Sprintf("%s (%d%%)", wip.Title, wip.Progress)
}
