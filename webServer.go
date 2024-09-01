package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"slices"
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

/* func Saving() Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {

			hf(rw, r)
		}
	}
} */

func Chain(hf http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		hf = m(hf)
	}
	return hf
}

func RunServer() {

	scr := Scraper{}
	// load template files
	tmplRoot := template.Must(template.ParseFiles("templates/root.html"))
	tmplChart := template.Must(template.ParseFiles("templates/chart.html"))

	//define route handler functions
	rootRoute := func(rw http.ResponseWriter, r *http.Request) {
		tmplRoot.Execute(rw, nil)
	}

	chartRoute := func(rw http.ResponseWriter, r *http.Request) {
		/* 		if r.Method != http.MethodPost {
			tmplRoot.Execute(rw, nil)
			return
		} */
		scr.Url = r.FormValue("url")
		scr.Semester = r.FormValue("semester")
		scr.createUrlOffset()
		scr.loadLectures()
		slices.SortFunc(scr.lectures, compareLecsByDays)
		lecsJson, err := json.Marshal(&scr.lectures)
		if err != nil {
			fmt.Printf("error with json serialization: %v\n", err)
			return
		}

		err = saveJsonToFile(string(lecsJson), "lectures.json")
		if err != nil {
			fmt.Printf("error saving JSON to file: %v\n", err)
			return
		}

		tmplChart.Execute(rw, struct {
			Lectures *[]Lecture
			Semester string
		}{Lectures: &scr.lectures, Semester: scr.lectures[0].Semester})
	}

	downloadRoute := func(rw http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("java", "-jar", "lfw.jar", "lectures.json")
		if err := cmd.Run(); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		rw.WriteHeader(http.StatusOK)
	}

	// Bind handler functions and start server
	fmt.Println("Server is running on localhost:4567..")
	http.HandleFunc("/", rootRoute)
	http.HandleFunc("/chart", Chain(chartRoute, Logging()))
	http.HandleFunc("/download", Chain(downloadRoute, Logging()))
	http.ListenAndServe(":4567", nil)

}
