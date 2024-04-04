package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
}

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func main() {
	fmt.Println("hello Server")
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi you requested the following endpoint: %s\n", r.URL.Path)
	}).Methods("GET")

	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		userid := If(r.URL.Query().Has("id"), r.URL.Query().Get("id"), "null")
		println(userid)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html><h2>Hello</h2></html>"))
	})

	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	fs := http.FileServer(http.Dir("static/"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))
	r.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":5000", r)
}
