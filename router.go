package main

import (
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/submit", submitJobHandler).Methods("POST")
	router.HandleFunc("/api/status", getJobInfoHandler).Methods("GET")
	return router
}