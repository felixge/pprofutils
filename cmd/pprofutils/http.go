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
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/felixge/pprofutils/v2/internal"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/julienschmidt/httprouter"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const maxPostSize = 128 * 1024 * 1024

func newHTTPServer() http.Handler {
	router := httptrace.New()
	router.HandlerFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		addSpanTags(r)
		w.Header().Set("Location", "https://github.com/felixge/pprofutils#readme")
		w.WriteHeader(http.StatusFound)
	})

	for _, util := range internal.Utils {
		router.Handler("POST", "/"+util.Name, utilHandler(util))
	}

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := addSpanTags(r)
		defer span.Finish(tracer.WithError(errors.New(http.StatusText(http.StatusNotFound))))
		http.NotFoundHandler().ServeHTTP(w, r)
	})

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := addSpanTags(r)
		defer span.Finish(tracer.WithError(errors.New(http.StatusText(http.StatusMethodNotAllowed))))
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Service-Version", version)
		m := httpsnoop.CaptureMetrics(router, w, r)
		log.Printf("%d %s %s %s", m.Code, r.Method, r.URL, m.Duration)
	})
}

func utilHandler(util internal.Util) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addSpanTags(r)

		var err error
		span, _ := tracer.SpanFromContext(r.Context())
		defer func() {
			span.Finish(tracer.WithError(err))
		}()

		out := &bytes.Buffer{}
		a := &internal.UtilArgs{Output: out}

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
			for name, flag := range util.Flags {
				a.Flags[name] = flag.Default
				if _, ok := r.URL.Query()[name]; !ok {
					continue
				}

				qVal := r.URL.Query().Get(name)
				switch flag.Default.(type) {
				case time.Duration:
					dur, err := time.ParseDuration(qVal)
					if err != nil {
						return fmt.Errorf("bad query param %s: %w", name, err)
					}
					a.Flags[name] = dur
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
		err = util.Execute(execCtx, a)
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
