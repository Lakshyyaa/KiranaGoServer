package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

func jobSimulation(reqVar JobRequest, curJobId int) {
	var wg sync.WaitGroup
	var jobErrors []ErrorType
	var errorMu sync.Mutex
	for _, val := range reqVar.Visits {
		_, found := storeMaster[val.Store_ID]
		if !found {
			errorMu.Lock()
			fmt.Printf("store id not found\n")
			jobErrors = append(jobErrors, ErrorType{
				Store_ID: val.Store_ID,
				Error:    "store_id does not exists",
			})
			errorMu.Unlock()
			continue
		}
		for _, img := range val.Image_URL {
			wg.Add(1)
			go func(image string, storeID string) {
				defer wg.Done()
				err := imageSimulation(image)
				if err != nil {
					errorMu.Lock()
					jobErrors = append(jobErrors, ErrorType{
						Store_ID: storeID,
						Error:    "image processing failed",
					})
					errorMu.Unlock()
				}
			}(img, val.Store_ID)
		}
	}
	wg.Wait()

	// Update job status
	Statusmu.Lock()
	job, exists := jobs[curJobId]
	if exists {
		if len(jobErrors) > 0 {
			job.Status = "failed"
			job.Error = jobErrors
		} else {
			job.Status = "completed"
		}
	}
	Statusmu.Unlock()
}
func imageSimulation(img string) error {
	// Calculate perimeter (random for simulation)
	perimeter := 2 * (rand.Intn(100) + 100)
	// Simulate processing time between 0.1 and 0.4 seconds
	time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond)
	log.Printf("Processed image: %s, Perimeter: %d", img, perimeter)
	// Simulate random failures (10% chance) though not necessary
	if rand.Float32() < 0.1 {
		return fmt.Errorf("failed to process image")
	}
	return nil
}