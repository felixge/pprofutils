package legacy

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/pprof/profile"
	"github.com/matryer/is"
)

func TestProtobufConvert(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		is := is.New(t)
		data, err := ioutil.ReadFile(filepath.Join("test-fixtures", "pprof.samples.cpu.001.pb.gz"))
		is.NoErr(err)

		proto, err := profile.Parse(bytes.NewReader(data))
		is.NoErr(err)

		out := bytes.Buffer{}
		is.NoErr(Protobuf{}.Convert(proto, &out))
		want := strings.TrimSpace(`
golang.org/x/sync/errgroup.(*Group).Go.func1;main.run.func2;main.computeSum 19
runtime.mcall;runtime.park_m;runtime.resetForSleep;runtime.resettimer;runtime.modtimer;runtime.wakeNetPoller;runtime.netpollBreak;runtime.write;runtime.write1 7
golang.org/x/sync/errgroup.(*Group).Go.func1;main.run.func2;main.computeSum;runtime.asyncPreempt 5
runtime.mstart;runtime.mstart1;runtime.sysmon;runtime.usleep 3
runtime.mcall;runtime.park_m;runtime.schedule;runtime.findrunnable;runtime.stopm;runtime.notesleep;runtime.semasleep;runtime.pthread_cond_wait 2
runtime.mcall;runtime.gopreempt_m;runtime.goschedImpl;runtime.schedule;runtime.findrunnable;runtime.stopm;runtime.notesleep;runtime.semasleep;runtime.pthread_cond_wait 1
runtime.mcall;runtime.park_m;runtime.schedule;runtime.findrunnable;runtime.checkTimers;runtime.nanotime;runtime.nanotime1 1
`) + "\n"
		is.Equal(out.String(), want)
	})

	t.Run("different lines per function", func(t *testing.T) {
		is := is.New(t)
		data, err := ioutil.ReadFile(filepath.Join("test-fixtures", "pprof.lines.pb.gz"))
		is.NoErr(err)

		proto, err := profile.Parse(bytes.NewReader(data))
		is.NoErr(err)

		out := bytes.Buffer{}
		is.NoErr(Protobuf{}.Convert(proto, &out))
		want := strings.TrimSpace(`
main.run.func1;main.threadKind.Run;main.goGo1;main.goHog 85
main.run.func1;main.threadKind.Run;main.goGo2;main.goHog 78
main.run.func1;main.threadKind.Run;main.goGo3;main.goHog 72
main.run.func1;main.threadKind.Run;main.goGo0;main.goHog 72
main.run.func1;main.threadKind.Run;main.goGo0;main.goHog;runtime.asyncPreempt 1
`) + "\n"
		is.Equal(out.String(), want)
	})
}
