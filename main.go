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
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
)

type Book struct {
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
	name := strings.TrimPrefix(url, "/api/books/")
	log.Println("GetBook name:", name)
	m := a.repo.findAllBooks()

	if b, ok := m[name]; ok {
		data, _ := json.Marshal(b)
		fmt.Fprintln(w, string(data))
	} else {
		if req.URL.Query().Has("fuzzy") {
			log.Println("fuzzy search")
			titles := []string{}
			for _, b := range m {
				titles = append(titles, b.Title)
			}
			log.Printf("titles:%#v \n", titles)

			// matches := fuzzy.Find(name, titles)
			matches := fuzzy.RankFindFold(name, titles)
			sort.Sort(matches)

			log.Println("fuzzy matches:", matches)

			if len(matches) > 0 {
				b := m[matches[0].Target]
				data, _ := json.Marshal(b)
				fmt.Fprintln(w, string(data))
				return
			}
		}
		fmt.Fprintln(w, "[]")

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
	version := flag.Bool("version", false, "prints app version and exits")
	flag.Parse()
	if *version {
		fmt.Println("bookline:", Version, "git:", GitRev)
		os.Exit(0)
	}

	r := NewInMemoryRepo()
	app := &App{
		repo: r,
		met:  NewMetrics(),
	}
	app.met.numOfBooks.Set(float64(len(app.repo.findAllBooks())))

	mux := http.NewServeMux()
	mux.HandleFunc("/api/books", app.handleListBook) // GET list books
	mux.HandleFunc("/api/books/", app.handleGetBook) // GET list books
	// mux.HandleFunc("/api/books/", app.handleApi) // POST create book
	// mux.HandleFunc("/api/books/", app.handleApi) // DELETE create book

	// mux.HandleFunc("/api/authors", app.handleApi) //list authors
	// mux.HandleFunc("/api/authors/ID", app.handleApi) //list authors

	mux.HandleFunc("/api/", app.handleApi)
	mux.HandleFunc("/version", app.handleVersion)
	mux.HandleFunc("/index.html", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/html;krumpli=9")
		fmt.Fprintln(w, "<h2>todo ...</h2>")
	})
	mux.Handle("/metrics", promhttp.Handler())

	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	http.ListenAndServe(":8888", std.Handler("", mdlw, mux))

}
