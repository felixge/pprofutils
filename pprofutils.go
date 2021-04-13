package pprofutils

import (
	_ "embed"
	"strings"
)

var (
	//go:embed version.txt
	version string
	Version string = strings.TrimSpace(version)
)
