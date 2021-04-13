package main

import (
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
	return pprofutils.Text2PPROF(os.Stdin, os.Stdout)
}
