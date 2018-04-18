package web

import (
	"net/http"
	"sync"

	_ "net/http/pprof"
)

var setUrlsOnce sync.Once

func SetUrls() {
	setUrlsOnce.Do(func() {
		// Static resources server:
		srs := http.FileServer(GlobalSettings.StaticResourcesRoot)
		http.Handle("/sr/", http.StripPrefix("/sr/", srs))

		// Default handler:
		http.HandleFunc("/", DefaultHandler)

		// Following URLs are added in package "net/http/pprof":
		//   http.HandleFunc("/debug/pprof/", Index)
		//   http.HandleFunc("/debug/pprof/cmdline", Cmdline)
		//   http.HandleFunc("/debug/pprof/profile", Profile)
		//   http.HandleFunc("/debug/pprof/symbol", Symbol)
		//   http.HandleFunc("/debug/pprof/trace", Trace)
	})
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		IndexHandler(w, r)
	} else {
		// Response 404.
		HandleNotFound(w, r)
	}
}
