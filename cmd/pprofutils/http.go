package main

import (
	"bytes"
	"errors"
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
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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
	handler := func(w http.ResponseWriter, r *http.Request) error {
		var in io.Reader
		out := &bytes.Buffer{}
		a := &UtilArgs{Output: out}
		contentType := r.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(maxPostSize); err != nil {
				return fmt.Errorf("upload too big: %w", err)
			}

			var first *multipart.FileHeader
			for _, files := range r.MultipartForm.File {
				for _, file := range files {
					if first != nil {
						return errors.New("only one file is expected to be uploaded")
					}
					first = file
				}
			}

			file, err := first.Open()
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()

			in = file
		} else {
			in = io.LimitReader(r.Body, maxPostSize)
		}

		inBuf, err := ioutil.ReadAll(in)
		if err != nil {
			return fmt.Errorf("upload error: %w", err)
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
					return fmt.Errorf("bad query param %s: %w", name, err)
				}
				a.Flags[name] = val
			case string:
				a.Flags[name] = qVal
			}
		}

		if err := cmd.Execute(r.Context(), a); err != nil {
			return err
		}
		_, _ = io.Copy(w, out)
		return nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, _ := tracer.SpanFromContext(r.Context())
		span.SetTag("http.full_url", r.URL.String())
		span.SetTag("http.content_length", r.Header.Get("Content-Length"))
		span.SetTag("user.ip", r.Header.Get("Fly-Client-IP"))
		span.SetTag("user.agent", r.Header.Get("User-Agent"))
		if err := handler(w, r); err != nil {
			http.Error(w, fmt.Sprintf("error: %s\n", err), http.StatusBadRequest)
			span.Finish(tracer.WithError(err))
		}
	})
}
