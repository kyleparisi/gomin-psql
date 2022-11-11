package main

import (
	"gomin/src/framework"
	"gomin/src/login"
	"gomin/src/register"
	"log"
	"net/http"
)

// Routes
func main() {
	router := framework.NewRouter()
	router.Get("/login", login.GetHandler)
	router.Post("/login", login.PostHandler)
	router.Get("/register", register.GetHandler)
	router.Post("/register", register.PostHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}
