package progress

import "fmt"

// WorkInProgress represents each work and its progress
type WorkInProgress struct {
	Title    string `json:"title"`
	Progress int    `json:"progress"`
}

func (wip *WorkInProgress) String() string {
	return fmt.Sprintf("%s (%d%%)", wip.Title, wip.Progress)
}
