package utils

import (
	"bytes"
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/pprof/profile"
	"github.com/stretchr/testify/require"
)

func TestHeapage(t *testing.T) {
	data, err := ioutil.ReadFile(filepath.Join("testdata", "heapage.pprof"))
	require.NoError(t, err)

	buf := &bytes.Buffer{}
	h := &Heapage{Input: data, Output: buf, Period: 10 * time.Second}
	require.NoError(t, h.Execute(context.Background()))

	_, err = profile.Parse(buf)
	require.NoError(t, err)

	// TODO: finish testing :)
}
