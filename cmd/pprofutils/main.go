package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/felixge/pprofutils/v2/internal"
	"github.com/peterbourgon/ff/v3/ffcli"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

//go:generate bash -c "../../scripts/generate_version.bash > version.go"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() error {
	var addr = "localhost:8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	var (
		rootFlagSet  = flag.NewFlagSet("pprofutils", flag.ExitOnError)
		ffCommands   []*ffcli.Command
		serveFlagSet = flag.NewFlagSet("pprofutils serve", flag.ExitOnError)
		serveAddr    = serveFlagSet.String("addr", addr, "HTTP listen addr.")
		profiling    = serveFlagSet.Bool("profiling", false, "Enable profiling.")
		tracing      = serveFlagSet.Bool("tracing", false, "Enable tracing.")
	)

	for _, util := range internal.Utils {
		ffCommands = append(ffCommands, ffCommand(util))
	}

	ffCommands = append(ffCommands, &ffcli.Command{
		Name:       "serve",
		FlagSet:    serveFlagSet,
		ShortUsage: "pprofutils serve [-addr addr]",
		ShortHelp:  "Serves pprofutils as a HTTP REST API",
		Exec: func(_ context.Context, _ []string) error {
			if *profiling {
				log.Printf("Starting datadog profiler")
				profilerOptions := []profiler.Option{
					profiler.WithVersion(version),
					profiler.CPUDuration(60 * time.Second),
					profiler.WithPeriod(60 * time.Second),
					profiler.WithProfileTypes(
						profiler.CPUProfile,
						profiler.HeapProfile,
						profiler.BlockProfile,
						profiler.MutexProfile,
						profiler.GoroutineProfile,
					),
				}
				if err := profiler.Start(profilerOptions...); err != nil {
					return err
				}
				defer profiler.Stop()
			}

			if *tracing {
				log.Printf("Starting datadog tracer")
				tracer.Start(
					tracer.WithServiceVersion(version),
				)
				defer tracer.Stop()
			}

			log.Printf("Serving pprofutils %s via http at %s", version, *serveAddr)
			return http.ListenAndServe(*serveAddr, newHTTPServer())
		},
	})

	ffCommands = append(ffCommands, &ffcli.Command{
		Name:       "version",
		ShortUsage: "pprofutils version",
		ShortHelp:  "Print version and exit",
		Exec: func(_ context.Context, _ []string) error {
			os.Stdout.WriteString(version + "\n")
			return nil
		},
	})

	sort.Slice(ffCommands, func(i, j int) bool {
		return ffCommands[i].Name < ffCommands[j].Name
	})

	var rootCmd *ffcli.Command
	rootCmd = &ffcli.Command{
		ShortUsage:  "pprofutils <subcommand>",
		FlagSet:     rootFlagSet,
		Subcommands: ffCommands,
		Exec: func(_ context.Context, _ []string) error {
			os.Stdout.WriteString(rootCmd.UsageFunc(rootCmd))
			return nil
		},
	}

	return rootCmd.ParseAndRun(context.Background(), os.Args[1:])
}

func ffCommand(util internal.Util) *ffcli.Command {
	fs := flag.NewFlagSet("pprofutils "+util.Name, flag.ExitOnError)
	flags := map[string]interface{}{}
	for name, bf := range util.Flags {
		val := bf.Default
		switch vt := val.(type) {
		case time.Duration:
			fs.DurationVar(&vt, name, vt, bf.Usage)
			flags[name] = &vt
		case bool:
			fs.BoolVar(&vt, name, vt, bf.Usage)
			flags[name] = &vt
		case string:
			fs.StringVar(&vt, name, vt, bf.Usage)
			flags[name] = &vt
		}
	}

	return &ffcli.Command{
		Name:       util.Name,
		ShortUsage: fmt.Sprintf("pprofutils %s %s", util.Name, util.ShortUsage),
		ShortHelp:  util.ShortHelp,
		LongHelp:   util.LongHelp,
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			in, out, err := openInputOutput(args)
			if err != nil {
				return err
			}
			defer in.Close()
			defer out.Close()

			a := &internal.UtilArgs{}
			inBuf, err := ioutil.ReadAll(in)
			if err != nil {
				return err
			}
			a.Inputs = append(a.Inputs, inBuf)
			a.Output = out
			a.Flags = make(map[string]interface{})
			for k, v := range flags {
				switch vt := v.(type) {
				case *time.Duration:
					a.Flags[k] = *vt
				case *bool:
					a.Flags[k] = *vt
				case *string:
					a.Flags[k] = *vt
				}
			}

			return util.Execute(ctx, a)
		},
	}
}

func openInputOutput(args []string) (io.ReadCloser, io.WriteCloser, error) {
	inputPath := "-"
	if len(args) >= 1 {
		inputPath = args[0]
	}
	in, err := openInput(inputPath)
	if err != nil {
		return nil, nil, err
	}

	outputPath := "-"
	if len(args) >= 2 {
		outputPath = args[1]
	}
	out, err := openOutput(outputPath)
	if err != nil {
		return nil, nil, err
	}
	return in, out, nil
}

func openInput(path string) (io.ReadCloser, error) {
	if path == "-" {
		return io.NopCloser(os.Stdin), nil
	}
	return os.Open(path)
}

func openOutput(path string) (io.WriteCloser, error) {
	if path == "-" {
		return nopWriteCloser{os.Stdout}, nil
	}
	return os.Create(path)
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }
