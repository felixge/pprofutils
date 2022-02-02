package internal

import (
	"context"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/felixge/pprofutils/v2/utils"
)

const commonSuffix = "\n\n" + `The input and output file default to "-" which means stdin or stdout.`

var Utils = []Util{
	{
		Name:       "json",
		Flags:      map[string]UtilFlag{},
		ShortUsage: "<input file> <output file>",
		ShortHelp:  "Converts from pprof to json and vice versa",
		LongHelp: strings.TrimSpace(`
	Converts from pprof to json and vice vera. The input format is automatically
	detected and used to determine the output format.
	`) + commonSuffix,
		Examples: []Example{
			{Name: "Convert pprof to json", In: []string{"pprof"}, Out: []string{"json"}},
			{Name: "Convert json to pprof", In: []string{"json"}, Out: []string{"pprof"}},
		},
		Execute: func(ctx context.Context, a *UtilArgs) error {
			return (&utils.JSON{Input: a.Inputs[0], Output: a.Output}).Execute(ctx)
		},
	},
	{
		Name:       "raw",
		ShortUsage: "<input file> <output file>",
		ShortHelp:  "Converts pprof to the same text format as go tool pprof -raw",
		LongHelp: strings.TrimSpace(`
Converts pprof to the same text format as go tool pprof -raw.
`) + commonSuffix,
		Examples: []Example{
			{Name: "Convert pprof to raw", In: []string{"pprof"}, Out: []string{"txt"}},
		},
		Execute: func(ctx context.Context, a *UtilArgs) error {
			return (&utils.Raw{
				Input:  a.Inputs[0],
				Output: a.Output,
			}).Execute(ctx)
		},
	},
	{
		Name: "anon",
		Flags: map[string]UtilFlag{
			"whitelist": {"^runtime;^net;^encoding", "Semicolon separated pkg name regex list"},
		},
		ShortUsage: "[-whitelist=<regex>] <input file> <output file>",
		ShortHelp:  "Anonymizes a pprof profile",
		LongHelp: strings.TrimSpace(`
Takes a pprof profile and anonymizes it by replacing pkg, file and function
names with human readable hashes. The whitelist can be used to prevent certain
packages from being anonymized.

TODO: Ignore all stdlib packages by default and maybe also popular OSS libs.
`) + commonSuffix,
		Examples: []Example{
			{Name: "Anonymize a CPU profile", In: []string{"pprof", "png"}, Out: []string{"pprof", "png"}},
		},
		Execute: func(ctx context.Context, a *UtilArgs) error {
			return (&utils.Anon{
				Input:     a.Inputs[0],
				Output:    a.Output,
				Whitelist: a.Flags["whitelist"].(string),
			}).Execute(ctx)
		},
	},
	{
		Name:       "avg",
		ShortUsage: "<input file> <output file>",
		ShortHelp:  "Creates a profile with the average value per sample",
		LongHelp: strings.TrimSpace(`
Takes a block or mutex profile and creates a profile that contains the average
time per contention by dividing the nanoseconds or value in the profile by the
sample count value.

TODO: Support memory profiles.
`) + commonSuffix,
		Examples: []Example{
			{Name: "Convert block profile to avg time", In: []string{"pprof", "png"}, Out: []string{"pprof", "png"}},
		},
		Execute: func(ctx context.Context, a *UtilArgs) error {
			return (&utils.Avg{
				Input:  a.Inputs[0],
				Output: a.Output,
			}).Execute(ctx)
		},
	},
	{
		Name: "folded",
		Flags: map[string]UtilFlag{
			"headers": {false, "Add header column for each sample type"},
		},
		ShortUsage: "[-headers] <input file> <output file>",
		ShortHelp:  "Converts pprof to Brendan Gregg's folded text format and vice versa",
		LongHelp: strings.TrimSpace(`
Converts pprof to Brendan Gregg's folded text format and vice versa. The input
format is automatically detected and used to determine the output format.
`) + commonSuffix,
		Examples: []Example{
			{Name: "Convert folded text to pprof", In: []string{"txt"}, Out: []string{"pprof", "png"}},
			{Name: "Convert pprof to folded text", In: []string{"pprof", "png"}, Out: []string{"txt"}},
		},
		Execute: func(ctx context.Context, a *UtilArgs) error {
			return (&utils.Folded{
				Input:   a.Inputs[0],
				Output:  a.Output,
				Headers: a.Flags["headers"].(bool),
			}).Execute(ctx)
		},
	},
	{
		Name: "labelframes",
		Flags: map[string]UtilFlag{
			"label": {"mylabel", "The label key to turn into virtual frames."},
		},
		ShortUsage: "-label=<label> <input file> <output file>",
		ShortHelp:  "Adds virtual root frames for the given pprof label",
		LongHelp: strings.TrimSpace(`
Adds virtual root frames for the each value of the selected pprof label. This
is useful to visualize label values in a flamegraph.
`) + commonSuffix,
		Examples: []Example{
			{Name: "Add root frames for pprof label values", In: []string{"pprof", "png"}, Out: []string{"pprof", "png"}},
		},
		Execute: func(ctx context.Context, a *UtilArgs) error {
			return (&utils.Labelframes{
				Input:  a.Inputs[0],
				Output: a.Output,
				Label:  a.Flags["label"].(string),
			}).Execute(ctx)
		},
	},
	{
		Name: "heapage",
		Flags: map[string]UtilFlag{
			"period": {10 * time.Second, "The time period covered by the heap profile."},
		},
		ShortUsage: "-period=<period> <input file> <output file>",
		ShortHelp:  "Adds virtual frames showing the average allocation lifetime for Go memory allocations.",
		LongHelp: strings.TrimSpace(`
Adds virtual frames showing the average allocation lifetime for Go memory allocations.
`) + commonSuffix,
		Examples: []Example{
			{Name: "Calculate Avg Inuse Object Age", In: []string{"pprof", "png"}, Out: []string{"pprof", "png"}},
		},
		Execute: func(ctx context.Context, a *UtilArgs) error {
			return (&utils.Heapage{
				Input:  a.Inputs[0],
				Output: a.Output,
				Period: a.Flags["period"].(time.Duration),
			}).Execute(ctx)
		},
	},
}

func init() {
	sort.Slice(Utils, func(i, j int) bool {
		return Utils[i].Name < Utils[j].Name
	})
}

type Example struct {
	Name string
	In   []string
	Out  []string
}

type Util struct {
	Name       string
	Flags      map[string]UtilFlag
	ShortUsage string
	ShortHelp  string
	LongHelp   string
	Examples   []Example
	Execute    func(context.Context, *UtilArgs) error
}

type UtilArgs struct {
	Inputs [][]byte
	Output io.Writer
	Flags  map[string]interface{}
}

type UtilFlag struct {
	Default interface{}
	Usage   string
}
