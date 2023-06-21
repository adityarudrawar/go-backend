package main

import "github.com/adityarudrawar/go-backend/app"

func main(){
	err := app.SetupAndRunApp()
	if err != nil {
		panic(err)
	}
}