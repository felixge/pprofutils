package utils

import (
	"context"
	"io"

	"github.com/google/pprof/profile"
)

type Raw struct {
	Input  io.Reader
	Output io.Writer
}

func (r *Raw) Execute(ctx context.Context) error {
	prof, err := profile.Parse(r.Input)
	if err != nil {
		return err
	}
	_, _ = io.WriteString(r.Output, prof.String())
	return nil
}
