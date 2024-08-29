package geoprovider

import (
	"errors"
	"geoservice/internal/modules/geo/dto"
	"geoservice/internal/modules/geo/infrastructure/geoprovider/grpc_v1"
	"geoservice/internal/modules/geo/infrastructure/geoprovider/jsonrpc_v1"
	"geoservice/internal/modules/geo/infrastructure/geoprovider/rpc_v1"
)

type GeoProvider interface {
	AddressSearch(input string) ([]dto.Address, error)
	GeoCode(lat, lng string) ([]dto.Address, error)
}

func NewGeoProvider(host, port, serviceName, rpcProtocol string) (GeoProvider, error) {
	switch rpcProtocol {
	case "rpc":
		return rpc_v1.NewProvider(host, port, serviceName)
	case "json-rpc":
		return jsonrpc_v1.NewProvider(host, port, serviceName)
	case "grpc":
		return grpc_v1.NewProvider(host, port)
	default:
		return nil, errors.New("unknown rpc protocol: " + rpcProtocol)
	}
}
