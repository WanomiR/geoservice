package jsonrpc_v1

import (
	"errors"
	"geoprovider/internal/dto"
	"testing"
)

var geoController *GeoController

func init() {
	geoController = NewGeoController(NewMockGeoProvider())
}

func TestGeoController_AddressSearch(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"normal case", "Улица Ленина", false},
		{"empty input", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var reply dto.Addresses
			if err := geoController.AddressSearch(&AddressSearchArgs{Address: tc.input}, &reply); (err != nil) != tc.wantErr {
				t.Errorf("AddressSearch() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestGeoController_GeoCode(t *testing.T) {
	testCases := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"normal case", []string{"55.098765", "37.97259"}, false},
		{"bad coordinates", []string{"55.098765", ""}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var reply dto.Addresses
			if err := geoController.GeoCode(&GeoCodeArgs{Lat: tc.args[0], Lng: tc.args[1]}, &reply); (err != nil) != tc.wantErr {
				t.Errorf("GeoCode() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

type MockGeoProvider struct{}

func NewMockGeoProvider() *MockGeoProvider {
	return &MockGeoProvider{}
}

func (m *MockGeoProvider) AddressSearch(input string) ([]dto.Address, error) {
	if input == "" {
		return []dto.Address{}, errors.New("no address provided")
	}
	return []dto.Address{}, nil
}

func (m *MockGeoProvider) GeoCode(lat, lng string) ([]dto.Address, error) {
	if lat == "" || lng == "" {
		return []dto.Address{}, errors.New("bad coordinates")
	}
	return []dto.Address{}, nil
}
