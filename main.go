package main

import (
	"fmt"
	"log"
	"net/http"
)

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
}

func main() {
	fmt.Println("hello World")
	http.HandleFunc("/hello", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/html")
		rw.Write([]byte("<html><h2>Hello</h2></html>"))
	})

	log.Fatal(http.ListenAndServe(":5000", nil))
}
