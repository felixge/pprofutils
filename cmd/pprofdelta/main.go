package main

import (
	"errors"
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
	var (
		versionF = flag.Bool("version", false, "Print version and exit.")
		outF     = flag.String("o", "delta.pprof.gz", "Output profile name.")
	)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s -o <output> <profA> <profB>:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *versionF {
		fmt.Printf("%s\n", pprofutils.Version)
		return nil
	} else if flag.NArg() != 2 {
		flag.Usage()
		return errors.New("2 arguments required")
	}

	profA, err := os.Open(flag.Arg(0))
	if err != nil {
		return err
	}
	defer profA.Close()

	profB, err := os.Open(flag.Arg(1))
	if err != nil {
		return err
	}
	defer profB.Close()

	out, err := os.Create(*outF)
	if err != nil {
		return err
	}
	defer out.Close()

	return pprofutils.PPROFDelta(profA, profB, out)
}
