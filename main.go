package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var version = "latest"
var port = "8080"
var dbpath = "/mnt/id1db"
var pubsub PubSub

func main() {
	if godotenv.Load(".env") == nil {
		port = os.Getenv("PORT")
		dbpath = os.Getenv("DBPATH")
	}
	fmt.Printf("id1 API build %s, port: %s, dbpath: %s\n", version, port, dbpath)
	pubsub = NewPubSub()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ok200(w, fmt.Appendf(nil, "id1 api v.%s", version))
	})

	http.HandleFunc("/{id}/{key...}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			ok200(w, []byte{})
			return
		}
		req := NewRequestProps(r)

		id := ""
		if claims, _ := validateToken(req.Token, ""); len(claims.Subject) > 0 {
			id = claims.Subject
			if _, err := validateToken(req.Token, generateSecret(id)); err != nil {
				id = ""
			}
		}
		req.Cmd.Args["x-id"] = id

		authOk := auth(id, req.Cmd)

		if !authOk {
			if len(id) > 0 {
				err403(w, "")
			} else if pubKey, err := CmdGet(KK(req.Id, "pub", "key")).Exec(); err == nil {
				if challenge, err := generateChallenge(req.Id, string(pubKey)); err == nil {
					err401(w, challenge)
				} else {
					log.Printf("%s: %s", err, id)
					err500(w, err.Error())
				}
			} else {
				err404(w, err.Error())
			}
			return
		}

		if req.IsWebSocket {
			upgrader := websocket.Upgrader{}
			upgrader.CheckOrigin = func(r *http.Request) bool { return true }
			wsHandler := webSocketHandler{
				upgrader: upgrader,
			}
			wsHandler.Handle(w, r)
		} else if data, err := req.Cmd.Exec(); err == nil {
			ok200(w, data)
		} else if errors.Is(err, ErrNotFound) {
			err404(w, "")
		} else if errors.Is(err, ErrLimitExceeded) {
			err413(w, "")
		} else {
			err400(w, err.Error())
		}
	})

	go func() {
		for {
			go dotAfter(dbpath)
			time.Sleep(time.Second * 10)
		}
	}()

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("error starting service: %s", err)
	}
}
