package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/pprof/profile"
)

type JSON struct {
	Input  io.Reader
	Output io.Writer
	Simple bool
}

func (j *JSON) Execute(ctx context.Context) error {
	if j.Simple {
		return fmt.Errorf("simple format is not implemented yet")
	}
	prof, err := profile.Parse(j.Input)
	if err != nil {
		return err
	}
	return toFullJSON(prof, j.Output)
}

func toFullJSON(prof *profile.Profile, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(prof)
}
