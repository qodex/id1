package main

import (
	"fmt"
	"log"
	"net/http"
)

var version = "latest"
var port = "8080"

func main() {
	log.Printf("id1 API %s, port: %s\n\n", version, port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "id1 api v", version)
	})

	http.ListenAndServe(":"+port, nil)
}
