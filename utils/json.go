package utils

import (
	"context"
	"io"
)

type JSON struct {
	Input  io.Reader
	Output io.Writer
	Simple bool
}

func (j *JSON) Execute(ctx context.Context) error {
	return nil
}
