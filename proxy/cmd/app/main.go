package main

import (
	"log"
	"os"
	"proxy/internal/app"
)

// @title Microservice Geoservice
// @version 2.0.0
// @description Geoservice API

// @host localhost:8080
// @basePath /api

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(a.Run())
}
