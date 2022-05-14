package main

import (
	_ "embed"
	"log"

	"github.com/d2jvkpn/goapp/cmd"
	"github.com/d2jvkpn/goapp/internal"
	"github.com/d2jvkpn/goapp/pkg/misc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:generate bash go_build.sh

var (
	//go:embed Project.yaml
	projectStr string
)

func init() {
}

func main() {
	var (
		err       error
		buildInfo [][2]string
		project   *viper.Viper
	)

	if project, err = misc.ReadConfigString("project", projectStr, "yaml"); err != nil {
		log.Fatalln(err)
	}

	extract := func(key string) [2]string {
		return [2]string{key, project.GetString(key)}
	}
	buildInfo = misc.BuildInfo(extract("project"), extract("version"))
	internal.BuildInfo = buildInfo

	root := &cobra.Command{Use: project.GetString("usage")}

	root.AddCommand(cmd.NewVersion(buildInfo))
	root.AddCommand(cmd.NewConfig("config"))
	root.AddCommand(cmd.NewServe())
	root.AddCommand(cmd.NewApiTest())

	root.Execute()
}
