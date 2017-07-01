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

	historyFilePath := "/var/log/SanderServer/history.json"

	if *checkPgrsPtr {
		fmt.Println("Progress will poll: " + config.ProgressCheckInterval)
		pm := progress.Monitor{
			History: &progress.JSONFileReadWriter{
				FilePath: historyFilePath,
			},
		}
		pm.ScheduleProgressCheckJob()
	}

	server.Start(config.Port)
}
