package web

import (
	"net/http"
)

func init() {
	http.HandleFunc("/", DefaultHandler)
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		IndexHandler(w, r)
	} else {
		// Response 404.
		HandleNotFound(w, r)
	}
}
