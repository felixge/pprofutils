package pprofutils

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestPPROF2Text(t *testing.T) {
	is := is.New(t)
	data, err := os.ReadFile(filepath.Join("test-fixtures", "pprof.samples.cpu.001.pb.gz"))
	is.NoErr(err)

	out := bytes.Buffer{}
	is.NoErr(PPROF2Text(bytes.NewReader(data), &out))
	want := strings.TrimSpace(`
main.computeSum;main.run.func2;golang.org/x/sync/errgroup.(*Group).Go.func1 19
runtime.asyncPreempt;main.computeSum;main.run.func2;golang.org/x/sync/errgroup.(*Group).Go.func1 5
runtime.pthread_cond_wait;runtime.semasleep;runtime.notesleep;runtime.stopm;runtime.findrunnable;runtime.schedule;runtime.goschedImpl;runtime.gopreempt_m;runtime.mcall 1
runtime.nanotime1;runtime.nanotime;runtime.findrunnable;runtime.schedule;runtime.park_m;runtime.mcall 1
runtime.pthread_cond_wait;runtime.semasleep;runtime.notesleep;runtime.stopm;runtime.findrunnable;runtime.schedule;runtime.park_m;runtime.mcall 2
runtime.write1;runtime.write;runtime.wakeNetPoller;runtime.modtimer;runtime.resettimer;runtime.park_m;runtime.mcall 7
runtime.usleep;runtime.sysmon;runtime.mstart1;runtime.mstart 3
`) + "\n"
	is.Equal(out.String(), want)
}
