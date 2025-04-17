package main

import (
	"net/http"
)

func cors(rw *http.ResponseWriter) {
	w := *rw
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, PUT, PATCH, POST, DELETE, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "*")
}

func ok200(rw http.ResponseWriter, data []byte) {
	cors(&rw)
	rw.Write(data)
}

func err404(rw http.ResponseWriter, s string) {
	cors(&rw)
	http.Error(rw, s, http.StatusNotFound)
}

func err401(rw http.ResponseWriter, s string) {
	cors(&rw)
	http.Error(rw, s, http.StatusUnauthorized)
}

func err403(rw http.ResponseWriter, s string) {
	cors(&rw)
	http.Error(rw, s, http.StatusForbidden)
}

func err413(rw http.ResponseWriter, s string) {
	cors(&rw)
	http.Error(rw, s, http.StatusRequestEntityTooLarge)
}

func err400(rw http.ResponseWriter, s string) {
	cors(&rw)
	http.Error(rw, s, http.StatusBadRequest)
}
