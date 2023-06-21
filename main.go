package main

import (
	"github.com/adityarudrawar/app"
)

func main(){

	// TODO: Add support for env 
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}