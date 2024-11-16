package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
)

func init() {
	filename := "StoreMasterAssignment.csv"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error loading store master:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading CSV:", err)
	}

	for _, record := range records[1:] {
		storeMaster[record[0]] = Store{
			StoreID:   record[0],
			StoreName: record[1],
			AreaCode:  record[2],
		}
	}
}

func main() {
	fmt.Println("hello pls")
	r := Router()
	fmt.Printf("listening at port 4000\n")

	log.Fatal(http.ListenAndServe(":4000", r))
}
