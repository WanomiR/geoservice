package grpc_v1

import (
	"context"
	"fmt"
	"geoservice/internal/modules/geo/dto"
	pb "geoservice/internal/modules/geo/infrastructure/geoprovider/grpc_v1/pkg/geoprovider_v1"
	"github.com/wanomir/e"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Provider struct {
	address string
}

func NewProvider(host, port string) (*Provider, error) {
	return &Provider{address: fmt.Sprintf("%s:%s", host, port)}, nil
}

func (p *Provider) AddressSearch(query string) ([]dto.Address, error) {
	client, conn, err := p.createClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.AddressSearch(ctx, &pb.AddressRequest{Query: query})
	if err != nil {
		return nil, e.Wrap("error calling AddressSearch", err)
	}

	addresses := make([]dto.Address, 0)
	for _, a := range resp.Addresses {
		address := dto.Address{
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

func (p *Provider) GeoCode(lat, lng string) ([]dto.Address, error) {
	client, conn, err := p.createClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.GeoCode(ctx, &pb.GeoRequest{Lat: lat, Lng: lng})
	if err != nil {
		return nil, e.Wrap("error calling GeoCode", err)
	}

	addresses := make([]dto.Address, 0)
	for _, a := range resp.Addresses {
		address := dto.Address{
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

func (p *Provider) createClient() (client pb.GeoProviderV1Client, conn *grpc.ClientConn, err error) {
	conn, err = grpc.NewClient(p.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, e.Wrap("error creating grpc connection", err)
	}
	client = pb.NewGeoProviderV1Client(conn)

	return client, conn, nil
}
