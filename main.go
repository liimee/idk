package main

import (
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/liimee/idk/board"
)

type H struct {
	clis  map[*Cli]bool
	bc    chan []byte
	reg   chan *Cli
	unreg chan *Cli
}

type Cli struct {
	n *H

	co *websocket.Conn

	send chan []byte

	game string
	id   string
}

type User struct {
	Name  string
	Money int
	Id    string
}

type Game struct {
	Players []User
}

//go:embed build/*
var f embed.FS

var dev = os.Getenv("DEV") == "y"

var gs = map[string]Game{}

func main() {
	e := mux.NewRouter()

	if !dev {
		e.PathPrefix("/").HandlerFunc(fe)
	}

	e.Path("/api/new").Methods("POST").HandlerFunc(New)
	e.Path("/api/board").Methods("GET").HandlerFunc(func(w http.ResponseWriter, h *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		b, _ := json.Marshal(board.Board)
		w.Write([]byte(b))
	})

	hu := &H{
		bc:    make(chan []byte),
		reg:   make(chan *Cli),
		unreg: make(chan *Cli),
		clis:  make(map[*Cli]bool),
	}
	go hu.Run()
	e.Path("/api/ws").HandlerFunc(func(w http.ResponseWriter, h *http.Request) {
		Ws(hu, w, h)
	})

	s := ":3000"
	if dev {
		s = ":4000"
	}

	fmt.Println("running")

	http.Handle("/", e)
	http.ListenAndServe(s, nil)
}

func New(w http.ResponseWriter, h *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", os.Getenv("CL"))
	id := GetRandom()
	gs[id] = Game{Players: []User{}}
	m, _ := json.Marshal(map[string]string{"id": id})
	w.Write([]byte(m))
}

func GetRandom() string {
	b := make([]byte, 6)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func Ws(hu *H, w http.ResponseWriter, h *http.Request) {
	s := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	c, _ := s.Upgrade(w, h, nil)

	client := &Cli{n: hu, co: c, send: make(chan []byte, 256)}
	client.n.reg <- client

	go client.ReadWs()
}

func fe(w http.ResponseWriter, r *http.Request) {
	s, _ := fs.Sub(f, "build")
	http.FileServer(http.FS(s)).ServeHTTP(w, r)
}

func (h *H) Run() {
	for {
		select {
		case client := <-h.reg:
			h.clis[client] = true
		case client := <-h.unreg:
			if _, ok := h.clis[client]; ok {
				delete(h.clis, client)
				close(client.send)
			}
		}
	}
}

func (c *Cli) ReadWs() {
	for {
		_, m, err := c.co.ReadMessage()
		if err != nil {
			break
		}

		var s map[string]string
		json.Unmarshal(m, &s)

		if _, s := gs[s["id"]]; !s {
			panic("No.")
			//some error message via ws
		}

		if s["s"] == "join" {
			r := GetRandom()
			l := gs[s["id"]]
			l.Players = append(gs[s["id"]].Players, User{
				Name:  s["as"],
				Money: 1500,
				Id:    r,
			})
			gs[s["id"]] = l
			c.id = r
			c.game = s["id"]

			c.co.WriteMessage(websocket.TextMessage, []byte(r))
		}
	}
}
