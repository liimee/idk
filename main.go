package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
)

//go:embed build/*
var f embed.FS

func main() {
	e := mux.NewRouter()

	e.PathPrefix("/").HandlerFunc(fe)

	fmt.Println("running")

	http.Handle("/", e)
	http.ListenAndServe(":3000", nil)
}

func fe(w http.ResponseWriter, r *http.Request) {
	s, _ := fs.Sub(f, "build")
	http.FileServer(http.FS(s)).ServeHTTP(w, r)
}
