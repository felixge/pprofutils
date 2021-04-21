package pprofutils

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
		if strings.TrimSpace(line) == "" {
			continue
		}

		// custom extension: first line can contain header that looks like this:
		// "samples/count duration/nanoseconds" to describe the sample types
		if n == 0 && !containsDigit(line) {
			p.SampleType = nil
			for _, sampleType := range strings.Split(line, " ") {
				parts := strings.Split(sampleType, "/")
				if len(parts) != 2 {
					return fmt.Errorf("bad header: %d: %q", n, line)
				}
				p.SampleType = append(p.SampleType, &profile.ValueType{
					Type: parts[0],
					Unit: parts[1],
				})
			}
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != len(p.SampleType)+1 {
			return fmt.Errorf("bad line: %d: %q", n, line)
		}

		stack := strings.Split(parts[0], ";")
		sample := &profile.Sample{}
		for _, valS := range parts[1:] {
			val, err := strconv.ParseInt(valS, 10, 64)
			if err != nil {
				return fmt.Errorf("bad line: %d: %q: %w", n, line, err)
			}
			sample.Value = append(sample.Value, val)
		}

		for i := range stack {
			frame := stack[len(stack)-i-1]
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

func containsDigit(s string) bool {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}
	return false
}
