package pprofutils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestText2PPROF(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		is := is.New(t)
		textIn := strings.TrimSpace(`
main;foo 5
main;foo;bar 3
main;foobar 4
`)
		pprofOut := bytes.Buffer{}
		is.NoErr(Text2PPROF(strings.NewReader(textIn), &pprofOut))
		textOut := bytes.Buffer{}
		is.NoErr(PPROF2Text(&pprofOut, &textOut))
		is.Equal(textIn+"\n", textOut.String())
	})

	t.Run("multiple sample types", func(t *testing.T) {
		is := is.New(t)
		textIn := strings.TrimSpace(`
samples/count duration/nanoseconds
main;foo 5 50000000
main;foo;bar 3 30000000
main;foobar 4 40000000
	`)
		pprofOut := bytes.Buffer{}
		is.NoErr(Text2PPROF(strings.NewReader(textIn), &pprofOut))
		textOut := bytes.Buffer{}
		is.NoErr(PPROF2TextConfig{SampleTypes: true}.Convert(&pprofOut, &textOut))
		is.Equal(textIn+"\n", textOut.String())
	})
}
