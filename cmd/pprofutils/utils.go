package main

import (
	"context"
	"io"
	"strings"

	"github.com/felixge/pprofutils/utils"
)

const commonSuffix = "\n\n" + `The input and output file default to "-" which means stdin or stdout.`

var utilCommands = []UtilCommand{
	{
		Name:       "json",
		ShortUsage: "<input file> <output file>",
		ShortHelp:  "Converts from pprof to json and vice versa",
		LongHelp: strings.TrimSpace(`
Converts from pprof to json and vice vera. The input format is automatically
detected and used to determine the output format.
`) + commonSuffix,
		Execute: func(ctx context.Context, a *UtilArgs) error {
			return (&utils.JSON{
				Input:  a.Inputs[0],
				Output: a.Output,
			}).Execute(ctx)
		},
	},
	{
		Name:       "raw",
		ShortUsage: "<input file> <output file>",
		ShortHelp:  "Converts pprof to the same text format as go tool pprof -raw",
		LongHelp: strings.TrimSpace(`
Converts pprof to the same text format as go tool pprof -raw
`) + commonSuffix,
		Execute: func(ctx context.Context, a *UtilArgs) error {
			return (&utils.Raw{
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
			"label": {"", "The label key to turn into virtual frames."},
		},
		ShortUsage: "-label=<label> <input file> <output file>",
		ShortHelp:  "Adds virtual root frames for the given pprof label",
		LongHelp: strings.TrimSpace(`
Adds virtual root frames for the given pprof label. This is useful to visualize
label values in a flamegraph.
`) + commonSuffix,
		Execute: func(ctx context.Context, a *UtilArgs) error {
			return (&utils.Labelframes{
				Input:  a.Inputs[0],
				Output: a.Output,
				Label:  a.Flags["label"].(string),
			}).Execute(ctx)
		},
	},
}

type UtilCommand struct {
	Name       string
	ShortUsage string
	ShortHelp  string
	LongHelp   string
	Flags      map[string]UtilFlag
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
