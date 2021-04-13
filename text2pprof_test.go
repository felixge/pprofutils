package pprofutils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestText2PPROF(t *testing.T) {
	is := is.New(t)
	textIn := strings.TrimSpace(`main;foo 5
main;foo;bar 3
main;foobar 4
`)
	pprofOut := bytes.Buffer{}
	is.NoErr(Text2PPROF(strings.NewReader(textIn), &pprofOut))
	textOut := bytes.Buffer{}
	is.NoErr(PPROF2Text(&pprofOut, &textOut))
	is.Equal(textIn+"\n", textOut.String())
}
