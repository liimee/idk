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
	"github.com/liimee/idk/board"
)

type User struct {
	Name string
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
		b, _ := json.Marshal(board.Board)
		w.Write([]byte(b))
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
	b := make([]byte, 6)
	rand.Read(b)
	id := hex.EncodeToString(b)
	gs[id] = Game{Players: []User{}}
	m, _ := json.Marshal(map[string]string{"id": id})
	w.Write([]byte(m))
}

func fe(w http.ResponseWriter, r *http.Request) {
	s, _ := fs.Sub(f, "build")
	http.FileServer(http.FS(s)).ServeHTTP(w, r)
}
