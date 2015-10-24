package main

import (
	"net/http"
	"os"
)

func main() {

	// Router
	// ----------------------------
	router := newRouter()

	// Server
	// ----------------------------
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	http.ListenAndServe(":"+port, router)

}
