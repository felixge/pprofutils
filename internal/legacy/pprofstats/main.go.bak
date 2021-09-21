package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/felixge/pprofutils/internal"
	"github.com/google/pprof/profile"
	"github.com/olekukonko/tablewriter"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var (
		versionF    = flag.Bool("version", false, "Print version and exit.")
		recursionsF = flag.Int("r", 0, "Print top N recursive locations.")
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

	recursionLocations := map[*profile.Location]int64{}
	var recursions int64
	for _, s := range inProf.Sample {
		seen := map[uint64]struct{}{}
		for _, l := range s.Location {
			if _, ok := seen[l.Address]; ok {
				recursions++
				recursionLocations[l] += 1
			} else {
				seen[l.Address] = struct{}{}
			}
		}
	}

	fmt.Printf("Sample Types: %d\n", len(inProf.SampleType))
	fmt.Printf("Samples: %d (%d locations doing %d recursions)\n", len(inProf.Sample), len(recursionLocations), recursions)
	fmt.Printf("Locations: %d\n", len(inProf.Location))
	fmt.Printf("Functions: %d\n", len(inProf.Function))
	fmt.Printf("Mappings: %d\n", len(inProf.Mapping))
	fmt.Printf("Comments: %d\n", len(inProf.Comments))

	if *recursionsF > 0 {
		fmt.Printf("\nTop %d Recursive Locations:\n\n", *recursionsF)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Count", "Method", "File", "Line"})
		table.SetBorder(false)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)

		for _, l := range topLocations(recursionLocations, *recursionsF) {
			line := l.Line[0]
			fn := line.Function
			table.Append([]string{
				fmt.Sprintf("%d", l.ID),
				fmt.Sprintf("%d", l.Count),
				fn.Name,
				fn.Filename,
				fmt.Sprintf("%d", line.Line),
			})
		}
		table.Render()
	}

	return nil
}

func printSample(s *profile.Sample) {
	var out []string
	for _, l := range s.Location {
		for _, line := range l.Line {
			out = append(out, line.Function.Name)
		}
	}
	fmt.Printf("%s\n", strings.Join(out, ";"))
}

func topLocations(locations map[*profile.Location]int64, n int) []TopLocation {
	var top []TopLocation
	for l, v := range locations {
		top = append(top, TopLocation{l, v})
	}
	sort.Slice(top, func(i, j int) bool {
		return top[i].Count > top[j].Count
	})
	return top[0:n]
}

type TopLocation struct {
	*profile.Location
	Count int64
}
