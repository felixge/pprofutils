package main

import (
	"flag"
	"fmt"
	"os"

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
	fmt.Printf("Sample Types: %d\n", len(inProf.SampleType))
	fmt.Printf("Samples: %d\n", len(inProf.Sample))
	fmt.Printf("Locations: %d\n", len(inProf.Location))
	fmt.Printf("Functions: %d\n", len(inProf.Function))
	fmt.Printf("Mappings: %d\n", len(inProf.Mapping))
	fmt.Printf("Comments: %d\n", len(inProf.Comments))
	return nil
}
