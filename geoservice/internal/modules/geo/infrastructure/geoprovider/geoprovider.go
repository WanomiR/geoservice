package geoprovider

import (
	"geoservice/internal/modules/geo/entity"
	"github.com/wanomir/e"
	"net/rpc"
)

type Provider struct {
	client      *rpc.Client
	serviceName string
}

func NewProvider(host, port, serviceName string) (*Provider, error) {
	client, err := rpc.Dial("tcp", host+":"+port)
	if err != nil {
		return nil, e.Wrap("error dialing rpc service", err)
	}

	return &Provider{client: client, serviceName: serviceName}, nil
}

func (p *Provider) AddressSearch(input string) ([]entity.Address, error) {
	var addresses entity.Addresses
	if err := p.client.Call(p.serviceName+".AddressSearch", input, &addresses); err != nil {
		return nil, e.Wrap("error calling "+p.serviceName+".AddressSearch", err)
	}

	return addresses.Addresses, nil
}

func (p *Provider) GeoCode(lat, lng string) ([]entity.Address, error) {
	var addresses entity.Addresses
	if err := p.client.Call(p.serviceName+".GeoCode", []string{lat, lng}, &addresses); err != nil {
		return nil, e.Wrap("error calling "+p.serviceName+".Geocode", err)
	}

	return addresses.Addresses, nil
}
