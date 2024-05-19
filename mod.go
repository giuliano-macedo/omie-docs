// NOTE: debug.ReadBuildInfo() only works when it is being ran on main package AND it is not unit tests (see: https://github.com/golang/go/issues/33976)
package mod

import (
	_ "embed"
	"strings"
)

const Version = "0.0.1"

//go:embed go.mod
var moduleContent string

func mustGetModuleUrl() string {
	lines := strings.Split(moduleContent, "\n")
	for _, line := range lines {
		lineStripped := strings.TrimSpace(line)
		if lineStripped == "" {
			continue
		}
		if strings.HasPrefix(line, "module") {
			_, url, _ := strings.Cut(line, " ")
			return url
		}
	}
	panic("Couldn't find module URL")
}

var Url = "https://" + mustGetModuleUrl()
