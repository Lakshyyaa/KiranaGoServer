package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Visit struct {
	Store_ID   string   `json:"store_id"`
	Image_URL  []string `json:"image_url"`
	Visit_Time string   `json:"visit_time"`
}
type JobRequest struct {
	Count  int   `json:"count"`
	Visits []Visit `json:"visits"`
}

func submitJobHandler(w http.ResponseWriter, r *http.Request) {
	var reqVar JobRequest
	err := json.NewDecoder(r.Body).Decode(&reqVar)
	if err != nil {
		http.Error(w, "Invalid payload request", http.StatusBadRequest)
		return
	}
	if reqVar.Count != len(reqVar.Visits) {
		http.Error(w, "Wrong count sent with the request", http.StatusBadRequest)
		return
	}	
	go jobSimulation(reqVar)
}

func getJobInfoHandler(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("jobid")
	if jobID == "" {
		log.Println("No jobID")
		return
	}
	fmt.Println(jobID, "is ok")
}
