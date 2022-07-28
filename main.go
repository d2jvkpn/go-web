package main

import (
	_ "embed"
	"log"

	"github.com/d2jvkpn/go-web/internal"
	"github.com/d2jvkpn/go-web/internal/cmd"
	"github.com/d2jvkpn/go-web/pkg/misc"
	"github.com/d2jvkpn/go-web/pkg/wrap"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:generate bash scripts/go_build.sh

var (
	//go:embed project.yaml
	_Project string
)

func init() {
	misc.RegisterLogPrinter()
}

func main() {
	var (
		err       error
		buildInfo [][2]string
		project   *viper.Viper
	)

	if project, err = wrap.ReadConfigString("project", _Project, "yaml"); err != nil {
		log.Fatalln(err)
	}

	extract := func(key string) [2]string {
		return [2]string{key, project.GetString(key)}
	}
	buildInfo = misc.BuildInfo(extract("project"), extract("version"))
	internal.LoadBuildInfo(buildInfo)

	internal.AppendStaticDir(wrap.ServeStatic(
		"/uploads", "./data/uploads", false,
	))

	root := &cobra.Command{Use: project.GetString("usage")}

	root.AddCommand(
		cmd.NewServe(),
		cmd.NewApiTest(),
		cmd.NewWsTest(),
		cmd.NewVersion(buildInfo),
		cmd.NewPrint("print"),
		cmd.NewLog2Tsv(),
	)

	root.Execute()
}
