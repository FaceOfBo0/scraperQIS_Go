package main

import (
	"html/template"
	"net/http"
	"slices"
)

/* type Middleware func(http.HandlerFunc) http.HandlerFunc

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
} */

func RunServer() {

	// load template files
	tmplRoot := template.Must(template.ParseFiles("templates/root.html"))
	tmplChart := template.Must(template.ParseFiles("templates/chart.html"))

	mux := http.NewServeMux()

	rootRoute := func(rw http.ResponseWriter, r *http.Request) {
		tmplRoot.Execute(rw, nil)
	}

	chartRoute := func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmplRoot.Execute(rw, nil)
			return
		}

		scraper := Scraper{url: r.FormValue("url"), semester: r.FormValue("semester")}
		scraper.createUrlOffset()
		lecturesList := scraper.getLectures()
		slices.SortFunc(lecturesList, compareLecsByDays)
		tmplChart.Execute(rw, &lecturesList)
	}
	mux.HandleFunc("/", rootRoute)
	mux.HandleFunc("/chart", chartRoute)
	http.ListenAndServe(":4567", mux)
}
