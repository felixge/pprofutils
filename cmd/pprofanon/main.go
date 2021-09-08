package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/felixge/pprofutils/internal"
	"github.com/google/pprof/profile"
	"github.com/wolfeidau/humanhash"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var (
		versionF   = flag.Bool("version", false, "Print version and exit.")
		whitelistF = flag.String("w", `^runtime;^net;^encoding`, "Semicolon separated whitelist of func.pkg regexes to allow.")
	)
	flag.Parse()
	if *versionF {
		fmt.Printf("%s\n", internal.Version)
		return nil
	}
	inProf, err := profile.Parse(os.Stdin)
	if err != nil {
		return err
	}

	var whitelisted []*regexp.Regexp
	for _, w := range strings.Split(*whitelistF, ";") {
		r, err := regexp.Compile(w)
		if err != nil {
			return err
		}
		whitelisted = append(whitelisted, r)
	}

outer:
	for _, f := range inProf.Function {
		for _, w := range whitelisted {
			if w.MatchString(f.Name) {
				continue outer
			}
		}

		fmt.Fprintf(os.Stderr, "%s\n", f.Name)
		h := sha1.Sum([]byte(f.Name))
		f.Name, err = humanhash.Humanize([]byte(h[:]), 3)
		if err != nil {
			return err
		}
		f.SystemName = ""

		parts := strings.Split(f.Filename, string(filepath.Separator))
		for i, p := range parts {
			if p != "" {
				h := sha1.Sum([]byte(p))
				parts[i], err = humanhash.Humanize(h[:], 1)
				if err != nil {
					return err
				}
			}
		}
		f.Filename = strings.Join(parts, string(filepath.Separator))
	}

	return inProf.Write(os.Stdout)
}
