package main

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"go.uber.org/zap"
)

var (
	debug *bool
)

func TestMain(m *testing.M) {
	debug = flag.Bool("debug", false, "Enable observability for debugging.")
	flag.Parse()

	fmt.Println(">>>")
	if *debug {
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		zap.ReplaceGlobals(logger)
	}
	os.Exit(m.Run())
}

func TestA1(t *testing.T) {
	fmt.Println("~~~", *debug)
}
