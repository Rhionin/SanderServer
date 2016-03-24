package main

import (
	"flag"
	"fmt"
	"rhionin.com/Rhionin/SanderServer/server"
)

func main() {

	checkPgrsPtr := flag.Bool("checkPgrs", false, "Run progress check job")

	flag.Parse()

	config := GetConfig()

	server.Start(config.Port)

	if *checkPgrsPtr {
		fmt.Println("Progress will poll: " + config.ProgressCheckInterval)
		ScheduleProgressCheckJob()
	}
}
