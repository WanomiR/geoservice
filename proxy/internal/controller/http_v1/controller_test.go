package http_v1

import (
	"bytes"
	"encoding/json"
	"github.com/wanomir/rr"
	"net/http"
	"net/http/httptest"
	"proxy/internal/dto"
	"testing"
)

var cntrl *Controller

func init() {
	cntrl = NewController(NewMockUseCase(), rr.NewReadResponder())
}

func TestController_AddressSearch(t *testing.T) {
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

			cntrl.AddressSearch(wr, req)

			r := wr.Result()

			if r.StatusCode != tc.wantCode {
				t.Errorf("got %d, want %d", r.StatusCode, tc.wantCode)
			}
		})
	}
}

func TestController_AddressGeocode(t *testing.T) {
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

			cntrl.AddressGeocode(wr, req)

			r := wr.Result()

			if r.StatusCode != tc.wantCode {
				t.Errorf("got %d, want %d", r.StatusCode, tc.wantCode)
			}
		})
	}
}

type MockUseCase struct{}

func NewMockUseCase() *MockUseCase {
	return &MockUseCase{}
}

func (m *MockUseCase) AddressSearch(_ string) ([]dto.Address, error) {
	return []dto.Address{}, nil
}

func (m *MockUseCase) GeoCode(_, _ string) ([]dto.Address, error) {
	return []dto.Address{}, nil
}
