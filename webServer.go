package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Logging() Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() { fmt.Println(r.URL.Path, time.Since(start)) }()

			hf(w, r)
		}
	}
}

func Chain(hf http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		hf = m(hf)
	}
	return hf
}

func StartSever() {

	// load template files
	tmplRoot := template.Must(template.ParseFiles("templates/root.html"))
	tmplChart := template.Must(template.ParseFiles("templates/chart.html"))

	rootRoute := func(rw http.ResponseWriter, r *http.Request) {
		tmplRoot.Execute(rw, nil)
	}

	chartRoute := func(wr http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmplRoot.Execute(wr, nil)
			return
		}
		tmplChart.Execute(wr, nil)
	}

	http.HandleFunc("/", Chain(rootRoute, Logging()))
	http.HandleFunc("/chart", Chain(chartRoute, Logging()))
	http.ListenAndServe(":4567", nil)
}
