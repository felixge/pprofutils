package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/felixge/pprofutils"
	"github.com/felixge/pprofutils/internal"
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
	return pprofutils.Protobuf{SampleTypes: *multiSampleTypes}.Convert(os.Stdin, os.Stdout)
}
