package grpc_v1

import (
	"context"
	"fmt"
	"github.com/wanomir/e"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"proxy/internal/dto"
	pb "proxy/internal/infrastructure/api_clients/geo_v1/pkg/geo_v1"
	"time"
)

type GeoProvider struct {
	address string
}

func NewGeoProvider(host, port string) *GeoProvider {
	return &GeoProvider{address: fmt.Sprintf("%s:%s", host, port)}
}

func (p *GeoProvider) AddressSearch(query string) ([]dto.Address, error) {
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

func (p *GeoProvider) GeoCode(lat, lng string) ([]dto.Address, error) {
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

func (p *GeoProvider) createClient() (client pb.GeoProviderV1Client, conn *grpc.ClientConn, err error) {
	conn, err = grpc.NewClient(p.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, e.Wrap("error creating grpc connection", err)
	}
	client = pb.NewGeoProviderV1Client(conn)

	return client, conn, nil
}
