package main

import (
	"bytes"
	"encoding/json"
	_ "geoservice/docs"
	"geoservice/internal/app"
	v1 "geoservice/internal/modules/geo/controller/http_v1"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"
)

// @title GeoService
// @version 1.0.0
// @description Geoservice API

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

	go simulateLoad()

	a.Shutdown()
}

func simulateLoad() {
	client := http.Client{}

	for {
		var geocode bytes.Buffer
		json.NewEncoder(&geocode).Encode(v1.RequestAddressGeocode{Lat: "55.753214", Lng: "37.642589"})
		req1, _ := http.NewRequest(http.MethodPost, "http://localhost:8888/address/geocode", &geocode)

		var address bytes.Buffer
		json.NewEncoder(&address).Encode(v1.RequestAddressSearch{Query: "Подкопаевский переулок"})
		req2, _ := http.NewRequest(http.MethodPost, "http://localhost:8888/address/search", &address)

		requests := []http.Request{*req1, *req2}

		reqId := rand.Intn(2)
		client.Do(&requests[reqId])

		sleepFactor := rand.Intn(11) * 15
		time.Sleep(time.Duration(sleepFactor) * time.Millisecond)
	}

}
