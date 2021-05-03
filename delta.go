package pprofutils

import (
	"github.com/google/pprof/profile"
)

type Delta struct {
	// SampleTypes will limit delta calcultion to the given sample types. Other
	// sample types will retain the values of profile b.
	SampleTypes []SampleType
}

// TODO modifies a
func (d Delta) Convert(a, b *profile.Profile) (*profile.Profile, error) {
	ratios := make([]float64, len(a.SampleType))
	for i, st := range a.SampleType {
		// Empty c.SampleTypes means we calculate the delta for every st
		if len(d.SampleTypes) == 0 {
			ratios[i] = -1
			continue
		}

		// Otherwise we only calcuate the delta for any st that is listed in
		// c.SampleTypes. st's not listed in there will default to ratio 0, which
		// means we delete them from pa, so only the pb values remain in the final
		// profile.
		for _, deltaSt := range d.SampleTypes {
			if deltaSt.Type == st.Type && deltaSt.Unit == st.Unit {
				ratios[i] = -1
			}
		}
	}
	a.ScaleN(ratios)

	delta, err := profile.Merge([]*profile.Profile{a, b})
	if err != nil {
		return nil, err
	}
	return delta, delta.CheckValid()
}
