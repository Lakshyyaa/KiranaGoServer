package main

import (
	"fmt"
	"net/http"
	"log"
)

func main() {
	fmt.Println("hello pls")
	r := Router()
	log.Fatal(http.ListenAndServe(":4000", r))
	fmt.Printf("listening at port 4000")
}