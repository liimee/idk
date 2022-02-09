package main

import (
	"embed"
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

	s := ":3000"
	if dev {
		s = ":4000"
	}

	fmt.Println("running")

	http.Handle("/", e)
	http.ListenAndServe(s, nil)
}

func fe(w http.ResponseWriter, r *http.Request) {
	s, _ := fs.Sub(f, "build")
	http.FileServer(http.FS(s)).ServeHTTP(w, r)
}
