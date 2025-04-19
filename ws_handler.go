package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

func (t webSocketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx, cancel := context.WithCancel(context.Background())
	cmdIn := make(chan (Command))

	if conn, err := t.upgrader.Upgrade(w, r, nil); err != nil {
		log.Printf("error upgrading to websocket. %s", err)
	} else {
		session := Session{
			Id:     id,
			Conn:   conn,
			CmdOut: pubsub.Subscribe(id),
			CmdIn:  cmdIn,
			Ctx:    ctx,
			Cancel: cancel,
		}
		defer session.Disconnect()
		session.OnConnect()
	}
	<-ctx.Done()
	cancel()
}

type Session struct {
	Id     string
	Conn   *websocket.Conn
	CmdIn  chan Command
	CmdOut chan Command
	Ctx    context.Context
	Cancel func()
}

func (t *Session) OnConnect() {
	go t.handleCommands()
	go t.readCommands()
	go t.writeCommands()
	go t.ping()

	log.Printf("connected: %s", t.Id)
	if _, err := CmdSet(KK(t.Id, "online"), []byte{}).Exec(); err != nil {
		log.Printf("cmd set error: %s", err)
	}
}

func (t *Session) Disconnect() {
	t.Conn.Close()
	pubsub.Unsubscribe(t.Id, t.CmdOut)
	if _, err := CmdDel(KK(t.Id, "online")).Exec(); err != nil {
		log.Printf("cmd del error: %s", err)
	}
	log.Printf("disconnected: %s", t.Id)
}

func (t *Session) ping() {
	for {
		select {
		case <-time.After(time.Second * 60):
			t.CmdOut <- CmdGet(KK(t.Id, ".ping"))
		case <-t.Ctx.Done():
			return
		}
	}
}

func (t *Session) readCommands() {
	for {
		if _, data, err := t.Conn.ReadMessage(); err != nil {
			t.Conn.Close()
			t.Cancel()
			return
		} else if cmd, err := ParseCommand(data); err != nil {
			log.Printf("error parsing websocket message: %s", err)
		} else {
			t.CmdIn <- cmd
		}
	}
}

func (t *Session) writeCommands() {
	for {
		select {
		case cmd := <-t.CmdOut:
			if err := t.Conn.WriteMessage(websocket.BinaryMessage, cmd.Bytes()); err != nil {
				t.Cancel()
				return
			}
		case <-t.Ctx.Done():
			return
		}
	}
}

func (t *Session) handleCommands() {
	timeout := time.Second * 600
	for {
		select {
		case cmd := <-t.CmdIn:
			authOk := auth(t.Id, cmd)

			if !authOk {
				if pubKey, err := CmdGet(KK(t.Id, "pub", "key")).Exec(); err == nil {
					if challenge, err := generateChallenge(t.Id, string(pubKey)); err == nil {
						t.CmdOut <- CmdSet(KK(t.Id, "auth"), []byte(challenge))
					} else {
						log.Println(err)
					}
				} else {
					log.Println(err)
					t.CmdOut <- CmdDel(KK(t.Id, "pub", "key"))
				}
				continue
			}

			if data, err := cmd.Exec(); err == nil {
				t.CmdOut <- CmdSet(cmd.Key, data)
			} else if errors.Is(ErrNotFound, err) {
				t.CmdOut <- CmdDel(cmd.Key)
			} else {
				log.Printf("error executing command: %s", err)
			}
		case <-time.After(timeout):
			log.Printf("no commands, disconnecting...")
			t.Cancel()
			return
		case <-t.Ctx.Done():
			return
		}
	}
}
