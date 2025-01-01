package progress

import (
	"fmt"
)

// WorkInProgressSimple represents each work and its progress
type WorkInProgressSimple struct {
	Title    string `json:"title"`
	Progress int    `json:"progress"`
}

// WorkInProgress represents each work and its progress, and a comparison against the previous progress
type WorkInProgress struct {
	Title        string `json:"title"`
	Progress     int    `json:"progress"`
	PrevProgress int    `json:"prevProgress"`
}

// ToString prints to string
func (wip *WorkInProgress) ToString() string {
	return fmt.Sprintf("%s (%d%%)", wip.Title, wip.Progress)
}
