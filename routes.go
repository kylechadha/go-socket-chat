package main

import "github.com/julienschmidt/httprouter"

func NewRouter() *httprouter.Router {

	// Create a new httprouter Router.
	router := httprouter.New()

	// Define the routes.
	// router.GET("/logs", utils.AppHandler(lc.HomeHandler))
	// router.POST("/logs", utils.AppHandler(lc.LogHandler))

	return router
}
