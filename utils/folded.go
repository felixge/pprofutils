package utils

import (
	"context"
	"io"

	"github.com/felixge/pprofutils/internal/legacy"
	"github.com/google/pprof/profile"
)

type Folded struct {
	Input  []byte
	Output io.Writer
}

func (f *Folded) Execute(ctx context.Context) error {
	prof, err := profile.ParseData(f.Input)
	if err != nil {
		return err
	}
	p := legacy.Protobuf{}
	return p.Convert(prof, f.Output)
}
