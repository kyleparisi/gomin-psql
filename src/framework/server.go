package framework

import (
	"database/sql"
	"fmt"
	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/ucarion/urlpath"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"
)

type Router struct {
	GET  []map[string]Handler
	POST []map[string]PostHandler
}
type Redirect struct {
	location string
	code     int
}
type Response struct {
	StatusCode int
	Data       interface{}
	Template   *template.Template
	Redirect   interface{}
	Session    *sessions.Session
}
type Handler func(*sql.Tx, *sessions.Session) func(urlpath.Match) Response
type PostHandler func(*sql.Tx, *sessions.Session) func(urlpath.Match, io.Reader) Response

func NewRouter() *Router {
	router := new(Router)
	router.GET = make([]map[string]Handler, 0)
	router.POST = make([]map[string]PostHandler, 0)
	return router
}

func NewDatabaseConnection() *sql.DB {
	dsn := fmt.Sprintf("sslmode=disable ")
	username := os.Getenv("DB_USERNAME")
	if username == "" {
		username = "app"
	}
	dsn += "user=" + username + " "
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "secret"
	}
	dsn += "password=" + password + " "
	database := os.Getenv("DB_DATABASE")
	if database == "" {
		database = "app_test"
	}
	dsn += "database=" + database + " "
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func NewSessionStore() *pgstore.PGStore {
	dsn := fmt.Sprintf("sslmode=disable ")
	username := os.Getenv("DB_USERNAME")
	if username == "" {
		username = "tron_local"
	}
	dsn += "user=" + username + " "
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "secret"
	}
	dsn += "password=" + password + " "
	database := os.Getenv("DB_DATABASE")
	if database == "" {
		database = "tron_local"
	}
	dsn += "database=" + database + " "
	store, err := pgstore.NewPGStore(dsn, []byte(os.Getenv("APP_SECRET_KEY")))
	if err != nil {
		panic(err.Error())
	}
	return store
}

func (r *Router) Get(path string, f Handler) {
	r.GET = append(r.GET, map[string]Handler{path: f})
}

func (r *Router) Post(path string, f PostHandler) {
	r.POST = append(r.POST, map[string]PostHandler{path: f})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method, req.URL.Path)
	db := NewDatabaseConnection()
	defer db.Close()
	store := NewSessionStore()
	defer store.StopCleanup(store.Cleanup(time.Minute * 5))
	defer store.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	session, _ := store.Get(req, "user")

	if req.Method == "GET" {
		for _, element := range r.GET {
			for path, handler := range element {
				route := urlpath.New(path)
				match, ok := route.Match(req.URL.Path)
				if !ok {
					break
				}
				response := handler(tx, session)(match)
				// check for redirect first
				if redirect, ok := response.Redirect.(Redirect); ok {
					http.Redirect(w, req, redirect.location, redirect.code)
					goto Done
				}
				response.Template.Execute(w, response.Data)
				goto Done
			}
		}
	}

	if req.Method == "POST" {
		for _, element := range r.POST {
			for path, post := range element {
				route := urlpath.New(path)
				match, ok := route.Match(req.URL.Path)
				if !ok {
					break
				}
				response := post(tx, session)(match, req.Body)
				if response.Session != nil {
					err := response.Session.Save(req, w)
					if err != nil {
						fmt.Println("problem saving session", err.Error())
					}
				}
				// check for redirect first
				if redirect, ok := response.Redirect.(Redirect); ok {
					http.Redirect(w, req, redirect.location, redirect.code)
					goto Done
				}
				if response.Template != nil {
					err := response.Template.Execute(w, response.Data)
					if err != nil {
						fmt.Println(err)
					}
				}
				goto Done
			}
		}
	}

	// 404
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
Done:
}
