package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

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
		typesF   = flag.String("t", "", "A space separated list of type/unit sample types to calculate deltas for. Any other sample type will retain it's value from profB.")
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

	config := pprofutils.Delta{}
	if *typesF != "" {
		for _, sampleTypeS := range strings.Split(*typesF, " ") {
			parts := strings.Split(sampleTypeS, "/")
			if len(parts) != 2 {
				return fmt.Errorf("bad -t option: %q", *typesF)
			}
			config.SampleTypes = append(config.SampleTypes, pprofutils.SampleType{
				Type: parts[0],
				Unit: parts[1],
			})
		}
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

	return config.Convert(profA, profB, out)
}
