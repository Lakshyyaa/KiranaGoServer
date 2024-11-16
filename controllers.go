package main

import (
	"encoding/json"
	"net/http"
	"strconv"
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

type JobStateOK struct {
	Status string `json:"status"`
	Job_ID int    `json:"job_id"`
}
type JobResponseError struct {
	Error string `json:"error"`
}

var jobs = make(map[int]*JobState)
var storeMaster = make(map[string]Store)
var Statusmu sync.RWMutex
var jobCount = 0

func submitJobHandler(w http.ResponseWriter, r *http.Request) {
	var reqVar JobRequest
	err := json.NewDecoder(r.Body).Decode(&reqVar)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JobResponseError{Error: "Invalid request format"})
		return
	}
	// Check for missing or invalid fields
	if reqVar.Count <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JobResponseError{Error: "Count must be greater than 0"})
		return
	}

	if len(reqVar.Visits) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JobResponseError{Error: "Visits array is empty"})
		return
	}

	// Validate each visit
	for _, visit := range reqVar.Visits {
		if visit.Store_ID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(JobResponseError{Error: "Missing store_id"})
			return
		}

		if len(visit.Image_URL) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(JobResponseError{Error: "Missing image_url"})
			return
		}

		// Check each image URL in the array
		for _, url := range visit.Image_URL {
			if url == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(JobResponseError{Error: "Empty image URL "})
				return
			}
		}

		if visit.Visit_Time == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(JobResponseError{Error: "Missing visit_time"})
			return
		}
	}

	if reqVar.Count != len(reqVar.Visits) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JobResponseError{Error: "Count does not match number of visits"})
		return
	}

	// adding lock on the critical section: jobCount
	Statusmu.Lock()
	jobCount++
	curJobId := jobCount
	jobs[curJobId] = &JobState{
		Status: "ongoing",
		Job_ID: curJobId,
		Error:  make([]ErrorType, 0),
	}
	Statusmu.Unlock()
	go jobSimulation(reqVar, curJobId)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JobResponseOK{JobID: curJobId})
}

func getJobInfoHandler(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("jobid")
	w.Header().Set("Content-Type", "application/json")
	jobIDInt, err := strconv.Atoi(jobID)
	if jobID == "" || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{}{})
		return
	}
	Statusmu.RLock()
	job, found := jobs[jobIDInt]
	Statusmu.RUnlock()
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(struct{}{})
		return
	}
	if len(job.Error) > 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(job)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(JobStateOK{
			Status: job.Status,
			Job_ID: jobIDInt,
		})
		return
	}
}
