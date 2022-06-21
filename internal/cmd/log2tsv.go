package cmd

import (
	"fmt"
	"os"

	"github.com/d2jvkpn/go-web/pkg/resp"

	"github.com/spf13/cobra"
)

func NewLog2Tsv() (command *cobra.Command) {

	command = &cobra.Command{
		Use:   "log2tsv",
		Short: "convert log to tsv format",
		Long:  `convert log to tsv format and output to stdout`,

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Fprintf(os.Stderr, "provide a log file path\n")
				os.Exit(1)
			}

			if err := resp.Log2Tsv(args[0], os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "failed to process: %v\n", err)
				os.Exit(1)
			}
		},
	}

	return command
}
