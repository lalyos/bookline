package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Book struct {
	//ID     uint `gorm:"primary_key"`
	gorm.Model
	Title  string
	Author string
}

type mymetrics struct {
	numOfBooks prometheus.Gauge
}

func NewMetrics() *mymetrics {
	m := &mymetrics{
		numOfBooks: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "books_number",
			Help: "Current number of Books.",
		}),
	}
	prometheus.MustRegister(m.numOfBooks)
	return m
}

type App struct {
	repo Repository
	met  *mymetrics
}

var Version = "0.1.0"
var GitRev = ""

func (a *App) handleApi(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "api todo ..")
}
func (a *App) handleVersion(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, Version)
}

func (a *App) handleListBook(w http.ResponseWriter, req *http.Request) {
	// fake change
	a.met.numOfBooks.Dec()
	books := a.repo.findAllBooks()
	list := []Book{}
	for _, b := range books {
		list = append(list, b)
	}
	data, _ := json.Marshal(list)
	fmt.Fprintln(w, string(data))
}

func (a *App) handleGetBook(w http.ResponseWriter, req *http.Request) {
	// fake change
	a.met.numOfBooks.Inc()
	url := req.URL.Path
	id := strings.TrimPrefix(url, "/api/books/")
	log.Println("GetBook ID:", id)
	m := a.repo.findAllBooks()

	var idx uint
	fmt.Sscanf(id, "%d", &idx)
	if b, ok := m[idx]; ok {
		data, _ := json.Marshal(b)
		fmt.Fprintln(w, string(data))
	} else {
		fmt.Fprintln(w, "[]")
	}
}

func (a *App) handleSearchByTitle(w http.ResponseWriter, req *http.Request) {
	url := req.URL.Path
	name := strings.TrimPrefix(url, "/api/search/")
	log.Println("GetBook name:", name)
	books := a.repo.findAllBooks()

	titles := []string{}
	byTitle := map[string]Book{}
	for _, b := range books {
		titles = append(titles, b.Title)
		byTitle[b.Title] = b
	}

	if b, ok := byTitle[name]; ok {
		data, _ := json.Marshal(b)
		fmt.Fprintln(w, string(data))
	} else {
		if req.URL.Query().Has("fuzzy") {
			log.Println("fuzzy search")

			// matches := fuzzy.Find(name, titles)
			matches := fuzzy.RankFindFold(name, titles)
			sort.Sort(matches)

			log.Println("fuzzy matches:", matches)

			if len(matches) > 0 {
				b := byTitle[matches[0].Target]
				data, _ := json.Marshal(b)
				fmt.Fprintln(w, string(data))
				return
			}
		}
		fmt.Fprintln(w, "[]")

	}
}

type DBRepo struct {
	DB *gorm.DB
}

func (r *DBRepo) ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&Book{})
	r.DB = database
}

func (r *DBRepo) findAllBooks() map[uint]Book {
	ret := map[uint]Book{}

	var books []Book
	r.DB.Find(&books)
	for _, b := range books {
		ret[b.ID] = b
	}

	return ret
}

type Repository interface {
	findAllBooks() map[uint]Book
}

type InMemoryRepo struct {
	books map[uint]Book
}

func (r *InMemoryRepo) findAllBooks() map[uint]Book {
	// mutex .....
	return r.books
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		books: map[uint]Book{
			1: {Title: "War and Peace", Author: "Tolsztoj"},
			2: {Title: "Harry Potter I", Author: "J.K."},
		},
	}
}

// simple middleware
func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)

		f(w, r)
	}
}

func main() {
	version := flag.Bool("version", false, "prints app version and exits")
	flag.Parse()
	if *version {
		fmt.Println("bookline:", Version, "git:", GitRev)
		os.Exit(0)
	}

	db := DBRepo{}
	db.ConnectDatabase()
	// db.DB.Create(&Book{Author: "Tolsztoj", Title: "War and Peace"})
	// db.DB.Create(&Book{Author: "J.K.", Title: "Harry Potter I"})
	// db.DB.Create(&Book{Author: "Jozsef", Title: "Bible"})

	// books := []Book{}
	// db.DB.Find(&books)
	// fmt.Println("books", books)
	// os.Exit(0)

	// r := NewInMemoryRepo()
	app := &App{
		// repo: r,
		repo: &db,
		met:  NewMetrics(),
	}
	app.met.numOfBooks.Set(float64(len(app.repo.findAllBooks())))

	mux := http.NewServeMux()
	mux.HandleFunc("/api/books", logging(app.handleListBook)) // GET list books
	mux.HandleFunc("/api/books/", logging(app.handleGetBook)) // GET list books
	// mux.HandleFunc("/api/books/", logging(app.handleApi)) // POST create book
	// mux.HandleFunc("/api/books/", logging(app.handleApi)) // DELETE create book

	// mux.HandleFunc("/api/authors", logging(app.handleApi)) //list authors
	// mux.HandleFunc("/api/authors/ID", alogging(pp.handleApi)) //list authors

	mux.HandleFunc("/api/search/", logging(app.handleSearchByTitle))
	mux.HandleFunc("/api/", logging(app.handleApi))
	mux.HandleFunc("/version", logging(app.handleVersion))
	mux.HandleFunc("/index.html", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/html;krumpli=9")
		fmt.Fprintln(w, "<h2>todo ...</h2>")
	})
	mux.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":8888", mux)
}
