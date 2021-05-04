package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/felixge/pprofutils"
	"github.com/felixge/pprofutils/internal"
	"github.com/google/pprof/profile"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var (
		versionF         = flag.Bool("version", false, "Print version and exit.")
		multiSampleTypes = flag.Bool("m", false, "Extract multiple sample types and write header for them.")
	)
	flag.Parse()
	if *versionF {
		fmt.Printf("%s\n", internal.Version)
		return nil
	}
	inProf, err := profile.Parse(os.Stdin)
	if err != nil {
		return err
	}
	return pprofutils.Protobuf{SampleTypes: *multiSampleTypes}.Convert(inProf, os.Stdout)
}
