package utils

import (
	"context"
	"crypto/sha1"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/pprof/profile"
	"github.com/wolfeidau/humanhash"
)

type Anon struct {
	Input     []byte
	Output    io.Writer
	Whitelist string
}

func (a *Anon) Execute(ctx context.Context) error {
	prof, err := profile.ParseData(a.Input)
	if err != nil {
		return err
	}

	var whitelisted []*regexp.Regexp
	wl := strings.TrimSpace(a.Whitelist)
	if len(wl) > 0 {
		for _, w := range strings.Split(wl, ";") {
			r, err := regexp.Compile(w)
			if err != nil {
				return err
			}
			whitelisted = append(whitelisted, r)
		}
	}

outer:
	for _, f := range prof.Function {
		for _, w := range whitelisted {
			if w.MatchString(f.Name) {
				continue outer
			}
		}

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

	return prof.Write(a.Output)
}
