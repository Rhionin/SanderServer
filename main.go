package main

import (
	"rhionin.com/Rhionin/SanderServer/progress"
	"time"
)

func main() {
	progress.ScheduleProgressCheckJob()
	time.Sleep(10000 * time.Minute)
}
