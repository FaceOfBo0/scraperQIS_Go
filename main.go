package main

import (
	"fmt"
	"net/http"
	"text/template"

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

func Tern[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func main() {

	tmpl_todo := template.Must(template.ParseFiles("templates/todos.html"))
	tmpl_form := template.Must(template.ParseFiles("templates/forms.html"))

	fmt.Println("hello Server")
	r := mux.NewRouter()

	r.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
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
	})

	r.HandleFunc("/todos", func(wrt http.ResponseWriter, req *http.Request) {
		dataTodo := TodoPageData{
			PageTitle: "My TODO list",
			Todos: []Todo{
				{Title: "Task 1", Done: false},
				{Title: "Task 2", Done: false},
				{Title: "Task 3", Done: false},
			},
		}

		tmpl_todo.Execute(wrt, dataTodo)
	})

	r.HandleFunc("/", func(wrt http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(wrt, "Hi you requested the following endpoint: %s\n", req.URL.Path)
	}).Methods("GET")

	r.HandleFunc("/hello", func(wrt http.ResponseWriter, req *http.Request) {
		userid := Tern(req.URL.Query().Has("id"), req.URL.Query().Get("id"), "null")
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
	r.Handle("/statics", http.StripPrefix("/statics", fs))
	http.ListenAndServe(":5000", r)
}
