package main

import "net/http"

func main() {

	// Router
	// ----------------------------
	router := NewRouter()

	// Server
	// ----------------------------
	http.ListenAndServe(":8080", router)

}
