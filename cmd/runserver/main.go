package main

import (
	"flag"
	"fmt"
	cfg "rhionin.com/Rhionin/SanderServer/config"
	"rhionin.com/Rhionin/SanderServer/progress"
	"rhionin.com/Rhionin/SanderServer/server"
)

func main() {

	checkPgrsPtr := flag.Bool("checkPgrs", false, "Run progress check job")

	flag.Parse()

	config := cfg.GetConfig()

	if *checkPgrsPtr {
		fmt.Println("Progress will poll: " + config.ProgressCheckInterval)
		progress.ScheduleProgressCheckJob()
	}

	server.Start(config.Port)
}
