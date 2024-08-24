package main

import (
	"geoprovider/internal/app"
	"log"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	go a.Start()

	a.Shutdown()
}
