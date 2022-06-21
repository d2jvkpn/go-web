package cmd

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	//go:embed config.yaml
	_Config string
	//go:embed deploy.yaml
	_Deploy string
)

func NewPrint(name string) (command *cobra.Command) {
	return &cobra.Command{
		Use:   name,
		Short: "print serve config files",
		Long:  "print serve config files: config.yaml, docker-compose.yaml, etc.",

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Fprintln(os.Stderr, "provide a object name: config, deploy")
				os.Exit(1)
			}

			switch args[0] {
			case "config":
				fmt.Println(_Config)
			case "deploy":
				fmt.Println(_Deploy)
			default:
				fmt.Fprintln(os.Stderr, "unknown object name")
				os.Exit(1)
			}
		},
	}
}
