package pprofutils

import (
	"bufio"
	"fmt"
	"io"

	"github.com/google/pprof/profile"
)

func PPROF2Text(pprof io.Reader, text io.Writer) error {
	prof, err := profile.Parse(pprof)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(text)
	for _, sample := range prof.Sample {
		for i, loc := range sample.Location {
			frame := loc.Line[0].Function.Name
			w.WriteString(frame)
			if i+1 < len(sample.Location) {
				w.WriteString(";")
			}
		}
		fmt.Fprintf(w, " %d\n", sample.Value[0])
	}
	return w.Flush()
}
