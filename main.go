package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
)

func init() {
	err := loadStoreMaster("StoreMasterAssignment.csv")
	if err != nil {
		log.Fatalf("Failed to load store master: %v", err)
	}
}

func loadStoreMaster(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("err", err)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Println("err", err)
		return err
	}
	for _, record := range records[1:] {
		storeMaster[record[2]] = Store{
			AreaCode:  record[0],
			StoreName: record[1],
			StoreID:   record[2],
		}
	}
	return nil
}

func main() {
	fmt.Println("starting the server")
	r := Router()
	fmt.Printf("listening at port 4000\n")

	log.Fatal(http.ListenAndServe(":4000", r))
}
