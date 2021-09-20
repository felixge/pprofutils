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

	"github.com/peterbourgon/ff/v3/ffcli"
)

var version = "vN/A"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() error {
	var (
		rootFlagSet  = flag.NewFlagSet("pprofutils", flag.ExitOnError)
		ffCommands   []*ffcli.Command
		serveFlagSet = flag.NewFlagSet("pprofutils serve", flag.ExitOnError)
		serveAddr    = serveFlagSet.String("addr", "localhost:8080", "HTTP listen addr.")
	)

	for _, util := range utilCommands {
		ffCommands = append(ffCommands, ffCommand(util))
	}

	ffCommands = append(ffCommands, &ffcli.Command{
		Name:       "serve",
		FlagSet:    serveFlagSet,
		ShortUsage: "pprofutils serve [-addr addr]",
		ShortHelp:  "Serves pprofutils as a HTTP REST API",
		Exec: func(_ context.Context, _ []string) error {
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

func ffCommand(cmd UtilCommand) *ffcli.Command {
	fs := flag.NewFlagSet("pprofutils "+cmd.Name, flag.ExitOnError)
	flags := map[string]interface{}{}
	for name, bf := range cmd.Flags {
		val := bf.Default
		switch vt := val.(type) {
		case bool:
			fs.BoolVar(&vt, name, vt, bf.Usage)
			flags[name] = &vt
		case string:
			fs.StringVar(&vt, name, vt, bf.Usage)
			flags[name] = &vt
		}
	}

	return &ffcli.Command{
		Name:       cmd.Name,
		ShortUsage: fmt.Sprintf("pprofutils %s %s", cmd.Name, cmd.ShortUsage),
		ShortHelp:  cmd.ShortHelp,
		LongHelp:   cmd.LongHelp,
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			in, out, err := openInputOutput(args)
			if err != nil {
				return err
			}
			defer in.Close()
			defer out.Close()

			a := &UtilArgs{}
			inBuf, err := ioutil.ReadAll(in)
			if err != nil {
				return err
			}
			a.Inputs = append(a.Inputs, inBuf)
			a.Output = out
			a.Flags = make(map[string]interface{})
			for k, v := range flags {
				switch vt := v.(type) {
				case *bool:
					a.Flags[k] = *vt
				case *string:
					a.Flags[k] = *vt
				}
			}

			return cmd.Execute(ctx, a)
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
