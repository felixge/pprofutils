package pprofutils

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/google/pprof/profile"
)

func PPROF2Text(pprof io.Reader, text io.Writer) error {
	prof, err := profile.Parse(pprof)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(text)
	for _, sample := range prof.Sample {
		var frames []string
		for i := range sample.Location {
			loc := sample.Location[len(sample.Location)-i-1]
			for j := range loc.Line {
				line := loc.Line[len(loc.Line)-j-1]
				frames = append(frames, line.Function.Name)
			}
		}
		fmt.Fprintf(w, "%s %d\n", strings.Join(frames, ";"), sample.Value[0])
	}
	return w.Flush()
}
