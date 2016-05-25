package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"golang.org/x/net/websocket"
)

type handler struct {
	connections []*websocket.Conn
}

type commands struct {
	ID        string `json:"id"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Direction string `json:"direction"`
	Type      string `json:"type"`
}

func newHandler() *handler {
	return &handler{}
}

func (h *handler) socket(ws *websocket.Conn) {
	h.connections = append(h.connections, ws)
	for {
		var c commands
		if err := websocket.JSON.Receive(ws, &c); err != nil {
			fmt.Println(err)
			break
		}

		h.sendCommands(&c)
		fmt.Println("Received commands:", c)
	}
}

func (h *handler) sendCommands(c *commands) {
	for _, conn := range h.connections {
		if err := websocket.JSON.Send(conn, c); err != nil {
			fmt.Println(err)
		}
	}
}

func (h *handler) newRouter() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/socket", websocket.Handler(h.socket))
	return r
}

func main() {
	h := newHandler()
	r := h.newRouter()
	http.ListenAndServe(":9000", r)
}
