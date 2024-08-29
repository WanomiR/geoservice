package grpc_v1

import (
	"context"
	"geoprovider/internal/dto"
	"geoprovider/pkg/geoprovider_v1"
	"github.com/wanomir/e"
)

type GeoProvider interface {
	AddressSearch(query string) ([]dto.Address, error)
	GeoCode(lat, lng string) ([]dto.Address, error)
}

type GeoController struct {
	geoprovider_v1.UnimplementedGeoProviderV1Server // safety fallback for non-implemented controllers
	service                                         GeoProvider
}

func NewGeoController(service GeoProvider) *GeoController {
	return &GeoController{service: service}
}

func (c *GeoController) AddressSearch(_ context.Context, in *geoprovider_v1.AddressRequest) (*geoprovider_v1.AddressesResponse, error) {
	addresses, err := c.service.AddressSearch(in.GetQuery())
	if err != nil {
		return nil, e.Wrap("error fetching addresses", err)
	}

	response := &geoprovider_v1.AddressesResponse{
		Addresses: make([]*geoprovider_v1.AddressResponse, 0),
	}

	for _, address := range addresses {
		addressResponse := &geoprovider_v1.AddressResponse{
			City:   address.City,
			Street: address.Street,
			House:  address.House,
			Lat:    address.Lat,
			Lon:    address.Lon,
		}
		response.Addresses = append(response.Addresses, addressResponse)
	}

	return response, nil
}

func (c *GeoController) GeoCode(_ context.Context, in *geoprovider_v1.GeoRequest) (*geoprovider_v1.AddressesResponse, error) {
	addresses, err := c.service.GeoCode(in.GetLat(), in.GetLng())
	if err != nil {
		return nil, e.Wrap("error fetching address", err)
	}

	response := &geoprovider_v1.AddressesResponse{
		Addresses: make([]*geoprovider_v1.AddressResponse, 0),
	}

	for _, address := range addresses {
		addressResponse := &geoprovider_v1.AddressResponse{
			City:   address.City,
			Street: address.Street,
			House:  address.House,
			Lat:    address.Lat,
			Lon:    address.Lon,
		}
		response.Addresses = append(response.Addresses, addressResponse)
	}

	return response, nil
}
