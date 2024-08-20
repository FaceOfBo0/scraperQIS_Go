package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
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

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func authMiddleware() Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "cookie-name")
			if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			hf(w, r)
		}
	}
}

func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	fmt.Fprintln(w, "The cake is a lie!")
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	//Authentication methods...

	session.Values["authenticated"] = true
	sessions.Save(r, w)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	//Revoke authentication
	session.Values["authenticated"] = false
	session.Save(r, w)

}

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

func RunTutServer() {

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

	r.HandleFunc("/logout", Chain(logout, Logging()))
	r.HandleFunc("/login", Chain(login, Logging()))
	r.HandleFunc("/secret", Chain(secret, Logging()))
	r.HandleFunc("/form", Chain(formRoute, authMiddleware(), Logging()))
	r.HandleFunc("/todos", Chain(todoRoute, Method("GET"), authMiddleware(), Logging()))

	r.HandleFunc("/", func(wrt http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(wrt, "Hi you requested the following endpoint: %s\n", req.URL.Path)
	}).Methods("GET")

	r.HandleFunc("/hello", func(wrt http.ResponseWriter, req *http.Request) {
		// userid := Tern(req.URL.Query().Has("id"), req.URL.Query().Get("id"), "null")
		// println(userid)
		wrt.Header().Set("Content-Type", "text/html")
		wrt.Write([]byte("<html><h2>Hello</h2></html>"))
	})

	r.HandleFunc("/books/{title}/page/{page}", func(wrt http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(wrt, "You've requested the book: %s on page %s\n", title, page)
	})

	r.HandleFunc("/decode", func(w http.ResponseWriter, r *http.Request) {
		var book Book
		json.NewDecoder(r.Body).Decode(&book)

		fmt.Fprintf(w, "%s is written by %s with %d pages.\n", book.Title, book.Author, book.Pages)
	}).Methods("POST")

	r.HandleFunc("/encode", func(w http.ResponseWriter, r *http.Request) {
		book := Book{
			Author: "Goethe",
			Title:  "Faust",
			Pages:  683,
		}

		json.NewEncoder(w).Encode(book)
	})

	fs := http.FileServer(http.Dir("static/"))
	r.Handle("/statics", http.StripPrefix("/statics", fs))
	http.ListenAndServe(":5000", r)
}
