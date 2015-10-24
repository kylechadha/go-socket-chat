package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {

	// Create a new mux Router.
	router := mux.NewRouter().StrictSlash(true)

	// WebSocket routes.
	router.HandleFunc("/ws", socketHandler)

	// Public routes.
	router.PathPrefix("/scripts").Handler(http.FileServer(http.Dir("./public/")))
	router.PathPrefix("/styles").Handler(http.FileServer(http.Dir("./public/")))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/views")))

	return router
}
