package main

import (
	_ "backend/docs"
	"backend/internal/app"
	"log"
	_ "net/http/pprof"
)

// @title GeoService
// @version 1.0.0
// @description Geoservice API

// @host localhost:8888
// @basePath /
func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	go a.Start()

	a.Shutdown()
}
