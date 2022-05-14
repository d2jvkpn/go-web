package cmd

import (
	_ "embed"
	"fmt"
	"log"
	"os"

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
		fp string
		// name string
		fSet *pflag.FlagSet
	)

	command = &cobra.Command{
		Use:   "api-test",
		Short: "api test",
		Long: fmt.Sprintf(
			"http request test by provide api name and a yaml config:\n```yaml\n%s```",
			apiTestConfig,
		),

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

			for _, v := range args {
				statusCode, body, err = rt.Request(v)
				e := "nil"
				if err != nil {
					e = err.Error()
				}

				fmt.Printf(
					">>> Api: %s, StatusCode: %d, Error: %q\n%s\n\n",
					v, statusCode, e, body,
				)

				if err != nil {
					os.Exit(1)
				}
			}
		},
	}

	fSet = command.Flags()
	fSet.StringVar(&fp, "config", "", "yaml config file path")
	// fSet.StringVar(&name, "name", "", "api name defined in config file")

	_ = cobra.MarkFlagRequired(fSet, "config")
	// _ = cobra.MarkFlagRequired(fSet, "name")

	return command
}
