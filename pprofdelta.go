package pprofutils

import (
	"io"

	"github.com/google/pprof/profile"
)

type DeltaConfig struct {
	// SampleTypes will limit delta calcultion to the given sample types. Other
	// sample types will retain the values of profile b.
	SampleTypes []SampleType
}

func (c DeltaConfig) Convert(a, b io.Reader, w io.Writer) error {
	pa, err := profile.Parse(a)
	if err != nil {
		return err
	}
	pb, err := profile.Parse(b)
	if err != nil {
		return err
	}

	ratios := make([]float64, len(pa.SampleType))
	for i, st := range pa.SampleType {
		// Empty c.SampleTypes means we calculate the delta for every st
		if len(c.SampleTypes) == 0 {
			ratios[i] = -1
			continue
		}

		// Otherwise we only calcuate the delta for any st that is listed in
		// c.SampleTypes. st's not listed in there will default to ratio 0, which
		// means we delete them from pa, so only the pb values remain in the final
		// profile.
		for _, deltaSt := range c.SampleTypes {
			if deltaSt.Type == st.Type && deltaSt.Unit == st.Unit {
				ratios[i] = -1
			}
		}
	}
	pa.ScaleN(ratios)

	delta, err := profile.Merge([]*profile.Profile{pa, pb})
	if err != nil {
		return err
	} else if err := delta.CheckValid(); err != nil {
		return err
	}
	return delta.Write(w)
}

func Delta(a, b io.Reader, w io.Writer) error {
	return DeltaConfig{}.Convert(a, b, w)
}
