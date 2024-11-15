package main

import (
	"fmt"
	"log"
	"net/http"
)

func submitJobHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %s", r.Method)
}

func getJobInfoHandler(w http.ResponseWriter, r *http.Request) {
    jobID := r.URL.Query().Get("jobid") 
    if jobID == "" {
        log.Println("No jobID")
        return
    }
    fmt.Println(jobID, "is ok")
}
