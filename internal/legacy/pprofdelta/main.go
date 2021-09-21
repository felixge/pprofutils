package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

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
		fmt.Printf("%s\n", internal.Version)
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
			config.SampleTypes = append(config.SampleTypes, pprofutils.ValueType{
				Type: parts[0],
				Unit: parts[1],
			})
		}
	}

	profA, err := loadProfile(flag.Arg(0))
	if err != nil {
		return err
	}

	profB, err := loadProfile(flag.Arg(1))
	if err != nil {
		return err
	}

	out, err := os.Create(*outF)
	if err != nil {
		return err
	}
	defer out.Close()

	delta, err := config.Convert(profA, profB)
	if err != nil {
		return err
	}
	return delta.Write(out)
}

func loadProfile(path string) (*profile.Profile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return profile.Parse(file)
}
