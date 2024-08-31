package geo

import (
	"errors"
	"proxy/internal/dto"
	"testing"
)

var usecase *GeoUsecase

func init() {
	usecase = NewGeoUsecase(NewMockProvider())
}

func TestGeoService_AddressSearch(t *testing.T) {
	testCases := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{"normal request", "Улица Ленина", false},
		{"incomplete query", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := usecase.AddressSearch(tc.query); (err != nil) != tc.wantErr {
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
		{"normal request", "55.12085", "37.10850", false},
		{"incomplete query", "", "35.79072", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := usecase.GeoCode(tc.lat, tc.lng); (err != nil) != tc.wantErr {
				t.Errorf("GeoCode() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (m *MockProvider) AddressSearch(query string) ([]dto.Address, error) {
	if query == "" {
		return []dto.Address{}, errors.New("invalid query")
	}
	return []dto.Address{}, nil
}

func (m *MockProvider) GeoCode(lat, lng string) ([]dto.Address, error) {
	if lat == "" || lng == "" {
		return []dto.Address{}, errors.New("invalid query")
	}
	return []dto.Address{}, nil
}
