package rpc_v1

import (
	"errors"
	"fmt"
	"geoprovider/internal/dto"
	"github.com/brianvoe/gofakeit/v7"
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
			if err := geoController.AddressSearch(tc.input, &reply); (err != nil) != tc.wantErr {
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
			if err := geoController.GeoCode(tc.args, &reply); (err != nil) != tc.wantErr {
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
	addresses := make([]dto.Address, 0)

	for i := 0; i < 10; i++ {
		fakeAddress := gofakeit.Address()
		address := dto.Address{
			City:   fakeAddress.City,
			Street: fakeAddress.Street,
			House:  fakeAddress.Address,
			Lat:    fmt.Sprintf("%f", fakeAddress.Latitude),
			Lon:    fmt.Sprintf("%f", fakeAddress.Longitude),
		}
		addresses = append(addresses, address)
	}

	if input == "" {
		return []dto.Address{}, errors.New("no address provided")
	}
	return addresses, nil
}

func (m *MockGeoProvider) GeoCode(lat, lng string) ([]dto.Address, error) {
	addresses := make([]dto.Address, 0)

	if lat == "" || lng == "" {
		return addresses, errors.New("bad coordinates")
	}

	for i := 0; i < 10; i++ {
		fakeAddress := gofakeit.Address()
		address := dto.Address{
			City:   fakeAddress.City,
			Street: fakeAddress.Street,
			House:  fakeAddress.Address,
			Lat:    fmt.Sprintf("%f", fakeAddress.Latitude),
			Lon:    fmt.Sprintf("%f", fakeAddress.Longitude),
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}
