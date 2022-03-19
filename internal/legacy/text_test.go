package legacy

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestTextConvert(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		is := is.New(t)
		textIn := strings.TrimSpace(`
main;foo 5
main;foobar 4
main;foo;bar 3
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
main;foobar 4
main;foo;bar 3
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
main;foobar 4 40000000
main;foo;bar 3 30000000
	`)
		proto, err := Text{}.Convert(strings.NewReader(textIn))
		is.NoErr(err)
		textOut := bytes.Buffer{}
		is.NoErr(Protobuf{SampleTypes: true}.Convert(proto, &textOut))
		is.Equal(textIn+"\n", textOut.String())
	})
}

func TestSplitLastN(t *testing.T) {
	for _, tc := range []struct {
		line        string
		n           int
		expected    []string
		expectedErr bool
	}{

		{"main;foo 5", 1, []string{"main;foo", "5"}, false},
		{"main;foo 5", 2, nil, true},
		{"main;foo 5 10", 2, []string{"main;foo", "5", "10"}, false},
		{"main;foo;foo bar 5", 1, []string{"main;foo;foo bar", "5"}, false},
		{"main;foo;foo bar 5 10", 2, []string{"main;foo;foo bar", "5", "10"}, false},
	} {

		t.Run(fmt.Sprintf("line=%q, n=%d", tc.line, tc.n), func(t *testing.T) {
			res, err := splitLastN(tc.line, tc.n)
			if tc.expectedErr && err == nil {
				t.Fatalf("line=%q, n=%d, expected err but passed", tc.line, tc.n)
			}
			if !tc.expectedErr && err != nil {
				t.Fatalf("line=%q, n=%d, unexpected err: %v", tc.line, tc.n, err)
			}
			if !reflect.DeepEqual(tc.expected, res) {
				t.Fatalf("expected=%v, got=%v", tc.expected, res)
			}
		})
	}
}
