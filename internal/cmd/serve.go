package cmd

import (
	// "fmt"
	"path/filepath"

	"github.com/d2jvkpn/go-web/internal"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewServe() (command *cobra.Command) {
	var (
		config, addr string
		release      bool
		fSet         *pflag.FlagSet
	)

	command = &cobra.Command{
		Use:   "serve",
		Short: "serve http",
		Long:  `running serve http`,

		Run: func(cmd *cobra.Command, args []string) {
			internal.Run(config, addr, release)
		},
	}

	fSet = command.Flags()
	fSet.StringVar(&config, "config", filepath.Join("configs", "local.yaml"), "config file path")
	fSet.StringVar(&addr, "addr", ":8080", "http serve address")
	fSet.BoolVar(&release, "release", false, "run in release mode")

	return command
}
