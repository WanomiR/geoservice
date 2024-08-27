package jsonrpc_v1

import (
	"geoservice/internal/modules/geo/dto"
	"github.com/wanomir/e"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type AddressSearchArgs struct {
	Address string
}

type GeoCodeArgs struct {
	Lat string
	Lng string
}

type Provider struct {
	client      *rpc.Client
	serviceName string
}

func NewProvider(host, port, serviceName string) (*Provider, error) {
	client, err := jsonrpc.Dial("tcp", host+":"+port)
	if err != nil {
		return nil, e.Wrap("error dialing json-rpc service", err)
	}

	return &Provider{client: client, serviceName: serviceName}, nil
}

func (p *Provider) AddressSearch(input string) ([]dto.Address, error) {
	calling := p.serviceName + ".AddressSearch"
	args := &AddressSearchArgs{Address: input}

	var addresses dto.Addresses
	if err := p.client.Call(calling, args, &addresses); err != nil {
		return nil, e.Wrap("error calling "+calling, err)
	}

	return addresses.Addresses, nil
}

func (p *Provider) GeoCode(lat, lng string) ([]dto.Address, error) {
	calling := p.serviceName + ".GeoCode"
	args := &GeoCodeArgs{Lat: lat, Lng: lng}

	var addresses dto.Addresses
	if err := p.client.Call(calling, args, &addresses); err != nil {
		return nil, e.Wrap("error calling "+calling, err)
	}

	return addresses.Addresses, nil
}
