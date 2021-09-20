package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/google/pprof/profile"
)

type JSON struct {
	Input  []byte
	Output io.Writer
	Simple bool
}

func (j *JSON) Execute(ctx context.Context) error {
	if j.Simple {
		return fmt.Errorf("simple format is not implemented yet")
	}
	prof, err := profile.ParseData(j.Input)
	if err == nil {
		return toFullJSON(prof, j.Output)
	}
	if err := fromFullJSON(j.Input, j.Output); err != nil {
		return errors.New("input format is neither pprof nor json")
	}
	return nil
}

func toFullJSON(prof *profile.Profile, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(prof)
}

func fromFullJSON(in []byte, out io.Writer) error {
	prof := &profile.Profile{}
	if err := json.Unmarshal(in, prof); err != nil {
		return err
	}
	return prof.Write(out)
}
