package progress

import "fmt"

// ProgressUpdate represents each work and its progress, and a comparison against the previous progress
type ProgressUpdate struct {
	Title        string `json:"title"`
	Progress     int    `json:"progress"`
	PrevProgress int    `json:"prevProgress"`
}

func (pu *ProgressUpdate) String() string {
	progressStr := fmt.Sprintf("%d%%", pu.Progress)
	if pu.PrevProgress > 0 && pu.PrevProgress != pu.Progress {
		progressStr = fmt.Sprintf("%d%% => %d%%", pu.PrevProgress, pu.Progress)
	}
	return fmt.Sprintf("%s (%s)", pu.Title, progressStr)
}

// ClampToForwardProgress returns latest with each work's progress raised to at
// least its previously stored value. Transient backward dips from the source are
// ignored so they never register as a change. Order and title set follow latest,
// so a work removed from the source still drops off.
func ClampToForwardProgress(stored, latest []WorkInProgress) []WorkInProgress {
	storedByTitle := make(map[string]int, len(stored))
	for _, s := range stored {
		storedByTitle[s.Title] = s.Progress
	}

	clamped := make([]WorkInProgress, len(latest))
	for i, w := range latest {
		if prev, ok := storedByTitle[w.Title]; ok && prev > w.Progress {
			w.Progress = prev
		}
		clamped[i] = w
	}
	return clamped
}

func GetProgressUpdate(latestProgress, prevProgress []WorkInProgress) []ProgressUpdate {
	updates := make([]ProgressUpdate, len(latestProgress))

	for i, latest := range latestProgress {
		updates[i] = ProgressUpdate{
			Title:    latest.Title,
			Progress: latest.Progress,
		}
		for _, prev := range prevProgress {
			if latest.Title == prev.Title {
				updates[i].PrevProgress = prev.Progress
				break // Found a match, no need to keep searching
			}
		}
	}

	return updates
}
