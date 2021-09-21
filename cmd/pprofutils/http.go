package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/felixge/httpsnoop"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/julienschmidt/httprouter"
)

const maxPostSize = 128 * 1024 * 1024

func newHTTPServer() http.Handler {
	router := httptrace.New()
	for _, cmd := range utilCommands {
		router.Handler("POST", "/"+cmd.Name, utilHandler(cmd))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(router, w, r)
		log.Printf("%d %s %s %s", m.Code, r.Method, r.URL, m.Duration)
	})
}

func utilHandler(cmd UtilCommand) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var in io.Reader
		out := &bytes.Buffer{}
		a := &UtilArgs{Output: out}
		contentType := r.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(maxPostSize); err != nil {
				http.Error(w, "upload too big: "+err.Error(), http.StatusBadRequest)
				return
			}

			var first *multipart.FileHeader
			for _, files := range r.MultipartForm.File {
				for _, file := range files {
					if first != nil {
						http.Error(w, "only one file is expected to be uploaded\n", http.StatusBadRequest)
						return
					}
					first = file
				}
			}

			file, err := first.Open()
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to open file: %s\n", err), http.StatusBadRequest)
				return
			}
			defer file.Close()

			in = file
		} else {
			in = io.LimitReader(r.Body, maxPostSize)
		}

		inBuf, err := ioutil.ReadAll(in)
		if err != nil {
			http.Error(w, fmt.Sprintf("upload error: %s\n", err), http.StatusBadRequest)
			return
		}
		a.Inputs = append(a.Inputs, inBuf)

		a.Flags = make(map[string]interface{})
		for name, flag := range cmd.Flags {
			a.Flags[name] = flag.Default
			if _, ok := r.URL.Query()[name]; !ok {
				continue
			}

			qVal := r.URL.Query().Get(name)
			switch flag.Default.(type) {
			case bool:
				val, err := strconv.ParseBool(qVal)
				if err != nil {
					http.Error(w, fmt.Sprintf("bad query param %q: %s", name, err), http.StatusBadRequest)
					return
				}
				a.Flags[name] = val
			}
		}

		if err := cmd.Execute(r.Context(), a); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, _ = io.Copy(w, out)
	})
}
