package pprofutils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestDelta(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		var (
			is        = is.New(t)
			profA     bytes.Buffer
			profB     bytes.Buffer
			delta     bytes.Buffer
			deltaText bytes.Buffer
		)

		is.NoErr(Text2PPROF(strings.NewReader(strings.TrimSpace(`
main;foo 5
main;foo;bar 3
main;foobar 4
`)), &profA))

		is.NoErr(Text2PPROF(strings.NewReader(strings.TrimSpace(`
main;foo 8
main;foo;bar 3
main;foobar 5
`)), &profB))
		is.NoErr(Delta(&profA, &profB, &delta))

		is.NoErr(PPROF2Text(&delta, &deltaText))
		is.Equal(deltaText.String(), strings.TrimSpace(`
main;foo 3
main;foobar 1
`)+"\n")
	})

	t.Run("sample types", func(t *testing.T) {
		var (
			is        = is.New(t)
			profA     bytes.Buffer
			profB     bytes.Buffer
			delta     bytes.Buffer
			deltaText bytes.Buffer
		)

		is.NoErr(Text2PPROF(strings.NewReader(strings.TrimSpace(`
x/count y/count
main;foo 5 10
main;foo;bar 3 6
main;foo;baz 9 0
main;foobar 4 8
`)), &profA))

		is.NoErr(Text2PPROF(strings.NewReader(strings.TrimSpace(`
x/count y/count
main;foo 8 16
main;foo;bar 3 6
main;foo;baz 9 0
main;foobar 5 10
`)), &profB))
		is.NoErr(DeltaConfig{SampleTypes: []SampleType{{Type: "x", Unit: "count"}}}.Convert(&profA, &profB, &delta))

		is.NoErr(PPROF2TextConfig{SampleTypes: true}.Convert(&delta, &deltaText))
		is.Equal(deltaText.String(), strings.TrimSpace(`
x/count y/count
main;foo 3 16
main;foo;bar 0 6
main;foobar 1 10
`)+"\n")
	})
}
