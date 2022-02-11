package main

import (
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	random "math/rand"
	"net/http"
	"os"
	"time"

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
	Pos   int
	Id    string
	Owns  []int
}

type Game struct {
	Players []User
	Start   bool
	Turn    string
}

//go:embed build/*
var f embed.FS

var dev = os.Getenv("DEV") == "y"

var gs = map[string]Game{}

var hu *H

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

	hu = &H{
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
	gs[id] = Game{Players: []User{}, Start: false, Turn: ""}
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

	client := &Cli{n: hu, co: c, send: make(chan []byte, 256), id: "", game: ""}
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
			if c.game != "" {
				d := gs[c.game]
				d.Players = RemovePlayer(c.game, c.id)
				gs[c.game] = d
				br, _ := json.Marshal(struct {
					S    string
					Data []User
				}{S: "data", Data: gs[c.game].Players})
				gs[c.game].BcGame(br)
			}
			hu.unreg <- c
			break
		}

		var s map[string]string
		json.Unmarshal(m, &s)

		if s["s"] == "join" {
			r := GetRandom()
			l := gs[s["id"]]
			if l.Start {
				c.co.Close()
				hu.unreg <- c
				return
			}
			l.Players = append(gs[s["id"]].Players, User{
				Name:  s["as"],
				Money: 1500,
				Id:    r,
				Pos:   0,
				Owns:  []int{},
			})
			gs[s["id"]] = l
			c.id = r
			c.game = s["id"]

			c.co.WriteJSON(map[string]string{"S": "id", "Id": c.id})
			br, _ := json.Marshal(struct {
				S    string
				Data []User
			}{S: "data", Data: gs[c.game].Players})
			gs[c.game].BcGame(br)
		} else if s["s"] == "start" {
			n := gs[c.game]
			n.Start = true
			n.Turn = gs[c.game].Players[0].Id
			gs[c.game] = n

			mr, _ := json.Marshal(struct {
				S     string
				Start bool
			}{S: "start", Start: gs[c.game].Start})

			gs[c.game].BcGame(mr)

			Turn(c.game)
		} else if s["s"] == "roll" {
			if c.id != gs[c.game].Turn {
				return
			}

			g := gs[c.game]
			ss := &g.Players[GetIndexById(g.Turn, g)]
			random.Seed(time.Now().UnixNano())
			ss.Pos += (random.Intn(5) + 1)
			if ss.Pos > 39 {
				ss.Pos -= 39
				ss.Money += 200
			}
			if board.Board[ss.Pos].Name == "Tax" || board.Board[ss.Pos].Name == "Tax(i)" {
				ss.Money -= board.Board[ss.Pos].Price
			}
			if WhoOwnsIt(c.game, ss.Pos) != c.id && WhoOwnsIt(c.game, ss.Pos) != "" {
				ss.Money -= board.Board[ss.Pos].Rent[0] // "0"
				gs[c.game].Players[GetIndexById(WhoOwnsIt(c.game, ss.Pos), gs[c.game])].Money += board.Board[ss.Pos].Rent[0]
			}
			gs[c.game] = g

			br, _ := json.Marshal(struct {
				S    string
				Data []User
			}{S: "data", Data: gs[c.game].Players})
			gs[c.game].BcGame(br)
		} else if s["s"] == "endturn" {
			if c.id != gs[c.game].Turn {
				return
			}

			n := gs[c.game]
			d := GetIndexById(n.Turn, n) + 1
			if d >= len(n.Players) {
				d = 0
			}
			n.Turn = n.Players[d].Id
			gs[c.game] = n

			Turn(c.game)
		} else if s["s"] == "buy" {
			if c.id != gs[c.game].Turn {
				return
			}

			//TODO: check

			s := gs[c.game]
			m := s.Players[GetIndexById(c.id, s)]
			m.Owns = append(m.Owns, m.Pos)
			m.Money -= board.Board[m.Pos].Price
			s.Players[GetIndexById(c.id, s)] = m
			gs[c.game] = s

			br, _ := json.Marshal(struct {
				S    string
				Data []User
			}{S: "data", Data: gs[c.game].Players})
			gs[c.game].BcGame(br)
		}
	}
}

func Turn(g string) {
	GetClisById(gs[g].Turn).co.WriteJSON(map[string]string{
		"S": "turn",
	})
}

func GetPlayerById(s string, g Game) User {
	var j User
	for _, v := range g.Players {
		if v.Id == s {
			j = v
			break
		}
	}
	return j
}

func GetIndexById(s string, g Game) int {
	var j int
	for v, d := range g.Players {
		if d.Id == s {
			j = v
			break
		}
	}
	return j
}

func GetClisById(s string) *Cli {
	var n *Cli
	for f := range hu.clis {
		if f.id == s {
			n = f
		}
	}
	return n
}

func (g Game) BcGame(d []byte) {
	for _, f := range g.Players {
		GetClisById(f.Id).co.WriteMessage(websocket.TextMessage, d)
	}
}

func WhoOwnsIt(g string, n int) string {
	var f string
	for _, s := range gs[g].Players {
		for _, j := range s.Owns {
			if n == j {
				f = s.Id
				break
			}
		}
	}
	return f
}

func RemovePlayer(s string, d string) []User {
	n := gs[s].Players
	for k, f := range gs[s].Players {
		if f.Id == d {
			n = append(n[:k], n[k+1:]...)
		}
	}
	return n
}
