package utils

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/google/pprof/profile"
)

type Labelframes struct {
	Input  []byte
	Output io.Writer
	Label  string
}

func (l *Labelframes) Execute(ctx context.Context) error {
	prof, err := profile.ParseData(l.Input)
	if err != nil {
		return err
	}

	var maxLocID uint64
	for _, loc := range prof.Location {
		if loc.ID > uint64(maxLocID) {
			maxLocID = loc.ID
		}
	}
	var maxFuncID uint64
	for _, fn := range prof.Function {
		if fn.ID > uint64(maxFuncID) {
			maxFuncID = fn.ID
		}
	}

	locIDX := map[string]*profile.Location{}
	for _, s := range prof.Sample {
		var labelVal = "N/A"
		for k, v := range s.Label {
			if k != l.Label {
				continue
			}

			labelVal = strings.Join(v, ",")
			break
		}

		frame := fmt.Sprintf("%s=%s", l.Label, labelVal)
		loc := locIDX[frame]
		if loc == nil {
			maxFuncID++
			fn := &profile.Function{
				ID:   maxFuncID,
				Name: frame,
			}
			prof.Function = append(prof.Function, fn)

			maxLocID++
			loc = &profile.Location{
				ID:   maxLocID,
				Line: []profile.Line{{Function: fn}},
			}
			prof.Location = append(prof.Location, loc)
			locIDX[frame] = loc
		}

		s.Location = append(s.Location, loc)
	}

	return prof.Write(l.Output)
}
