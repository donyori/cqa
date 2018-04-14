package web

import (
	"fmt"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	e := Render(w, "error.tmpl", &ErrorData{
		StatusCode: statusCode,
		Msg:        err.Error(),
	})
	if e != nil {
		// This is same as http.Error() without w.WriteHeader(statusCode).
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		fmt.Fprintln(w, err)
	}
}

func HandleBadRequest(w http.ResponseWriter, err error) {
	HandleError(w, err, http.StatusBadRequest)
}

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	HandleError(w, fmt.Errorf("%q is NOT found", r.URL.Path),
		http.StatusNotFound)
}

func HandleInternalServerError(w http.ResponseWriter, err error) {
	HandleError(w, err, http.StatusInternalServerError)
}
