package pprofutils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestTextConvert(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		is := is.New(t)
		textIn := strings.TrimSpace(`
main;foo 5
main;foo;bar 3
main;foobar 4
`)
		proto, err := Text{}.Convert(strings.NewReader(textIn))
		is.NoErr(err)
		textOut := bytes.Buffer{}
		is.NoErr(Protobuf{}.Convert(proto, &textOut))
		is.Equal(textIn+"\n", textOut.String())
	})

	t.Run("header with one sample type", func(t *testing.T) {
		is := is.New(t)
		textIn := strings.TrimSpace(`
samples/count
main;foo 5
main;foo;bar 3
main;foobar 4
	`)
		proto, err := Text{}.Convert(strings.NewReader(textIn))
		is.NoErr(err)
		textOut := bytes.Buffer{}
		is.NoErr(Protobuf{SampleTypes: true}.Convert(proto, &textOut))
		is.Equal(textIn+"\n", textOut.String())
	})

	t.Run("header with multiple sample types", func(t *testing.T) {
		is := is.New(t)
		textIn := strings.TrimSpace(`
samples/count duration/nanoseconds
main;foo 5 50000000
main;foo;bar 3 30000000
main;foobar 4 40000000
	`)
		proto, err := Text{}.Convert(strings.NewReader(textIn))
		is.NoErr(err)
		textOut := bytes.Buffer{}
		is.NoErr(Protobuf{SampleTypes: true}.Convert(proto, &textOut))
		is.Equal(textIn+"\n", textOut.String())
	})
}
