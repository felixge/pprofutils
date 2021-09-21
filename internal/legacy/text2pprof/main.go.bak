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
	versionF := flag.Bool("version", false, "Print version and exit.")
	flag.Parse()
	if *versionF {
		fmt.Printf("%s\n", internal.Version)
		return nil
	}
	outProf, err := pprofutils.Text{}.Convert(os.Stdin)
	if err != nil {
		return err
	}
	return outProf.Write(os.Stderr)
}
