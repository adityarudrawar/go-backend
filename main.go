package main

import "github.com/adityarudrawar/go-backend/app"

func main(){

	// TODO: Add support for env 
	err := app.SetupAndRunApp()
	if err != nil {
		panic(err)
	}
}