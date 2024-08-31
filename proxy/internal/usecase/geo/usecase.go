package geo

import "proxy/internal/dto"

type GeoProvider interface {
	AddressSearch(input string) ([]dto.Address, error)
	GeoCode(lat, lng string) ([]dto.Address, error)
}

type GeoUsecase struct {
	provider GeoProvider
}

func NewGeoUsecase(provider GeoProvider) *GeoUsecase {
	return &GeoUsecase{provider: provider}
}

func (s *GeoUsecase) AddressSearch(input string) (result []dto.Address, err error) {
	if result, err = s.provider.AddressSearch(input); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *GeoUsecase) GeoCode(lat, lng string) (result []dto.Address, err error) {
	if result, err = s.provider.GeoCode(lat, lng); err != nil {
		return nil, err
	}
	return result, nil
}
