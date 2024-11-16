package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Visit struct {
	Store_ID   string   `json:"store_id"`
	Image_URL  []string `json:"image_url"`
	Visit_Time string   `json:"visit_time"`
}
type JobRequest struct {
	Count  int     `json:"count"`
	Visits []Visit `json:"visits"`
}
type Store struct {
	AreaCode  string
	StoreName string
	StoreID   string
}
type ErrorType struct {
	Store_ID string `json:"store_id"`
	Error    string `json:"error"`
}
type JobState struct {
	Status string      `json:"status"`
	Job_ID int         `json:"job_id"`
	Error  []ErrorType `json:"error"`
}
type JobResponseOK struct {
	JobID int `json:"job_id"`
}

type JobResponseError struct {
	Error string `json:"error"`
}

var jobs = make(map[int]*JobState)
var storeMaster = make(map[string]Store)
var mu sync.RWMutex
var jobCount = 0

func submitJobHandler(w http.ResponseWriter, r *http.Request) {
	var reqVar JobRequest
	err := json.NewDecoder(r.Body).Decode(&reqVar)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JobResponseError{Error: "Invalid request format"})
		return
	}
	if reqVar.Count != len(reqVar.Visits) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JobResponseError{Error: "Count does not match number of visits"})
		return
	}
	// adding lock on the critical section: jobCount
	mu.Lock()
	jobCount++
	curJobId := jobCount
	jobs[curJobId] = &JobState{
		Status: "ongoing",
		Job_ID: curJobId,
	}
	mu.Unlock()
	go jobSimulation(reqVar, curJobId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JobResponseOK{JobID: curJobId})
}

func getJobInfoHandler(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("jobid")
	if jobID == "" {
		log.Println("No jobID")
		return
	}
	fmt.Println(jobID, "is ok")
}