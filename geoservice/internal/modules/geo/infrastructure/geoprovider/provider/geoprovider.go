package provider

import (
	"context"
	"geoservice/internal/modules/geo/entity"
	"geoservice/internal/modules/geo/infrastructure/geoprovider/pkg/geoprovider_rpc_v1"
	"github.com/wanomir/e"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Provider struct {
	address string
}

func NewProvider(host, port string) *Provider {
	return &Provider{address: host + ":" + port}
}

func (p *Provider) AddressSearch(input string) ([]entity.Address, error) {
	client, conn, err := p.createClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	request := geoprovider_rpc_v1.AddressRequest{Query: input}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.AddressSearch(ctx, &request)
	if err != nil {
		return nil, e.Wrap("error calling AddressSearch", err)
	}

	addresses := make([]entity.Address, 0)
	for _, a := range response.Addresses {
		address := entity.Address{
			City:   a.City,
			Street: a.Street,
			House:  a.House,
			Lat:    a.Lat,
			Lon:    a.Lon,
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (p *Provider) GeoCode(lat, lng string) ([]entity.Address, error) {
	client, conn, err := p.createClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	request := geoprovider_rpc_v1.GeoRequest{Lat: lat, Lng: lng}

	response, err := client.GeoCode(context.Background(), &request)
	if err != nil {
		return nil, e.Wrap("error calling GeoCode", err)
	}

	addresses := make([]entity.Address, 0)
	for _, a := range response.Addresses {
		address := entity.Address{
			City:   a.City,
			Street: a.Street,
			House:  a.House,
			Lat:    a.Lat,
			Lon:    a.Lon,
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (p *Provider) createClient() (client geoprovider_rpc_v1.GeoProviderV1Client, conn *grpc.ClientConn, err error) {
	conn, err = grpc.NewClient(p.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, e.Wrap("error creating grpc connection", err)
	}
	client = geoprovider_rpc_v1.NewGeoProviderV1Client(conn)

	return client, conn, nil
}
