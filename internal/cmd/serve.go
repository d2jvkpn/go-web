package cmd

import (
	// "fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/d2jvkpn/go-web/internal"
	"github.com/d2jvkpn/go-web/pkg/misc"

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
			var (
				value int
				err   error
				ch    chan int
			)

			if err = internal.Load(config, release); err != nil {
				log.Fatalln(err)
			}

			ch = make(chan int, 1)
			go misc.ListenOSSignal(ch) // goroutine1, send 0 when interrupted, send -1 otherwise

			go func() { // goroutine2
				parameters := map[string]interface{}{
					"config": config, "addr": addr, "release": release,
				}
				if err = internal.Serve(addr, parameters); err != nil {
					log.Println(err)
					ch <- 1
				} else {
					ch <- 2
				}
			}()

			value = <-ch
			switch value {
			case 1, 2: // goroutine2 exit
				internal.Down()
				ch <- value // send to goroutine1
				<-ch        // goroutine1 exit: -1
			case 0: // both goroutine1 exit by interrupted
				internal.Down()
				<-ch // goroutine2 exit: 1
			case -1: // both goroutine1 and goroutine2 exited
				internal.Down()
			}

			if value != 0 {
				os.Exit(1)
			}
		},
	}

	fSet = command.Flags()
	fSet.StringVar(&config, "config", filepath.Join("configs", "local.yaml"), "config file path")
	fSet.StringVar(&addr, "addr", ":8080", "http serve address")
	fSet.BoolVar(&release, "release", false, "run in release mode")

	return command
}
