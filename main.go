package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

var version = "latest"
var port = "8080"
var dbpath = "id1db"

func main() {
	log.Printf("id1 API %s, port: %s\n\n", version, port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ok200(w, fmt.Appendf(nil, "id1 api v.%s", version))
	})

	http.HandleFunc("/{id}/{key...}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			ok200(w, []byte{})
			return
		}

		req := NewRequestProps(r)
		if data, err := req.Cmd.Exec(); err == nil {
			ok200(w, data)
		} else if errors.Is(err, ErrNotFound) {
			err404(w, "")
		} else if errors.Is(err, ErrLimitExceeded) {
			err413(w, "")
		} else {
			err400(w, err.Error())
		}
	})

	http.ListenAndServe(":"+port, nil)
}
