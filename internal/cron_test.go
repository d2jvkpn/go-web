package internal

import (
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	Cron_At = "*/10 * * * * *"

	if err := _SetupCrons(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Minute)
	_Cron.Stop()
}
