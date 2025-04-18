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
		defer conn.Close()
		defer onDisconnect(id)
		onConnect(id)
		cmdOut := pubsub.Subscribe(id)
		defer pubsub.Unsubscribe(id, cmdOut)

		go handleCommands(cmdIn, cmdOut, ctx, cancel)
		go readCommands(conn, cmdIn, cancel)
		go writeCommands(cmdOut, conn, ctx, cancel)
		go ping(id, cmdOut, ctx)
	}
	<-ctx.Done()
	cancel()
}

func onConnect(id string) {
	log.Printf("connected: %s", id)
	if _, err := CmdSet(KK(id, "online"), []byte{}).Exec(); err != nil {
		log.Printf("cmd set error: %s", err)
	}
}

func onDisconnect(id string) {
	log.Printf("disconnected: %s", id)
	if _, err := CmdDel(KK(id, "online")).Exec(); err != nil {
		log.Printf("cmd del error: %s", err)
	}
}

func ping(id string, cmdOut chan (Command), ctx context.Context) {
	for {
		select {
		case <-time.After(time.Second * 60):
			cmdOut <- CmdGet(KK(id, ".ping"))
		case <-ctx.Done():
			return
		}
	}
}

func readCommands(conn *websocket.Conn, cmdInCh chan (Command), cancel func()) {
	for {
		if _, data, err := conn.ReadMessage(); err != nil {
			conn.Close()
			cancel()
			return
		} else if cmd, err := ParseCommand(data); err != nil {
			log.Printf("error parsing websocket message: %s", err)
		} else {
			cmdInCh <- cmd
		}
	}
}

func writeCommands(cmdOut chan (Command), conn *websocket.Conn, ctx context.Context, cancel func()) {
	for {
		select {
		case cmd := <-cmdOut:
			if err := conn.WriteMessage(websocket.BinaryMessage, cmd.Bytes()); err != nil {
				cancel()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func handleCommands(cmdIn chan (Command), cmdOut chan (Command), ctx context.Context, cancel func()) {
	timeout := time.Second * 600
	for {
		select {
		case cmd := <-cmdIn:
			if data, err := cmd.Exec(); err == nil {
				cmdOut <- CmdSet(cmd.Key, data)
			} else if errors.Is(ErrNotFound, err) {
				cmdOut <- CmdDel(cmd.Key)
			} else {
				log.Printf("error executing command: %s", err)
			}
		case <-time.After(timeout):
			log.Printf("no commands, disconnecting...")
			cancel()
			return
		case <-ctx.Done():
			return
		}
	}
}
