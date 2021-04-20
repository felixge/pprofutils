package pprofutils

import (
	"io"

	"github.com/google/pprof/profile"
)

func Delta(a, b io.Reader, w io.Writer) error {
	pa, err := profile.Parse(a)
	if err != nil {
		return err
	}
	pb, err := profile.Parse(b)
	if err != nil {
		return err
	}
	pa.Scale(-1)

	delta, err := profile.Merge([]*profile.Profile{pa, pb})
	if err != nil {
		return err
	} else if err := delta.CheckValid(); err != nil {
		return err
	}
	return delta.Write(w)
}
