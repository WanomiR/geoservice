package usecase

import (
	"errors"
	"geoprovider/internal/dto"
	"geoprovider/internal/entity"
)

type GeoService struct {
	geo *entity.Geo
}

func NewGeoService(apiKey, secretKey string) *GeoService {
	return &GeoService{entity.NewGeo(apiKey, secretKey)}
}

func (g *GeoService) AddressSearch(input string) ([]dto.Address, error) {
	if input == "" {
		return []dto.Address{}, errors.New("input is empty")
	}

	addresses, err := g.geo.SuggestByAddress(input)
	if err != nil {
		return []dto.Address{}, err
	}

	return addresses, nil
}

func (g *GeoService) GeoCode(lat, lng string) ([]dto.Address, error) {
	if lat == "" || lng == "" {
		return []dto.Address{}, errors.New("incomplete coordinates")
	}

	addresses, err := g.geo.SuggestByGeoCode(lat, lng)
	if err != nil {
		return []dto.Address{}, err
	}

	return addresses, nil
}
