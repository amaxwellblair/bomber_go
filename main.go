package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"golang.org/x/net/websocket"
)

type handler struct {
	mu   sync.Mutex
	room map[string][]*websocket.Conn
}

type commands struct {
	ID        string `json:"id"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Direction string `json:"direction"`
	Type      string `json:"type"`
}

func newHandler() *handler {
	return &handler{
		room: make(map[string][]*websocket.Conn),
	}
}

func (h *handler) socket(ws *websocket.Conn) {
	url := ws.Config().Location.String()
	h.mu.Lock()
	h.room[url] = append(h.room[url], ws)
	h.mu.Unlock()

	for {
		var c commands
		if err := websocket.JSON.Receive(ws, &c); err != nil {
			fmt.Println(err)
			break
		}

		h.sendCommands(&c, url)
		fmt.Println("Received commands:", c)
	}
}

func (h *handler) sendCommands(c *commands, url string) {
	h.mu.Lock()
	connections := h.room[url]
	h.mu.Unlock()

	for _, conn := range connections {
		if err := websocket.JSON.Send(conn, c); err != nil {
			fmt.Println(err)
		}
	}
}

func (h *handler) newRouter() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/socket/{id}", websocket.Handler(h.socket))
	return r
}

func main() {
	h := newHandler()
	r := h.newRouter()
	http.ListenAndServe(":9000", r)
}
