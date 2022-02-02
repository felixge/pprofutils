package utils

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/pprof/profile"
)

// Heapage uses Little's Law to add virtual heap age frames to a Go memory
// profile. Credit for this idea goes to Maxim Sokolov.
//
// avg(inuse_age) = inuse_objects / (alloc_objects/period)
type Heapage struct {
	Input  []byte
	Output io.Writer
	Period time.Duration
}

func (h *Heapage) Execute(ctx context.Context) error {
	prof, err := profile.ParseData(h.Input)
	if err != nil {
		return err
	}

	var (
		inuseObjects = profile.ValueType{Type: "inuse_objects", Unit: "count"}
		allocObjects = profile.ValueType{Type: "alloc_objects", Unit: "count"}
		inuseIDX     = sampleTypeIndex(prof, inuseObjects)
		allocIDX     = sampleTypeIndex(prof, allocObjects)
		fnID         = maxFuncID(prof)
		locID        = maxLocID(prof)
	)
	if inuseIDX < 0 {
		return fmt.Errorf("missing sample type: %s/%s", inuseObjects.Type, inuseObjects.Unit)
	} else if allocIDX < 0 {
		return fmt.Errorf("missing sample type: %s/%s", allocObjects.Type, allocObjects.Unit)
	}

	for _, s := range prof.Sample {
		allocs := s.Value[allocIDX]

		var ageS string
		if allocs > 0 {
			rate := float64(allocs) / float64(h.Period.Nanoseconds())
			age := time.Duration(float64(s.Value[inuseIDX]) / rate)
			ageS = age.Truncate(time.Second / 10).String()
		} else {
			ageS = fmt.Sprintf("∞ (> %s)", h.Period.String())
		}

		fnID++
		locID++
		fn := &profile.Function{
			ID:   fnID,
			Name: ageS,
		}
		loc := &profile.Location{
			ID: locID,
			Line: []profile.Line{{
				Function: fn,
			}},
		}
		prof.Location = append(prof.Location, loc)
		prof.Function = append(prof.Function, fn)
		s.Location = append([]*profile.Location{loc}, s.Location...)
	}

	return prof.Write(h.Output)
}

func maxFuncID(prof *profile.Profile) (max uint64) {
	for _, f := range prof.Function {
		if f.ID > max {
			max = f.ID
		}
	}
	return
}

func maxLocID(prof *profile.Profile) (max uint64) {
	for _, loc := range prof.Location {
		if loc.ID > max {
			max = loc.ID
		}
	}
	return
}

func stack(s *profile.Sample) string {
	var funcs []string
	for _, loc := range s.Location {
		for _, line := range loc.Line {
			funcs = append(funcs, line.Function.Name)
		}
	}
	var labels []string
	for k, v := range s.NumLabel {
		labels = append(labels, fmt.Sprintf("%v=%v", k, v))
	}

	return strings.Join(labels, ",") + " " + strings.Join(funcs, "<-")
}

func sampleTypeIndex(prof *profile.Profile, vt profile.ValueType) int {
	for i, ovt := range prof.SampleType {
		if vt.Type == ovt.Type && vt.Unit == ovt.Unit {
			return i
		}
	}
	return -1
}
