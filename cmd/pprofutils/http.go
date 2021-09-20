package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func newHTTPServer() http.Handler {
	router := httprouter.New()
	for _, cmd := range commands {
		router.POST("/"+cmd.Name, func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

		})
	}
	return router
}
