package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/felixge/pprofutils/utils"
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
		Name:       "version",
		ShortUsage: "pprofutils version",
		ShortHelp:  "Print version and exit.",
		Exec: func(_ context.Context, _ []string) error {
			os.Stdout.WriteString(version + "\n")
			return nil
		},
	})

	ffCommands = append(ffCommands, &ffcli.Command{
		Name:       "serve",
		FlagSet:    serveFlagSet,
		ShortUsage: "pprofutils serve [-addr addr]",
		ShortHelp:  "Serves pprofutils as a HTTP REST API.",
		Exec: func(_ context.Context, _ []string) error {
			log.Printf("Serving pprofutils %s via http at %s", version, *serveAddr)
			return http.ListenAndServe(*serveAddr, newHTTPServer())
		},
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

	if err := rootCmd.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	return nil
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
			a.Inputs = append(a.Inputs, in)
			a.Output = out
			a.Flags = make(map[string]interface{})
			for k, v := range flags {
				switch vt := v.(type) {
				case *bool:
					a.Flags[k] = *vt
				}
			}

			return cmd.Execute(ctx, a)
		},
	}
}

func run2() error {
	var (
		rootFlagSet  = flag.NewFlagSet("pprofutils", flag.ExitOnError)
		jsonFlagSet  = flag.NewFlagSet("pprofutils json", flag.ExitOnError)
		jsonSimple   = jsonFlagSet.Bool("simple", false, "Use simplified JSON format.")
		serveFlagSet = flag.NewFlagSet("pprofutils serve", flag.ExitOnError)
		serveAddr    = jsonFlagSet.String("addr", "localhost:8080", "HTTP listen addr.")
	)

	jsonCmd := &ffcli.Command{
		Name:       "json",
		FlagSet:    jsonFlagSet,
		ShortUsage: "pprofutils json [-simple] <input file> <output file>",
		ShortHelp:  "Converts from pprof to json and vice versa.",
		LongHelp: `The input and output file default to "-" which means stdin or stdout. If the` + "\n" +
			`input is pprof the output is json and for json inputs the output is pprof. This` + "\n" +
			`is automatically detected.` + "\n",
		Exec: func(ctx context.Context, args []string) error {
			in, out, err := openInputOutput(args)
			if err != nil {
				return err
			}
			defer in.Close()
			defer out.Close()

			return (&utils.JSON{
				Input:  in,
				Output: out,
				Simple: *jsonSimple,
			}).Execute(ctx)
		},
	}

	// TODO custom Command{} struct that can generate both http handlers
	// as well as cli

	serveCmd := &ffcli.Command{
		Name:       "serve",
		FlagSet:    serveFlagSet,
		ShortUsage: "pprofutils serve [-addr addr]",
		ShortHelp:  "Serves pprofutils as a HTTP REST API.",
		Exec: func(_ context.Context, _ []string) error {
			log.Printf("Serving pprofutils %s via http at %s", version, *serveAddr)
			return http.ListenAndServe(*serveAddr, newHTTPServer())
		},
	}

	versionCmd := &ffcli.Command{
		Name:       "version",
		FlagSet:    serveFlagSet,
		ShortUsage: "pprofutils version",
		ShortHelp:  "Print version and exit.",
		Exec: func(_ context.Context, _ []string) error {
			os.Stdout.WriteString(version + "\n")
			return nil
		},
	}

	var rootCmd *ffcli.Command
	rootCmd = &ffcli.Command{
		ShortUsage:  "pprofutils <subcommand>",
		FlagSet:     rootFlagSet,
		Subcommands: []*ffcli.Command{jsonCmd, serveCmd, versionCmd},
		Exec: func(_ context.Context, _ []string) error {
			os.Stdout.WriteString(rootCmd.UsageFunc(rootCmd))
			return nil
		},
	}

	if err := rootCmd.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	return nil
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
