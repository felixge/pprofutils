//build:ignore

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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
		"examples":   examples,
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

func examples(util internal.Util) (string, error) {
	var (
		b = &strings.Builder{}
	)

	for i, e := range util.Examples {
		pathTo := func(dir, format string) string {
			return filepath.Join("examples", util.Name+"."+dir+"."+format)
		}
		combination := func(in, out []string) bool {
			return stringsEquals(e.In, in) && stringsEquals(e.Out, out)
		}

		b.WriteString(fmt.Sprintf("#### Example %d: %s\n", i+1, e.Name))

		if len(e.In) == 1 && len(e.Out) == 1 {
			b.WriteString(simpleInOut(util.Name, pathTo("in", e.In[0]), pathTo("out", e.Out[0])))
		} else if combination([]string{"txt"}, []string{"pprof", "png"}) {
			inTxt, outPprof, outPng := pathTo("in", "txt"), pathTo("out", "pprof"), pathTo("out", "png")
			inData, err := ioutil.ReadFile(inTxt)
			if err != nil {
				return "", err
			}

			b.WriteString(shell(util.Name, inTxt, outPprof))
			b.WriteString(fmt.Sprintf("Converts [%s](./%s) with the following content:\n\n", inTxt, inTxt))
			b.WriteString(fmt.Sprintf("```\n%s\n```\n\n", inData))
			b.WriteString(outPprofImage(outPprof, outPng))
		} else if combination([]string{"pprof", "png"}, []string{"txt"}) {
			inPprof, inPng, outTxt := pathTo("in", "pprof"), pathTo("in", "png"), pathTo("out", "txt")
			outData, err := ioutil.ReadFile(outTxt)
			if err != nil {
				return "", err
			}

			b.WriteString(shell(util.Name, inPprof, outTxt))
			b.WriteString(inPprofImage(inPprof, inPng))
			b.WriteString(fmt.Sprintf("Into a new folded text file [%s](./%s) that looks like this:\n\n", outTxt, outTxt))
			b.WriteString(fmt.Sprintf("```\n%s\n```\n\n", outData))
		} else if combination([]string{"pprof", "png"}, []string{"pprof", "png"}) {
			inPprof, outPprof := pathTo("in", "pprof"), pathTo("out", "pprof")
			inPng, outPng := pathTo("in", "png"), pathTo("out", "png")
			b.WriteString(shell(util.Name, inPprof, outPprof))
			b.WriteString(inPprofImage(inPprof, inPng))
			b.WriteString(outPprofImage(outPprof, outPng))
		}
	}
	return b.String(), nil
}

func simpleInOut(util, in, out string) string {
	b := &strings.Builder{}
	b.WriteString(shell(util, in, out))
	b.WriteString(fmt.Sprintf("See [%s](./%s) and [%s](./%s) for more details.\n", in, in, out, out))
	return b.String()
}

func shell(util, in, out string) string {
	b := &strings.Builder{}
	b.WriteString("```shell\n")
	b.WriteString(fmt.Sprintf("pprofutils %s %s %s\n", util, in, out))
	b.WriteString("# or\n")
	b.WriteString(fmt.Sprintf("curl --data-binary @%s pprof.to/%s > %s\n", in, util, out))
	b.WriteString("```\n")
	return b.String()
}

func inPprofImage(inPprof, inPng string) string {
	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf("Converts the profile [%s](./%s) that looks like this:\n\n", inPprof, inPprof))
	b.WriteString(fmt.Sprintf("![](%s)\n\n", inPng))
	return b.String()
}

func outPprofImage(outPprof, outPng string) string {
	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf("Into a new profile [%s](./%s) that looks like this:\n\n", outPprof, outPprof))
	b.WriteString(fmt.Sprintf("![](%s)\n", outPng))
	return b.String()
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func stringsEquals(a, b []string) bool {
	aCopy := make([]string, len(a))
	copy(aCopy, a)
	bCopy := make([]string, len(b))
	copy(bCopy, b)
	sort.Strings(aCopy)
	sort.Strings(bCopy)
	return fmt.Sprint(aCopy) == fmt.Sprint(bCopy)
}
