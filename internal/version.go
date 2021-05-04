package internal

import (
	_ "embed"
	"strings"
)

var (
	//go:embed version.txt
	version string
	// Version holds the version number of the package.
	Version string = strings.TrimSpace(version)
)
