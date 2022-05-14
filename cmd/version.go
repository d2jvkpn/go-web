package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type VersionInfo [][2]string

func (items VersionInfo) String() string {
	slice := make([]string, 0, len(items))

	for i := range items {
		slice = append(slice, fmt.Sprintf("%s: %s", items[i][0], items[i][1]))
	}
	return strings.Join(slice, "\n")
}

func (items VersionInfo) JSON() []byte {
	mp := make(map[string]string, len(items))
	for i := range items {
		mp[items[i][0]] = items[i][1]
	}

	bts, _ := json.Marshal(mp)
	return bts
}

func NewVersion(slice [][2]string) (command *cobra.Command) {
	var (
		jsonFmt bool
		fSet    *pflag.FlagSet
	)

	command = &cobra.Command{
		Use:   "version",
		Short: "version inforamtion",
		Long:  `version and build information`,

		Run: func(cmd *cobra.Command, args []string) {
			items := VersionInfo(slice)

			if jsonFmt {
				fmt.Printf("%s\n", items.JSON())
			} else {
				fmt.Printf("%s\n", items)
			}
		},
	}

	fSet = command.Flags()
	fSet.BoolVar(&jsonFmt, "json", false, "output command in json object")

	return command
}
