package pprofutil

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/google/pprof/profile"
)

func Text2PPROF(text io.Reader, pprof io.Writer) error {
	var (
		functionID = uint64(1)
		locationID = uint64(1)
		p          = &profile.Profile{
			TimeNanos: time.Now().UnixNano(),
			SampleType: []*profile.ValueType{{
				Type: "samples",
				Unit: "count",
			}},
		}
	)

	m := &profile.Mapping{ID: 1, HasFunctions: true}
	p.Mapping = []*profile.Mapping{m}

	lines, err := io.ReadAll(text)
	if err != nil {
		return err
	}
	for n, line := range strings.Split(string(lines), "\n") {
		i := strings.LastIndex(line, " ")
		if i <= 0 {
			return fmt.Errorf("bad line: %d: %q", n, line)
		}
		stack := strings.Split(line[0:i], ";")
		count, err := strconv.ParseInt(line[i+1:], 10, 64)
		if err != nil {
			return fmt.Errorf("bad line: %d: %q: %w", n, line, err)
		}

		sample := &profile.Sample{Value: []int64{count}}
		for _, frame := range stack {
			function := &profile.Function{
				ID:   functionID,
				Name: frame,
			}
			p.Function = append(p.Function, function)
			functionID++

			location := &profile.Location{
				ID:      locationID,
				Mapping: m,
				Line:    []profile.Line{{Function: function}},
			}
			p.Location = append(p.Location, location)
			locationID++

			sample.Location = append(sample.Location, location)
		}

		p.Sample = append(p.Sample, sample)
	}

	if err := p.CheckValid(); err != nil {
		return err
	} else if err := p.Write(pprof); err != nil {
		return err
	}
	return nil
}
