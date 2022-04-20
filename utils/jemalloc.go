package utils

import (
	"bytes"
	"context"
	"io"

	"github.com/felixge/pprofutils/v2/internal/legacy"
)

type Jemalloc struct {
	Input  []byte
	Output io.Writer
}

func (f *Jemalloc) Execute(ctx context.Context) error {
	t := &legacy.Jemalloc{}
	prof, err := t.Convert(bytes.NewReader(f.Input))
	if err != nil {
		return err
	}
	return prof.Write(f.Output)
}
