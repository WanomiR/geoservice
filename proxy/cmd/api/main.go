package main

import (
	"log"
	"proxy/internal/app"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	go a.Start()

	a.Shutdown()
}
