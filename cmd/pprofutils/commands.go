package main

import (
	"context"
	"io"

	"github.com/felixge/pprofutils/utils"
)

var utilCommands = []UtilCommand{
	{
		Name:       "json",
		ShortUsage: "[-simple] <input file> <output file>",
		ShortHelp:  "Converts from pprof to json and vice versa.",
		LongHelp: `The input and output file default to "-" which means stdin or stdout. If the` + "\n" +
			`input is pprof the output is json and for json inputs the output is pprof. This` + "\n" +
			`is automatically detected.`,
		Flags: map[string]UtilFlag{
			"simple": {false, "Use simplified JSON format."},
		},
		Execute: func(ctx context.Context, a *UtilArgs) error {
			return (&utils.JSON{
				Input:  a.Inputs[0],
				Output: a.Output,
				Simple: a.Flags["simple"].(bool),
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
	Inputs []io.Reader
	Output io.Writer
	Flags  map[string]interface{}
}

type UtilFlag struct {
	Default interface{}
	Usage   string
}
