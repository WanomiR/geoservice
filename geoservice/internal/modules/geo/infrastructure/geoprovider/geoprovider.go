package geoprovider

import (
	"geoservice/internal/modules/geo/dto"
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

func (p *Provider) AddressSearch(input string) ([]dto.Address, error) {
	calling := p.serviceName + ".AddressSearch"
	var addresses dto.Addresses

	if err := p.client.Call(calling, input, &addresses); err != nil {
		return nil, e.Wrap("error calling "+calling, err)
	}

	return addresses.Addresses, nil
}

func (p *Provider) GeoCode(lat, lng string) ([]dto.Address, error) {
	calling := p.serviceName + ".GeoCode"
	var addresses dto.Addresses

	if err := p.client.Call(calling, []string{lat, lng}, &addresses); err != nil {
		return nil, e.Wrap("error calling "+calling, err)
	}

	return addresses.Addresses, nil
}
