package main

import (
	_ "backend/docs"
	"backend/internal/app"
	"log"
)

// @title GoLibrary
// @version 1.0.0
// @description Library API

// @host localhost:8888
// @basePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	go a.Start()

	a.Stop()
}
