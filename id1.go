package id1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var pubsub = NewPubSub()
var dbpath = "/mnt/id1db"
var version = "latest"

func Handle(path string, ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	dbpath = path

	go func() {
		for {
			select {
			case <-time.After(time.Second * 10):
				dotAfter(dbpath)
			case <-ctx.Done():
				return
			}
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			ok200(w, []byte{})
			return
		}

		req := NewRequestProps(r)

		if req.Id == "" {
			ok200(w, fmt.Appendf(nil, "Qodex id1 %s https://github.com/qodex/id1", version))
			return
		}

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
	}
}
