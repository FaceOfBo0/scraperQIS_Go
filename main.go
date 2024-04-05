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

	r.HandleFunc("/", func(wrt http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(wrt, "Hi you requested the following endpoint: %s\n", req.URL.Path)
	}).Methods("GET")

	r.HandleFunc("/hello", func(wrt http.ResponseWriter, req *http.Request) {
		userid := If(req.URL.Query().Has("id"), req.URL.Query().Get("id"), "null")
		println(userid)
		wrt.Header().Set("Content-Type", "text/html")
		wrt.Write([]byte("<html><h2>Hello</h2></html>"))
	})

	r.HandleFunc("/books/{title}/page/{page}", func(wrt http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(wrt, "You've requested the book: %s on page %s\n", title, page)
	})

	fs := http.FileServer(http.Dir("static/"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))
	r.Handle("/static", http.StripPrefix("/static", fs))
	http.ListenAndServe(":5000", r)
}
