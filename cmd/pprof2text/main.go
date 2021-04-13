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
	return pprofutils.PPROF2Text(os.Stdin, os.Stdout)
}
