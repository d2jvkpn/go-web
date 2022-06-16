package main

import (
	_ "embed"
	"log"
	"os"
	"strings"

	"github.com/d2jvkpn/go-web/internal"
	"github.com/d2jvkpn/go-web/internal/cmd"
	"github.com/d2jvkpn/go-web/pkg/misc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:generate bash scripts/build.sh

var (
	//go:embed project.yaml
	projectStr string
	//go:embed go.mod
	gomod string
)

func init() {
	misc.RegisterLogPrinter()
	strs := strings.Fields(strings.Split(gomod, "\n")[0])
	_ = os.Setenv("APP_Gomod", strs[len(strs)-1])
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
	internal.LoadBuildInfo(buildInfo)

	root := &cobra.Command{Use: project.GetString("usage")}

	root.AddCommand(cmd.NewVersion(buildInfo))
	root.AddCommand(cmd.NewConfig("config"))
	root.AddCommand(cmd.NewServe())
	root.AddCommand(cmd.NewApiTest())
	root.AddCommand(cmd.NewWsTest())

	root.Execute()
}
