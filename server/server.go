package server

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

// Run starts the epoxy web server
func Run(stop chan struct{}) {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	glog.Fatal(srv.ListenAndServe())
}
