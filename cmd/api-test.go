package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/d2jvkpn/goapp/pkg/misc"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	apiTestConfig string
)

func NewApiTest() (command *cobra.Command) {
	var (
		fp   string
		fSet *pflag.FlagSet
		tmp  *misc.RequestTmpls
	)

	command = &cobra.Command{
		Use:   "api-test",
		Short: "api test",
		Long: fmt.Sprintf(
			"http request test by provide api name and a yaml config:\n```yaml\n%s```",
			tmp.Yaml(),
		),

		Run: func(cmd *cobra.Command, args []string) {
			var (
				statusCode   int
				body, errStr string
				err          error
				tmpls        []*misc.RequestTmpl
				rt           *misc.RequestTmpls
			)

			if rt, err = misc.LoadRequestTmpls("api config", fp); err != nil {
				log.Fatalln(err)
			}

			if tmpls, err = rt.Match(args...); err != nil {
				log.Fatal(err)
			}

			errStr = "nil"
			for _, v := range tmpls {
				if statusCode, body, err = rt.Request(v); err != nil {
					errStr = err.Error()
				}

				fmt.Printf(
					">>> API: %s, Path: %s, StatusCode: %d, Error: %q\n%s\n\n",
					v.Name, v.Path, statusCode, errStr, body,
				)

				if err != nil {
					os.Exit(1)
				}
			}
		},
	}

	fSet = command.Flags()
	fSet.StringVar(&fp, "config", "", "yaml config file path")
	_ = cobra.MarkFlagRequired(fSet, "config")

	return command
}
