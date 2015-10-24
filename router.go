package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {

	// Create a new mux Router.
	router := mux.NewRouter().StrictSlash(true)

	// Public routes.
	router.PathPrefix("/scripts").Handler(http.FileServer(http.Dir("./public/")))
	router.PathPrefix("/styles").Handler(http.FileServer(http.Dir("./public/")))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/views")))

	// WebSocket routes.

	// Catch all route.
	// router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintln(w, "Oh man, what did you do? Folks we got a 404.")
	// })

	return router
}
