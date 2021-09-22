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
	router.HandlerFunc("GET", "/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Location", "https://github.com/felixge/pprofutils")
		w.WriteHeader(http.StatusFound)
	})

	for _, cmd := range utilCommands {
		router.Handler("POST", "/"+cmd.Name, utilHandler(cmd))
	}

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = addSpanTags(r)
		http.NotFoundHandler().ServeHTTP(w, r)
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Service-Version", version)
		m := httpsnoop.CaptureMetrics(router, w, r)
		log.Printf("%d %s %s %s", m.Code, r.Method, r.URL, m.Duration)
	})
}

func utilHandler(cmd UtilCommand) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		span := addSpanTags(r)
		defer func() {
			span.Finish(tracer.WithError(err))
		}()

		out := &bytes.Buffer{}
		a := &UtilArgs{Output: out}

		upload := func() error {
			var in io.Reader
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
			return nil
		}

		uploadSpan, _ := tracer.StartSpanFromContext(r.Context(), "upload")
		err = upload()
		uploadSpan.Finish(tracer.WithError(err))
		if err != nil {
			http.Error(w, fmt.Sprintf("error: %s\n", err), http.StatusBadRequest)
			return
		}

		execSpan, execCtx := tracer.StartSpanFromContext(r.Context(), "exec")
		err = cmd.Execute(execCtx, a)
		execSpan.Finish(tracer.WithError(err))
		if err != nil {
			http.Error(w, fmt.Sprintf("error: %s\n", err), http.StatusBadRequest)
			return
		}

		respondSpan, _ := tracer.StartSpanFromContext(r.Context(), "respond")
		_, err = io.Copy(w, out)
		respondSpan.Finish(tracer.WithError(err))
	})
}

func addSpanTags(r *http.Request) tracer.Span {
	span, _ := tracer.SpanFromContext(r.Context())
	span.SetTag("http.full_url", r.URL.String())
	span.SetTag("http.content_length", r.Header.Get("Content-Length"))
	span.SetTag("user.ip", r.Header.Get("Fly-Client-IP"))
	span.SetTag("user.agent", r.Header.Get("User-Agent"))
	return span
}
