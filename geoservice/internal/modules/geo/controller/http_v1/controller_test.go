package http_v1

import (
	"bytes"
	"encoding/json"
	"geoservice/internal/modules/geo/dto"
	"github.com/wanomir/rr"
	"net/http"
	"net/http/httptest"
	"testing"
)

var geoController *GeoController

func init() {
	geoController = NewGeoController(NewMockGeoService(), rr.NewReadResponder())
}

func TestGeoController_AddressSearch(t *testing.T) {
	testCases := []struct {
		name     string
		payload  RequestAddressSearch
		wantCode int
	}{
		{"normal request", RequestAddressSearch{"Подкопаевский переулок"}, 200},
		{"incomplete data", RequestAddressSearch{""}, 400},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body bytes.Buffer
			json.NewEncoder(&body).Encode(tc.payload)

			req := httptest.NewRequest(http.MethodGet, "/address/search", &body)
			wr := httptest.NewRecorder()

			geoController.AddressSearch(wr, req)

			r := wr.Result()

			if r.StatusCode != tc.wantCode {
				t.Errorf("got %d, want %d", r.StatusCode, tc.wantCode)
			}
		})
	}
}

func TestGeoController_AddressGeocode(t *testing.T) {
	testCases := []struct {
		name     string
		payload  RequestAddressGeocode
		wantCode int
	}{
		{"normal request", RequestAddressGeocode{"55.753214", "37.642589"}, 200},
		{"incomplete data", RequestAddressGeocode{"", "37.642589"}, 400},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body bytes.Buffer
			json.NewEncoder(&body).Encode(tc.payload)

			req := httptest.NewRequest(http.MethodGet, "/address/geocode", &body)
			wr := httptest.NewRecorder()

			geoController.AddressGeocode(wr, req)

			r := wr.Result()

			if r.StatusCode != tc.wantCode {
				t.Errorf("got %d, want %d", r.StatusCode, tc.wantCode)
			}
		})
	}
}

type MockGeoService struct{}

func NewMockGeoService() *MockGeoService {
	return &MockGeoService{}
}

func (m *MockGeoService) AddressSearch(_ string) ([]dto.Address, error) {
	return []dto.Address{}, nil
}

func (m *MockGeoService) GeoCode(_, _ string) ([]dto.Address, error) {
	return []dto.Address{}, nil
}
