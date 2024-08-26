package usecase

import (
	"log"
	"os"
	"testing"
)

var geoService *GeoService

func init() {
	apiKey := os.Getenv("DADATA_API_KEY")
	secretKey := os.Getenv("DADATA_SECRET_KEY")

	if apiKey == "" || secretKey == "" {
		log.Fatal("dadata env variables not set")
	}

	geoService = NewGeoService(apiKey, secretKey)
}

func TestGeoService_AddressSearch(t *testing.T) {
	testCases := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{"normal case", "Улица Ленина", false},
		{"empty query", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := geoService.AddressSearch(tc.query); (err != nil) != tc.wantErr {
				t.Errorf("AddressSearch() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestGeoService_GeoCode(t *testing.T) {
	testCases := []struct {
		name     string
		lat, lng string
		wantErr  bool
	}{
		{"normal case", "55.67195", "37.08907", false},
		{"incomplete coordinates", "", "37.08907", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := geoService.GeoCode(tc.lat, tc.lng); (err != nil) != tc.wantErr {
				t.Errorf("GeoCode() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
