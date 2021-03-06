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
	"path"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/liimee/idk/board"
	"github.com/liimee/idk/user"
)

type User = user.User

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

type Game struct {
	Players []User
	Start   bool
	Turn    string

	Bid    float64
	Bidd   User
	Biddi  bool
	Biddd  []string
	BidPos int
}

//go:embed build/*
var f embed.FS

var dev = os.Getenv("DEV") == "y"

var gs = map[string]Game{}

var hu *H

func main() {
	e := mux.NewRouter()

	e.Path("/api/new").Methods("POST").HandlerFunc(New)
	e.Path("/api/board").Methods("GET").HandlerFunc(func(w http.ResponseWriter, h *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		b, _ := json.Marshal(board.Board)
		w.Write([]byte(b))
	})
	e.Path("/api/exists/{game}").Methods("GET").HandlerFunc(GameExists)

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
	p, pp := os.LookupEnv("PORT")
	if pp {
		s = ":" + p
	}

	if !dev {
		e.PathPrefix("/").HandlerFunc(fe)
	}

	fmt.Println("running")

	http.Handle("/", e)
	http.ListenAndServe(s, nil)
}

func New(w http.ResponseWriter, h *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", os.Getenv("CL"))
	id := GetRandom()
	gs[id] = Game{Players: []User{}, Start: false, Turn: "", Bid: 0, Bidd: User{}, Biddi: false, Biddd: []string{}, BidPos: 0}
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
	d, _ := filepath.Abs(r.URL.Path)
	_, er := fs.Stat(f, path.Join("build", d))
	if os.IsNotExist(er) {
		g, _ := f.ReadFile("build/index.html")
		w.Write(g)
		return
	}
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
				c.co.WriteJSON(map[string]string{
					"S": "err",
					"E": "The game has started",
				})
				c.co.Close()
				hu.unreg <- c
				return
			}
			random.Seed(time.Now().UnixNano())
			cols := []string{"#B10606", "#B08205", "#05B014", "#0566B0"}
			l.Players = append(gs[s["id"]].Players, User{
				Name:   s["as"],
				Money:  1500,
				Id:     r,
				Pos:    0,
				Owns:   []int{},
				InJail: false,
				Color:  cols[random.Intn(len(cols))],
				Mo:     []int{},
				Ho:     map[int]int{},
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
			ss.Pos += (random.Intn(11) + 1)
			if ss.Pos > 39 {
				ss.Pos -= 40
				ss.Money += 200
			}
			if ss.Pos == 30 {
				ss.Pos = 10
				ss.InJail = true
			}
			if WhoOwnsIt(c.game, ss.Pos) != c.id && WhoOwnsIt(c.game, ss.Pos) != "" && !IsMortgaged(ss.Pos, gs[c.game]) && ss.Pos%10 != 5 {
				ss.Money -= board.Board[ss.Pos].Rent[gs[c.game].Players[GetIndexById(WhoOwnsIt(c.game, ss.Pos), gs[c.game])].Ho[ss.Pos]] // "0"
				gs[c.game].Players[GetIndexById(WhoOwnsIt(c.game, ss.Pos), gs[c.game])].Money += board.Board[ss.Pos].Rent[gs[c.game].Players[GetIndexById(WhoOwnsIt(c.game, ss.Pos), gs[c.game])].Ho[ss.Pos]]
			}
			if ss.Pos == 4 {
				ss.Money -= 200
			}
			if ss.Pos == 38 {
				ss.Money -= 100
			}
			if board.Board[ss.Pos].Name == "Chance" {
				random.Seed(time.Now().UnixNano())
				rk := random.Intn(len(board.ChanceCards))
				c.co.WriteJSON(map[string]string{
					"S":   "card",
					"Str": board.ChanceCards[rk].Str,
					"T":   "a",
				})
				g.Players[GetIndexById(g.Turn, g)] = board.ChanceCards[rk].Fun(GetPlayerById(c.id, g))
			}
			if board.Board[ss.Pos].Name == "Community Chest" {
				random.Seed(time.Now().UnixNano())
				rk := random.Intn(len(board.CommunityChestCards))
				c.co.WriteJSON(map[string]string{
					"S":   "card",
					"Str": board.CommunityChestCards[rk].Str,
					"T":   "b",
				})
				g.Players[GetIndexById(g.Turn, g)] = board.CommunityChestCards[rk].Fun(GetPlayerById(c.id, g))
			}
			if ss.Pos%10 == 5 {
				n := Stations(gs[c.game].Players[GetIndexById(WhoOwnsIt(c.game, ss.Pos), gs[c.game])].Id, c.game)
				ss.Money -= board.Board[ss.Pos].Rent[0] * n
				gs[c.game].Players[GetIndexById(WhoOwnsIt(c.game, ss.Pos), gs[c.game])].Money += board.Board[ss.Pos].Rent[0] * n
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

			if WhoOwnsIt(c.game, GetPlayerById(c.id, n).Pos) == "" && board.Board[GetPlayerById(c.id, n).Pos].Price > 0 {
				n.Bid = 0
				n.Bidd = User{}
				n.Biddi = true
				n.BidPos = GetPlayerById(c.id, n).Pos
				s, _ := json.Marshal(map[string]interface{}{
					"S":     "bid",
					"Bid":   n.Bid,
					"Biddi": n.Biddi,
					"Pa":    n.Biddd,
					"Biddd": n.BidPos,
				})
				gs[c.game].BcGame(s)
				gs[c.game] = n
			} else {
				d := GetIndexById(n.Turn, n) + 1
				if d >= len(n.Players) {
					d = 0
				}
				n.Turn = n.Players[d].Id
				gs[c.game] = n

				Turn(c.game)
			}
		} else if s["s"] == "buy" {
			if c.id != gs[c.game].Turn || GetPlayerById(c.id, gs[c.game]).Money < board.Board[GetPlayerById(c.id, gs[c.game]).Pos].Price {
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
		} else if s["s"] == "bid" {
			if gs[c.game].Biddi {
				var s map[string]interface{}
				json.Unmarshal(m, &s)
				f := gs[c.game]
				if s["bid"].(float64) > gs[c.game].Bid && s["bid"].(float64) <= float64(GetPlayerById(c.id, f).Money) {
					f.Bid = s["bid"].(float64)
					f.Bidd = GetPlayerById(c.id, f)
					f.Biddd = []string{}
					gs[c.game] = f

					s, _ := json.Marshal(map[string]interface{}{
						"S":     "bid",
						"Bid":   f.Bid,
						"Biddi": f.Biddi,
						"Pa":    f.Biddd,
						"Biddd": f.BidPos,
					})
					gs[c.game].BcGame(s)
				}
			}
		} else if s["s"] == "pass" {
			f := gs[c.game]
			f.Biddd = append(f.Biddd, c.id)

			if len(f.Biddd) == len(f.Players) {
				if f.Bid > 0 {
					mmm := f.Players[GetIndexById(f.Bidd.Id, f)]
					mmm.Money -= int(f.Bid)
					mmm.Owns = append(mmm.Owns, f.BidPos)
					f.Players[GetIndexById(f.Bidd.Id, f)] = mmm
				}

				f.Bid = 0
				f.Bidd = User{}
				f.Biddd = []string{}
				f.Biddi = false

				d := GetIndexById(f.Turn, f) + 1
				if d >= len(f.Players) {
					d = 0
				}
				f.Turn = f.Players[d].Id
				gs[c.game] = f

				Turn(c.game)
			}

			gs[c.game] = f

			s, _ := json.Marshal(map[string]interface{}{
				"S":     "bid",
				"Bid":   f.Bid,
				"Biddi": f.Biddi,
				"Pa":    f.Biddd,
				"Biddd": f.BidPos,
			})
			gs[c.game].BcGame(s)

			br, _ := json.Marshal(struct {
				S    string
				Data []User
			}{S: "data", Data: gs[c.game].Players})
			gs[c.game].BcGame(br)
		} else if s["s"] == "payjail" {
			m := GetPlayerById(c.id, gs[c.game])
			m.Money -= 50
			m.InJail = false
			gs[c.game].Players[GetIndexById(c.id, gs[c.game])] = m

			br, _ := json.Marshal(struct {
				S    string
				Data []User
			}{S: "data", Data: gs[c.game].Players})
			gs[c.game].BcGame(br)
		} else if s["s"] == "resign" {
			d := gs[c.game]
			mm := GetIndexById(c.id, d)
			d.Players = RemovePlayer(c.game, c.id)
			gs[c.game] = d
			br, _ := json.Marshal(struct {
				S    string
				Data []User
			}{S: "data", Data: gs[c.game].Players})
			gs[c.game].BcGame(br)

			if len(gs[c.game].Players) == 1 {
				co := GetClisById(d.Players[0].Id).co
				co.WriteJSON(map[string]string{
					"S": "win",
				})
				co.Close()
			} else {
				mm += 1
				if mm >= len(d.Players) {
					mm = 0
				}
				d.Turn = d.Players[mm].Id
				gs[c.game] = d

				Turn(c.game)
			}
		} else if s["s"] == "mortgage" {
			var s map[string]interface{}
			json.Unmarshal(m, &s)
			pos := int(s["pos"].(float64))

			or := false
			for _, v := range GetPlayerById(c.id, gs[c.game]).Owns {
				if board.Board[v].Set == board.Board[pos].Set {
					if GetPlayerById(c.id, gs[c.game]).Ho[v] > 0 {
						or = true
					}
				}
			}

			if WhoOwnsIt(c.game, pos) == c.id && !or {
				n := gs[c.game]
				u := &n.Players[GetIndexById(c.game, n)]
				if !IsMortgaged(pos, gs[c.game]) {
					u.Money += (board.Board[pos].Price / 2)
					u.Mo = append(u.Mo, pos)
				} else {
					ps := board.Board[pos].Price / 2
					u.Money -= (ps + int((0.1 * float64(ps))))
					index := 0
					for i, vv := range u.Mo {
						if vv == pos {
							index = i
							break
						}
					}
					u.Mo = append(u.Mo[:index], u.Mo[index+1:]...)
				}

				gs[c.game] = n

				br, _ := json.Marshal(struct {
					S    string
					Data []User
				}{S: "data", Data: gs[c.game].Players})
				gs[c.game].BcGame(br)
			}
		} else if s["s"] == "ho" {
			var s map[string]interface{}
			json.Unmarshal(m, &s)
			pos := int(s["pos"].(float64))

			pl := GetPlayerById(c.id, gs[c.game])

			if _, ok := pl.Ho[pos]; !ok {
				pl.Ho[pos] = 0
			}

			if !IsMortgaged(pos, gs[c.game]) && SameOwner(pos, pl, gs[c.game]) && pl.Ho[pos] < 5 {
				switch board.Board[pos].Set {
				case 1, 2:
					pl.Money -= 50
				case 3, 4:
					pl.Money -= 100
				case 5, 6:
					pl.Money -= 150
				case 7, 8:
					pl.Money -= 200
				}

				pl.Ho[pos] += 1
			}

			g := gs[c.game]
			g.Players[GetIndexById(c.id, g)] = pl
			gs[c.game] = g

			br, _ := json.Marshal(struct {
				S    string
				Data []User
			}{S: "data", Data: gs[c.game].Players})
			gs[c.game].BcGame(br)
		} else if s["s"] == "sell" {
			var s map[string]interface{}
			json.Unmarshal(m, &s)
			pos := int(s["pos"].(float64))

			pl := GetPlayerById(c.id, gs[c.game])

			if _, ok := pl.Ho[pos]; !ok {
				pl.Ho[pos] = 0
			}

			if !IsMortgaged(pos, gs[c.game]) && SameOwner(pos, pl, gs[c.game]) && pl.Ho[pos] > 0 {
				switch board.Board[pos].Set {
				case 1, 2:
					pl.Money += 50 / 2
				case 3, 4:
					pl.Money += 100 / 2
				case 5, 6:
					pl.Money += 150 / 2
				case 7, 8:
					pl.Money += 200 / 2
				}

				pl.Ho[pos] -= 1
			}

			g := gs[c.game]
			g.Players[GetIndexById(c.id, g)] = pl
			gs[c.game] = g

			br, _ := json.Marshal(struct {
				S    string
				Data []User
			}{S: "data", Data: gs[c.game].Players})
			gs[c.game].BcGame(br)
		}
	}
}

func SameOwner(p int, u User, gm Game) bool {
	set := board.Board[p].Set

	or := 0
	for _, v := range board.Board {
		if v.Set == set {
			or++
		}
	}

	pp := 0
	for _, v := range u.Owns {
		if board.Board[v].Set == set && !IsMortgaged(v, gm) {
			pp++
		}
	}

	return or == pp
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

func Stations(g string, j string) int {
	var s = 0
	for _, k := range []int{5, 15, 25, 35} {
		if WhoOwnsIt(j, k) == g {
			s++
		}
	}
	return s
}

func IsMortgaged(id int, g Game) bool {
	s := false
	for _, u := range g.Players {
		for _, p := range u.Mo {
			if p == id {
				s = true
				break
			}
		}
	}
	return s
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

func GameExists(w http.ResponseWriter, h *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", os.Getenv("CL"))
	p := mux.Vars(h)
	_, ok := gs[p["game"]]
	j, _ := json.Marshal(map[string]bool{
		"Exists": ok,
	})
	w.Write(j)
}
