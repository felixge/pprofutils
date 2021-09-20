package main

import (
	"context"
	"io"

	"github.com/felixge/pprofutils/utils"
)

var commands = []Command{
	{
		Name:   "json",
		Inputs: 1,
		BoolFlags: map[string]BoolFlag{
			"simple": {false, "Use simplified JSON format."},
		},
		Execute: func(ctx context.Context, a *Args) error {
			return (&utils.JSON{
				Input:  a.Inputs[0],
				Output: a.Output,
				Simple: a.BoolFlags["simple"],
			}).Execute(ctx)
		},
	},
}

type Command struct {
	Name      string
	Inputs    int
	BoolFlags map[string]BoolFlag
	Execute   func(context.Context, *Args) error
}

type Args struct {
	Inputs    []io.Reader
	Output    io.Writer
	BoolFlags map[string]bool
}

type BoolFlag struct {
	Default bool
	Usage   string
}
