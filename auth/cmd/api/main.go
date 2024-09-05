package main

import (
	"auth/internal/app"
	"log"
	"os"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	go a.ServeMetrics()

	os.Exit(a.Run())
}
