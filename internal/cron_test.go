package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCron(t *testing.T) {
	Cron_At = "*/10 * * * * *"

	err := _SetupCrons()
	require.NoError(t, err)
	_Cron.Start()

	time.Sleep(time.Minute)
	_Cron.Stop()
}
