package internal

import (
	"log"
)

func _SetupCrons() (err error) {
	// second, minute, hour, day(month), month, day(week)
	// at, jobname := "0 */1 * * * *", "JobName" // every minute, for testing only
	log.Printf(">>> Setup Cron %q: %s\n", Cron_Name, Cron_At)

	_, err = _Cron.AddFunc(Cron_At, func() {
		log.Printf(">>> Start Cron %q\n", Cron_Name)
	})

	if err != nil {
		return err
	}

	//... more cron jobs

	_Cron.Start()
	return nil
}
