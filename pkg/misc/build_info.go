package misc

import (
	"runtime/debug"
	"strings"
	"time"
)

func BuildInfo(vars ...[2]string) (info [][2]string) {
	info = make([][2]string, 0, 10)
	buildInfo, _ := debug.ReadBuildInfo()

	for i := range vars {
		info = append(info, [2]string{vars[i][0], vars[i][1]})
	}

	info = append(info, [2]string{"startTime", time.Now().Format(time.RFC3339)})
	info = append(info, [2]string{"goVersion", buildInfo.GoVersion})

	parseFlags := func(str string) {
		for _, v := range strings.Fields(str) {
			k, v, _ := strings.Cut(v, "=")
			if strings.HasPrefix(k, "main.") && v != "" {
				info = append(info, [2]string{k[5:], v})
			}
		}
	}

	for _, v := range buildInfo.Settings {
		if v.Key == "-ldflags" || v.Key == "--ldflags" {
			parseFlags(v.Value)
		}
	}

	return info
}
