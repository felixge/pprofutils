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
golang.org/x/sync/errgroup.(*Group).Go.func1;main.run.func2;main.computeSum 19
golang.org/x/sync/errgroup.(*Group).Go.func1;main.run.func2;main.computeSum;runtime.asyncPreempt 5
runtime.mcall;runtime.gopreempt_m;runtime.goschedImpl;runtime.schedule;runtime.findrunnable;runtime.stopm;runtime.notesleep;runtime.semasleep;runtime.pthread_cond_wait 1
runtime.mcall;runtime.park_m;runtime.schedule;runtime.findrunnable;runtime.checkTimers;runtime.nanotime;runtime.nanotime1 1
runtime.mcall;runtime.park_m;runtime.schedule;runtime.findrunnable;runtime.stopm;runtime.notesleep;runtime.semasleep;runtime.pthread_cond_wait 2
runtime.mcall;runtime.park_m;runtime.resetForSleep;runtime.resettimer;runtime.modtimer;runtime.wakeNetPoller;runtime.netpollBreak;runtime.write;runtime.write1 7
runtime.mstart;runtime.mstart1;runtime.sysmon;runtime.usleep 3
`) + "\n"
	is.Equal(out.String(), want)
}
