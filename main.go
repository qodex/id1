package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var version = "latest"
var port = "8080"
var dbpath = "/mnt/id1db"

func main() {
	if godotenv.Load(".env") == nil {
		port = os.Getenv("PORT")
		dbpath = os.Getenv("DBPATH")
	}

	fmt.Printf("id1 API build %s, port: %s, dbpath: %s\n", version, port, dbpath)

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

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("error starting service: %s", err)
	}
}
