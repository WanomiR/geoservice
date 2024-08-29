package grpc_v1

import (
	"context"
	"errors"
	"fmt"
	"geo/internal/dto"
	pb "geo/pkg/geo_v1"
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
			req := &pb.AddressRequest{Query: tc.input}
			if _, err := geoController.AddressSearch(context.Background(), req); (err != nil) != tc.wantErr {
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
			req := &pb.GeoRequest{Lat: tc.args[0], Lng: tc.args[1]}
			if _, err := geoController.GeoCode(context.Background(), req); (err != nil) != tc.wantErr {
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

	if input == "" {
		return addresses, errors.New("no address provided")
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
