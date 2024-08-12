package controller

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"log"
	"net/http/httptest"
	"proxy/internal/modules/geo/controller/mock_service"
	"proxy/internal/modules/geo/entities"
	"proxy/internal/utils/readresponder"
	"testing"
)

func TestGeo_AddressSearch(t *testing.T) {
	testCases := []struct {
		name        string
		body        any
		wantStatus  int
		wantMessage string
	}{
		{"successful request", entities.AddressSearch{"улица Ленина"}, 200, "search completed"},
		{"empty query", nil, 400, "query is required"},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := NewMockService(controller)
	geo := NewGeo(mockService, readresponder.NewReadRespond())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body bytes.Buffer
			_ = json.NewEncoder(&body).Encode(tc.body)

			req := httptest.NewRequest("POST", "/api/address/search", &body)
			wr := httptest.NewRecorder()

			geo.AddressSearch(wr, req)

			r := wr.Result()

			var resp readresponder.JSONResponse
			err := json.NewDecoder(r.Body).Decode(&resp)
			if err != nil {
				log.Println(err)
			}

			defer r.Body.Close()

			if r.StatusCode != tc.wantStatus {
				t.Errorf("got status code %d, want %d", r.StatusCode, tc.wantStatus)
			}

			if tc.wantMessage != resp.Message {
				t.Errorf("got message %s, want %s", resp.Message, tc.wantMessage)

			}
		})
	}
}

func TestGeo_AddressGeocode(t *testing.T) {
	testCases := []struct {
		name        string
		body        any
		wantStatus  int
		wantMessage string
	}{
		{"successful request", entities.AddressGeocode{"5.12501", "1.15080"}, 200, "search completed"},
		{"no latitude", entities.AddressGeocode{"", "1.14080"}, 400, "both lat and lng are required"},
		{"no longitude", entities.AddressGeocode{"5.15080", ""}, 400, "both lat and lng are required"},
		{"empty body", nil, 400, "both lat and lng are required"},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := NewMockService(controller)
	geo := NewGeo(mockService, readresponder.NewReadRespond())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body bytes.Buffer
			_ = json.NewEncoder(&body).Encode(tc.body)

			req := httptest.NewRequest("POST", "/api/address/geocode", &body)
			wr := httptest.NewRecorder()

			geo.AddressGeocode(wr, req)

			r := wr.Result()

			var resp readresponder.JSONResponse
			err := json.NewDecoder(r.Body).Decode(&resp)
			if err != nil {
				log.Println(err)
			}

			defer r.Body.Close()

			if r.StatusCode != tc.wantStatus {
				t.Errorf("got status code %d, want %d", r.StatusCode, tc.wantStatus)
			}

			if tc.wantMessage != resp.Message {
				t.Errorf("got message %s, want %s", resp.Message, tc.wantMessage)

			}
		})
	}
}

func NewMockService(controller *gomock.Controller) *mock_service.MockGeoServicer {
	mockService := mock_service.NewMockGeoServicer(controller)

	mockService.EXPECT().AddressSearch(gomock.Any()).Return([]*entities.Address{}, nil).AnyTimes()
	mockService.EXPECT().GeoCode(gomock.Any(), gomock.Any()).Return([]*entities.Address{}, nil).AnyTimes()

	return mockService
}
