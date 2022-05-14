package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	//go:embed config.yaml
	configStr string
)

func NewConfig(name string) (command *cobra.Command) {
	return &cobra.Command{
		Use:   name,
		Short: "config",
		Long:  `config file demo`,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(configStr)
		},
	}
}
