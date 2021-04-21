package pprofutils

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/google/pprof/profile"
)

type PPROF2TextConfig struct {
	SampleTypes bool
}

func (c PPROF2TextConfig) Convert(pprof io.Reader, text io.Writer) error {
	prof, err := profile.Parse(pprof)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(text)
	if c.SampleTypes {
		var sampleTypes []string
		for _, sampleType := range prof.SampleType {
			sampleTypes = append(sampleTypes, sampleType.Type+"/"+sampleType.Unit)
		}
		w.WriteString(strings.Join(sampleTypes, " ") + "\n")
	}
	for _, sample := range prof.Sample {
		var frames []string
		for i := range sample.Location {
			loc := sample.Location[len(sample.Location)-i-1]
			for j := range loc.Line {
				line := loc.Line[len(loc.Line)-j-1]
				frames = append(frames, line.Function.Name)
			}
		}
		var values []string
		for _, val := range sample.Value {
			values = append(values, fmt.Sprintf("%d", val))
			if !c.SampleTypes {
				break
			}
		}
		fmt.Fprintf(
			w,
			"%s %s\n",
			strings.Join(frames, ";"),
			strings.Join(values, " "),
		)
	}
	return w.Flush()
}

func PPROF2Text(pprof io.Reader, text io.Writer) error {
	return PPROF2TextConfig{}.Convert(pprof, text)
}
