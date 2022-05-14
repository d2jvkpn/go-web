package cmd

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/d2jvkpn/goapp/pkg/misc"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	//go:embed api-test.yaml
	apiTestConfig string
)

func NewApiTest() (command *cobra.Command) {
	var (
		fp, name string
		fSet     *pflag.FlagSet
	)

	command = &cobra.Command{
		Use:   "api test",
		Short: "api test",
		Long:  `http request test by provide a yaml config and api name`,

		Run: func(cmd *cobra.Command, args []string) {
			var (
				statusCode int
				body       string
				err        error
				rt         *misc.RequestTmpls
			)
			if rt, err = misc.LoadRequestTmpls("api config", fp); err != nil {
				log.Fatalln(err)
			}

			if statusCode, body, err = rt.Request(name); err != nil {
				fmt.Printf("!!! %v\n", err)
			}
			fmt.Printf("StatusCode: %d, Body:\n  %s\n", statusCode, body)
		},
	}

	fSet = command.Flags()
	fSet.StringVar(&fp, "config", "", "yaml config file path")
	fSet.StringVar(&name, "name", "", "api name defined in config file")

	_ = cobra.MarkFlagRequired(fSet, "config")
	_ = cobra.MarkFlagRequired(fSet, "name")

	return command
}
