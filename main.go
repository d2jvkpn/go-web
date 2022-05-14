package main

import (
	_ "embed"
	"log"

	"github.com/d2jvkpn/goapp/cmd"
	"github.com/d2jvkpn/goapp/pkg/misc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:generate bash go_build.sh

var (
	//go:embed Project.yaml
	projectStr string
	buildTime  string
	gitBranch  string
	gitCommit  string
	gitTime    string
)

func init() {
}

func main() {
	var (
		err     error
		project *viper.Viper
	)

	if project, err = misc.ReadConfigString("project", projectStr, "yaml"); err != nil {
		log.Fatalln(err)
	}

	extract := func(key string) [2]string {
		return [2]string{key, project.GetString(key)}
	}

	root := &cobra.Command{Use: project.GetString("usage")}

	root.AddCommand(cmd.NewVersion(
		"version",
		misc.BuildInfo(extract("project"), extract("version")),
	))

	root.AddCommand(cmd.NewConfig("config"))
	root.AddCommand(cmd.NewServe("serve"))

	root.Execute()
}
