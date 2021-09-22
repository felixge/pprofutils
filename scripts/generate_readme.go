//build:ignore

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/felixge/pprofutils/internal"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	input := &bytes.Buffer{}
	_, err := io.Copy(input, os.Stdin)
	if err != nil {
		return err
	}

	var fns = template.FuncMap{
		"queryflags": queryflags,
		"defaultval": defaultval,
		"example":    example,
	}

	tmpl, err := template.New("README.md").Funcs(fns).Parse(input.String())
	if err != nil {
		return err
	}

	return tmpl.Execute(os.Stdout, internal.Utils)
}

func defaultval(val interface{}) string {
	defaultVal := fmt.Sprintf("%v", val)
	if defaultVal == "" {
		defaultVal = "..."
	}
	return defaultVal
}

func queryflags(flags map[string]internal.UtilFlag) string {
	if len(flags) == 0 {
		return ""
	}

	var params []string
	for name, f := range flags {
		params = append(params, name+"="+defaultval(f.Default))
	}
	sort.Strings(params)
	return "?" + strings.Join(params, "&")
}

func example(util internal.Util) (string, error) {
	var (
		b        = &strings.Builder{}
		prefix   = filepath.Join("examples", util.Name)
		inPprof  = prefix + ".in.pprof"
		inPng    = prefix + ".in.png"
		inJson   = prefix + ".in.json"
		outPprof = prefix + ".out.pprof"
		outPng   = prefix + ".out.png"
		outJSON  = prefix + ".out.json"
	)

	if exists(inPprof) && exists(outJSON) {
		b.WriteString("```shell\n")
		b.WriteString(fmt.Sprintf("pprofutils %s %s %s\n", util.Name, inPprof, outJSON))
		b.WriteString("# or\n")
		b.WriteString(fmt.Sprintf("curl --data-binary @%s pprof.to/%s > %s\n", inPprof, util.Name, outJSON))
		b.WriteString("```\n")
		b.WriteString(fmt.Sprintf("Converts [%s](./%s) from pprof to json, see [%s](./%s).\n", inPprof, inPprof, outJSON, outJSON))
	}

	if exists(inJson) && exists(outPprof) {
		b.WriteString("```shell\n")
		b.WriteString(fmt.Sprintf("pprofutils %s %s %s\n", util.Name, inJson, outPprof))
		b.WriteString("# or\n")
		b.WriteString(fmt.Sprintf("curl --data-binary @%s pprof.to/%s > %s\n", inJson, util.Name, outPprof))
		b.WriteString("```\n")
		b.WriteString(fmt.Sprintf("Converts [%s](./%s) from json to pprof, see [%s](./%s).\n", inJson, inJson, outPprof, outPprof))
	}

	if exists(inPng) && exists(outPng) {
		fmt.Fprintf(
			b,
			"```\npprofutils %s %s %s\n```\n",
			util.Name,
			inPprof,
			outPprof,
		)
		fmt.Fprintf(b, "Converts a profile that looks like this:\n\n![](./%s)\n\n", inPng)
		fmt.Fprintf(b, "Into a new profile that looks like that:\n\n![](./%s)\n", outPng)
	}

	//if exists(inputPng) {
	//return fmt.Sprintf("![](./%s)", inputPng)
	//}

	//fmt.Fprintf(os.Stderr, "%s\n", util.Name)
	return b.String(), nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
