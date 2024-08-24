package rpc_v1

import (
	"context"
	"geoprovider/internal/entity"
	rpc_v1 "geoprovider/pkg/geoprovider_rpc_v1"
	"github.com/wanomir/e"
)

type GeoProvider interface {
	AddressSearch(input string) ([]entity.Address, error)
	GeoCode(lat, lng string) ([]entity.Address, error)
}

type Controller struct {
	rpc_v1.UnimplementedGeoProviderV1Server // safety fallback for non-implemented controllers
	service                                 GeoProvider
}

func NewController(service GeoProvider) *Controller {
	return &Controller{service: service}
}

func (c *Controller) AddressSearch(_ context.Context, in *rpc_v1.AddressRequest) (*rpc_v1.AddressesResponse, error) {
	addresses, err := c.service.AddressSearch(in.GetQuery())
	if err != nil {
		return nil, e.Wrap("error fetching addresses", err)
	}

	response := &rpc_v1.AddressesResponse{
		Addresses: make([]*rpc_v1.AddressResponse, 0),
	}

	for _, a := range addresses {
		addressResponse := &rpc_v1.AddressResponse{
			City:   a.City,
			Street: a.Street,
			House:  a.House,
			Lat:    a.Lat,
			Lon:    a.Lon,
		}
		response.Addresses = append(response.Addresses, addressResponse)
	}

	return response, nil
}

func (c *Controller) GeoCode(_ context.Context, in *rpc_v1.GeoRequest) (*rpc_v1.AddressesResponse, error) {
	addresses, err := c.service.GeoCode(in.GetLat(), in.GetLng())
	if err != nil {
		return nil, e.Wrap("error fetching addresses", err)
	}

	response := &rpc_v1.AddressesResponse{
		Addresses: make([]*rpc_v1.AddressResponse, 0),
	}

	for _, a := range addresses {
		addressResponse := &rpc_v1.AddressResponse{
			City:   a.City,
			Street: a.Street,
			House:  a.House,
			Lat:    a.Lat,
			Lon:    a.Lon,
		}
		response.Addresses = append(response.Addresses, addressResponse)
	}

	return response, nil
}
