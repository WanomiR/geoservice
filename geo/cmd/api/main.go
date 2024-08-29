package main

import (
	"geo/internal/app"
	"log"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	go a.Start()

	go a.ServeMetrics()

	a.Shutdown()
}
