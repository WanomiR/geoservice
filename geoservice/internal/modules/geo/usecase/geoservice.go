package usecase

import (
	"geoservice/internal/modules/geo/dto"
)

type GeoProvider interface {
	AddressSearch(input string) ([]dto.Address, error)
	GeoCode(lat, lng string) ([]dto.Address, error)
}

type GeoService struct {
	geoProvider GeoProvider
}

func NewGeoService(provider GeoProvider) *GeoService {
	return &GeoService{geoProvider: provider}
}

func (g *GeoService) AddressSearch(input string) ([]dto.Address, error) {
	res, err := g.geoProvider.AddressSearch(input)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (g *GeoService) GeoCode(lat, lng string) ([]dto.Address, error) {
	res, err := g.geoProvider.GeoCode(lat, lng)
	if err != nil {
		return nil, err
	}

	return res, nil
}
