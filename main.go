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
)

//go:embed build/*
var f embed.FS
var dev = os.Getenv("DEV") == "y"

func main() {
	e := mux.NewRouter()

	if !dev {
		e.PathPrefix("/").HandlerFunc(fe)
	}

	e.Path("/api/new").Methods("POST").HandlerFunc(New)

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
	m, _ := json.Marshal(map[string]string{"id": hex.EncodeToString(b)})
	w.Write([]byte(m))
}

func fe(w http.ResponseWriter, r *http.Request) {
	s, _ := fs.Sub(f, "build")
	http.FileServer(http.FS(s)).ServeHTTP(w, r)
}
