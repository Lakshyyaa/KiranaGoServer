package main

import (
	"log"
	"math/rand"
	"time"
)

func jobSimulation(reqVar JobRequest, curJobId int) {
	for _, val := range reqVar.Visits {
		_, found := storeMaster[val.Store_ID]
		if !found {
			log.Println("Invalid store ID: ", val.Store_ID)
			continue
		}
		for _, img := range val.Image_URL {
			go imageSimulation(img)
		}
	}
}

func imageSimulation(img string) {
	// fetch the image and find its perimeter
	perimeter := 2 * (rand.Intn(100) + 100)
	// simulating download and gpu processing
	time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond)
	log.Print(perimeter, img)
	// storing this in future
}
