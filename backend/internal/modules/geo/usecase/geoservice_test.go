package usecase

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

var (
	apiKey    string
	secretKey string
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey = os.Getenv("DADATA_API_KEY")
	secretKey = os.Getenv("DADATA_SECRET_KEY")
}

func TestGeoService_AddressSearch(t *testing.T) {
	testCases := []struct {
		name       string
		input      string
		wantErr    bool
		wantResult bool
	}{
		{"none-sense request", "adlgjaogu", false, false},
		{"adequate request", "улица Ленина", false, true},
	}

	g := NewGeoService(apiKey, secretKey)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := g.AddressSearch(tc.input)

			if (err != nil) != tc.wantErr {
				t.Errorf("AddressSearch() error = %v, wantErr %v", err, tc.wantErr)
			}

			if (len(res) > 0) != tc.wantResult {
				t.Errorf("AddressSearch() got = %v, want %v", res, tc.wantResult)
			}
		})
	}
}

func TestGeoService_GeoCode(t *testing.T) {
	type inputCoords struct {
		lat, lng string
	}
	testCases := []struct {
		name       string
		input      inputCoords
		wantErr    bool
		wantResult bool
	}{
		{"none-sense request", inputCoords{"-0.000", "+0.000"}, false, false},
		{"adequate request", inputCoords{"55.753214", "37.642589"}, false, true},
	}

	g := NewGeoService(apiKey, secretKey)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := g.GeoCode(tc.input.lat, tc.input.lng)

			if (err != nil) != tc.wantErr {
				t.Errorf("AddressSearch() error = %v, wantErr %v", err, tc.wantErr)
			}

			if (len(res) > 0) != tc.wantResult {
				t.Errorf("AddressSearch() got = %v, want %v", res, tc.wantResult)
			}
		})
	}
}
