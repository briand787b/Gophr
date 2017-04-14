package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"log"
)

func main() {
	router := NewRouter()

	router.Handle("GET", "/", HandleHome)

	router.ServeFiles(
		"/assets/*filepath",
		http.Dir("assets/"),
	)

	middleware := Middleware{}
	middleware.Add(router)
	log.Fatal(http.ListenAndServe(":3000", middleware))
}


func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){})
	return router
}