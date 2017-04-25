package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"log"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

func init() {
	// User store initialization
	store, err := NewFileUserStore("./data/users.json")
	if err != nil {
		panic(fmt.Errorf("Error creating user store: %s", err))
	}
	globalUserStore = store

	// Session store initialization
	sessionStore, err := NewFileSessionStore("./data/sessions.json")
	if err != nil {
		panic(fmt.Errorf("Error creating session store: %s", err))
	}
	globalSessionStore = sessionStore

	// DB connection initialization
	var dbCredentials struct{
		Username string
		Password string
	}

	file, err := ioutil.ReadFile("configuration/DBCredentials.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &dbCredentials)
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/gophr", dbCredentials.Username, dbCredentials.Password)
	db, err := NewMySQLDB(dsn)
	if err != nil {
		panic(err)
	}
	globalMySQLDB = db

	// Image store assignment
	globalImageStore = NewDBImageStore()
}

func main() {
	router := NewRouter()

	router.Handle("GET", "/", HandleHome)
	router.Handle("GET", "/register", HandleUserNew)
	router.Handle("POST", "/register", HandleUserCreate)
	router.Handle("GET", "/login", HandleSessionNew)
	router.Handle("POST", "/login", HandleSessionCreate)
	router.Handle("GET", "/image/:imageID", HandleImageShow)
	router.Handle("GET", "/user	/:userID", HandleUserShow)

	router.ServeFiles(
		"/assets/*filepath",
		http.Dir("assets/"),
	)

	router.ServeFiles(
		"/im/*filepath",
		http.Dir("data/images/"),
	)

	secureRouter := NewRouter()
	secureRouter.Handle("GET", "/sign-out", HandleSessionDestroy)
	secureRouter.Handle("GET", "/account", HandleUserEdit)
	secureRouter.Handle("POST", "/account", HandleUserUpdate)
	secureRouter.Handle("GET", "/images/new", HandleImageNew)
	secureRouter.Handle("POST", "/images/new", HandleImageCreate)

	middleware := Middleware{}
	middleware.Add(router)
	middleware.Add(http.HandlerFunc(RequireLogin))
	middleware.Add(secureRouter)

	log.Fatal(http.ListenAndServe(":3000", middleware))
}

// Creates a new router
func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	return router
}
