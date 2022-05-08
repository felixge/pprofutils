package utils

import (
	"bytes"
	"context"
	"io"

	"github.com/felixge/pprofutils/v2/internal/legacy"
	"github.com/google/pprof/profile"
)

type Folded struct {
	Input       []byte
	Output      io.Writer
	Headers     bool
	LineNumbers bool
}

func (f *Folded) Execute(ctx context.Context) error {
	prof, err := profile.ParseData(f.Input)
	if err == nil {
		p := legacy.Protobuf{
			SampleTypes: f.Headers,
			LineNumbers: f.LineNumbers,
		}
		return p.Convert(prof, f.Output)
	}

	// Invalid pprof, assume it's folded text.
	t := &legacy.Text{}
	prof, err = t.Convert(bytes.NewReader(f.Input))
	if err != nil {
		return err
	}
	return prof.Write(f.Output)
}
