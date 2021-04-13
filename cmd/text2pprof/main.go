package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/felixge/pprofutils"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	versionF := flag.Bool("version", false, "Print version and exit.")
	flag.Parse()
	if *versionF {
		fmt.Printf("%s\n", pprofutils.Version)
		return nil
	}
	return pprofutils.Text2PPROF(os.Stdin, os.Stdout)
}
