package rpc_v1

import (
	"geoprovider/internal/entity"
	"github.com/wanomir/e"
)

type GeoProvider interface {
	AddressSearch(input string) ([]entity.Address, error)
	GeoCode(lat, lng string) ([]entity.Address, error)
}

type GeoController struct {
	service GeoProvider
}

func NewController(service GeoProvider) *GeoController {
	return &GeoController{service: service}
}

func (c *GeoController) AddressSearch(input string, reply *entity.Addresses) error {
	addresses, err := c.service.AddressSearch(input)
	if err != nil {
		return e.Wrap("error fetching addresses", err)
	}

	*reply = entity.Addresses{
		Addresses: addresses,
	}
	return nil
}

func (c *GeoController) GeoCode(args []string, reply *entity.Addresses) error {
	addresses, err := c.service.GeoCode(args[0], args[1])
	if err != nil {
		return e.Wrap("error fetching addresses", err)
	}

	*reply = entity.Addresses{
		Addresses: addresses,
	}
	return nil
}
