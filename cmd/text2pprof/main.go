package main

import (
	"fmt"
	"os"

	pprofutil "github.com/felixge/pprof-util"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	return pprofutil.PPROF2Text(os.Stdin, os.Stdout)
}
