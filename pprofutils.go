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

type SampleType struct {
	Type string
	Unit string
}
