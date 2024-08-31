package geo

import "proxy/internal/dto"

type Provider interface {
	AddressSearch(input string) ([]dto.Address, error)
	GeoCode(lat, lng string) ([]dto.Address, error)
}

type UseCase struct {
	provider Provider
}

func NewUseCase(provider Provider) *UseCase {
	return &UseCase{provider: provider}
}

func (s *UseCase) AddressSearch(input string) (result []dto.Address, err error) {
	if result, err = s.provider.AddressSearch(input); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *UseCase) GeoCode(lat, lng string) (result []dto.Address, err error) {
	if result, err = s.provider.GeoCode(lat, lng); err != nil {
		return nil, err
	}
	return result, nil
}
