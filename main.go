package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func chuck() {
	resp, err := http.Get("https://api.chucknorris.io/jokes/random")
	if err != nil {
		log.Fatal(err)
	}
	// sc := bufio.NewScanner(resp.Body)
	// for sc.Scan() {
	// 	fmt.Println(sc.Text())
	// }
	data, _ := io.ReadAll(resp.Body)
	fmt.Println(string(data))
}

type Book struct {
	Title  string
	Author string
}

type App struct {
	repo Repository
}

var Version = "1.0.0"

func (a *App) handleApi(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "api todo ..")
}
func (a *App) handleVersion(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, Version)
}

func (a *App) handleListBook(w http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(a.repo.findAllBooks())
	fmt.Fprintln(w, string(data))
}

func (a *App) handleGetBook(w http.ResponseWriter, req *http.Request) {
	url := req.URL.Path
	name := strings.TrimPrefix(url, "/api/books/")
	m := a.repo.findAllBooks()
	b, ok := m[name]
	if ok {
		data, _ := json.Marshal(b)
		fmt.Fprintln(w, string(data))
	} else {
		fmt.Fprintln(w, "NOTFOUND")
	}

}

type Repository interface {
	findAllBooks() map[string]Book
}

type InMemoryRepo struct {
	books map[string]Book
}

func (r *InMemoryRepo) findAllBooks() map[string]Book {
	// mutex .....
	return r.books
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		books: map[string]Book{
			"War and Peace":  {"War and Peace", "Tolsztoj"},
			"Harry Potter I": {"Harry Potter I", "J.K."},
		},
	}
}

func main() {
	r := NewInMemoryRepo()
	app := &App{
		repo: r,
	}

	http.HandleFunc("/api/books", app.handleListBook) // GET list books
	http.HandleFunc("/api/books/", app.handleGetBook) // GET list books
	// http.HandleFunc("/api/books/", app.handleApi) // POST create book
	// http.HandleFunc("/api/books/", app.handleApi) // DELETE create book

	// http.HandleFunc("/api/authors", app.handleApi) //list authors
	// http.HandleFunc("/api/authors/ID", app.handleApi) //list authors

	http.HandleFunc("/api/", app.handleApi)
	http.HandleFunc("/version", app.handleVersion)
	http.HandleFunc("/index.html", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/html;krumpli=9")
		fmt.Fprintln(w, "<h2>todo ...</h2>")
	})
	http.ListenAndServe(":8888", nil)
}
