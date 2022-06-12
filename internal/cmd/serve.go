package cmd

import (
	"log"
	"os"
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
			var err error

			if err = internal.Load(config, release); err != nil {
				log.Fatalln(err)
			}

			err = internal.Serve(addr)
			internal.Down()
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
		},
	}

	fSet = command.Flags()
	fSet.StringVar(&addr, "addr", ":8080", "http serve address")
	fSet.StringVar(&config, "config", filepath.Join("configs", "local.yaml"), "config file path")
	fSet.BoolVar(&release, "release", false, "run in release mode")

	return command
}
