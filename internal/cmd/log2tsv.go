package cmd

import (
	// "fmt"
	"log"
	"os"
	"time"

	"github.com/d2jvkpn/go-web/pkg/misc"
	"github.com/d2jvkpn/go-web/pkg/resp"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewLog2Tsv() (command *cobra.Command) {
	var (
		start, end string
		fSet       *pflag.FlagSet
	)

	command = &cobra.Command{
		Use:   "log2tsv",
		Short: "convert log to tsv format",
		Long:  `convert log to tsv format and output to stdout`,

		Run: func(cmd *cobra.Command, args []string) {
			var (
				startTime, endTime time.Time
				err                error
			)

			if len(args) == 0 {
				log.Fatalln("provide a log file path")
			}
			if start != "" {
				if startTime, err = misc.ParseDatetime(start); err != nil {
					log.Fatalln(err)
				}
			}
			if end != "" {
				if endTime, err = misc.ParseDatetime(end); err != nil {
					log.Fatalln(err)
				}
			}

			err = resp.Log2Tsv(args[0], os.Stdout, startTime, endTime)
			if err != nil {
				log.Fatalf("failed to process: %v\n", err)
			}
		},
	}

	fSet = command.Flags()

	fSet.StringVar(&start, "start", "", "start time of log")
	fSet.StringVar(&end, "end", "", "start time of log")

	return command
}
