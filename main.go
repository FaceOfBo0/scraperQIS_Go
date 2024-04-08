package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
}

type ContactDetails struct {
	Email   string
	Subject string
	Message string
}

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Method(meth string) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != meth {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			hf(w, r)
		}
	}
}

func Logging() Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() { log.Println(r.URL.Path, time.Since(start)) }()

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

func Tern[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

// func logMiddleware(hf http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println(r.URL.Path)
// 		hf.ServeHTTP(w, r)
// 	})
// }

// func logging(hf http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Println(r.URL.Path)
// 		hf(w, r)
// 	}
// }

func main() {

	tmpl_todo := template.Must(template.ParseFiles("templates/todos.html"))
	tmpl_form := template.Must(template.ParseFiles("templates/forms.html"))

	formRoute := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl_form.Execute(w, nil)
			return
		}

		details := ContactDetails{
			Email:   r.FormValue("email"),
			Subject: r.FormValue("subject"),
			Message: r.FormValue("message"),
		}

		fmt.Println(details)

		tmpl_form.Execute(w, struct{ Success bool }{true})
	}

	todoRoute := func(wrt http.ResponseWriter, req *http.Request) {
		dataTodo := TodoPageData{
			PageTitle: "My TODO list",
			Todos: []Todo{
				{Title: "Task 1", Done: false},
				{Title: "Task 2", Done: false},
				{Title: "Task 3", Done: false},
			},
		}

		tmpl_todo.Execute(wrt, dataTodo)
	}

	fmt.Println("hello Server")
	r := mux.NewRouter()

	r.HandleFunc("/form", Chain(formRoute, Method("GET"), Logging()))
	r.HandleFunc("/todos", Chain(todoRoute, Method("GET"), Logging()))

	r.HandleFunc("/", func(wrt http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(wrt, "Hi you requested the following endpoint: %s\n", req.URL.Path)
	}).Methods("GET")

	r.HandleFunc("/hello", func(wrt http.ResponseWriter, req *http.Request) {
		//userid := Tern(req.URL.Query().Has("id"), req.URL.Query().Get("id"), "null")
		//println(userid)
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
	r.Handle("/statics", http.StripPrefix("/statics", fs))
	//r.Use(logMiddleware)
	http.ListenAndServe(":5000", r)
}
