package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	//go:embed serve.yaml
	serveConfig string
)

func NewConfig(name string) (command *cobra.Command) {
	return &cobra.Command{
		Use:   name,
		Short: "print serve config file",
		Long:  "print serve config file (yaml format)",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(serveConfig)
		},
	}
}
