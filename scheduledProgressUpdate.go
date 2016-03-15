package main

import (
	"fmt"

	"github.com/robfig/cron"
	"reflect"
	"rhionin.com/Rhionin/SanderServer/progress"
)

// ScheduleProgressCheckJob schedules a job to repeatedly check progress
// on Brandon Sanderson's books
func ScheduleProgressCheckJob() {
	var prevWips []progress.WorkInProgress

	config := GetConfig()

	c := cron.New()
	fmt.Println(config.ProgressCheckInterval)
	c.AddFunc(config.ProgressCheckInterval, func() {
		currentWips := progress.CheckProgress()

		if len(prevWips) > 0 {
			areEqual := reflect.DeepEqual(currentWips, prevWips)

			if !areEqual {
				fmt.Println("Update found! Pushing notification. Next check at", c.Entries()[0].Next)
				SendGCMUpdate(currentWips)
			} else {
				fmt.Println("No update. Next check at", c.Entries()[0].Next)
				// currentWips = currentWips[0:1]
			}
		}
		prevWips = currentWips

	})
	fmt.Println("First check at", c.Entries()[0].Next)
	c.Start()

}
