package progress

import (
	"fmt"

	"reflect"

	"github.com/robfig/cron"
	cfg "github.com/Rhionin/SanderServer/config"
	"github.com/Rhionin/SanderServer/gcm"
)

type (
	// Monitor monitors progress
	Monitor struct {
		LiveReader Reader
		History    ReadWriter
	}

	// Reader reads progress
	Reader interface {
		GetProgress() []WorkInProgress
	}

	// Writer writes progress
	Writer interface {
		WriteProgress(wips []WorkInProgress) error
	}

	// ReadWriter read and writes progress
	ReadWriter interface {
		Reader
		Writer
	}
)

// ScheduleProgressCheckJob schedules a job to repeatedly check progress
// on Brandon Sanderson's books
func (m *Monitor) ScheduleProgressCheckJob() {
	prevWips := m.History.GetProgress()

	config := cfg.GetConfig()

	c := cron.New()
	fmt.Println(config.ProgressCheckInterval)
	c.AddFunc(config.ProgressCheckInterval, func() {
		currentWips := CheckProgress() // m.LiveReader.GetProgress()
		if len(currentWips) > 0 {
			m.History.WriteProgress(currentWips)

			if len(prevWips) > 0 {
				areEqual := reflect.DeepEqual(currentWips, prevWips)

				if !areEqual {
					fmt.Println("Update found! Pushing notification. Next check at", c.Entries()[0].Next)

					// Get previous progress for existing works in progress
					wipsUpdate := make([]WorkInProgress, len(currentWips))
					copy(wipsUpdate, currentWips)
					for i := 0; i < len(wipsUpdate); i++ {
						currentWip := &wipsUpdate[i]
						for j := 0; j < len(prevWips); j++ {
							prevWip := &prevWips[j]
							if currentWip.Title == prevWip.Title && currentWip.Progress != prevWip.Progress {
								currentWip.PrevProgress = prevWip.Progress
							}
						}
					}

					SendGCMUpdate(wipsUpdate, "/topics/progress")
				} else {
					fmt.Println("No update. Next check at", c.Entries()[0].Next)
				}
			}
		}
		prevWips = currentWips
	})
	fmt.Println("First check at", c.Entries()[0].Next)
	c.Start()

}

// SendGCMUpdate pushes an update via GCM
func SendGCMUpdate(wips []WorkInProgress, recipient string) {

	message := gcm.Message{
		To: recipient,
		Data: map[string]interface{}{
			"worksInProgress": wips,
		},
	}

	gcm.Send(message)
}
