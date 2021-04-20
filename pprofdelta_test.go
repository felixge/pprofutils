package pprofutils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestDelta(t *testing.T) {
	is := is.New(t)

	var (
		profA     bytes.Buffer
		profB     bytes.Buffer
		delta     bytes.Buffer
		deltaText bytes.Buffer
	)

	is.NoErr(Text2PPROF(strings.NewReader(strings.TrimSpace(`main;foo 5
main;foo;bar 3
main;foobar 4
`)), &profA))

	is.NoErr(Text2PPROF(strings.NewReader(strings.TrimSpace(`main;foo 8
main;foo;bar 3
main;foobar 5
`)), &profB))
	is.NoErr(Delta(&profA, &profB, &delta))

	is.NoErr(PPROF2Text(&delta, &deltaText))
	is.Equal(deltaText.String(), strings.TrimSpace(`
main;foo 3
main;foobar 1
`)+"\n")
}
