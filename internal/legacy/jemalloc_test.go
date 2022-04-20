package legacy

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/matryer/is"
)

func TestJemallocConvert(t *testing.T) {
	is := is.New(t)
	f, err := os.Open("test-fixtures/jemalloc.heap")
	is.NoErr(err)
	proto, err := Jemalloc{}.Convert(f)
	is.NoErr(err)

	gold, err := ioutil.ReadFile("test-fixtures/jemalloc.heap.string")
	is.NoErr(err)

	js := proto.String()

	if js != string(gold) {
		comparePprof(t, []byte(js), gold)
	}
}

func comparePprof(t *testing.T, in, expected []byte) {
	f1, err := ioutil.TempFile("", "proto_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", "proto_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	f1.Write(in)
	f2.Write(expected)

	data, err := exec.Command("diff", "-u", f1.Name(), f2.Name()).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Errorf("diff: %s\n", data)
}
