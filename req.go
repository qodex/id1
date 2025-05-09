package id1

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var opMap = map[string]Op{
	http.MethodGet:    Get, // also list
	http.MethodPost:   Set,
	http.MethodDelete: Del,
	http.MethodPatch:  Add, // also mov
}

type RequestProps struct {
	Id          string
	Token       string
	IsWebSocket bool
	Cmd         Command
}

func (t RequestProps) String() string {
	return fmt.Sprintf("req cmd: %s %s (%s): %s", t.Cmd.Op, t.Cmd.Key, t.Cmd.Args, string(t.Cmd.Data))
}

func NewRequestProps(r *http.Request) RequestProps {
	key := K(r.URL.Path)
	req := RequestProps{
		Id: key.Id,
		Cmd: Command{
			Op:   opMap[r.Method],
			Key:  key,
			Args: map[string]string{},
			Data: []byte{},
		},
	}

	paramPairs := r.URL.Query()
	for key, values := range paramPairs {
		req.Cmd.Args[key] = values[0]
	}

	if r.Header["Authorization"] != nil && len(r.Header["Authorization"]) > 0 {
		req.Token = strings.TrimPrefix(r.Header["Authorization"][0], "Bearer ")
	}

	if r.Header["Upgrade"] != nil {
		req.IsWebSocket = true
	}

	if data, readErr := io.ReadAll(r.Body); readErr != nil {
		log.Printf("error reading request body, %s", readErr)
	} else {
		req.Cmd.Data = data
	}

	if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "*") {
		req.Cmd.Op = List
	}
	if r.Method == http.MethodPatch && len(r.Header["X-Move-To"]) > 0 {
		req.Cmd.Op = Mov
		req.Cmd.Data = []byte(r.Header["X-Move-To"][0])
	}

	return req
}
