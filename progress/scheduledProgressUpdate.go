package progress

import (
	"fmt"
	"github.com/robfig/cron"
)

func ScheduleProgressCheckJob() {
	c := cron.New()
	c.AddFunc("0 */5 * * * *", func() {
		wips := CheckProgress()
		fmt.Println(wips)
	})
	c.Start()
}
