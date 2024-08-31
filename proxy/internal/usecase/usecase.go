package usecase

import "proxy/internal/usecase/geo"

type Usecases struct {
	*geo.GeoUsecase
}

func NewUsecases(geoProvider geo.GeoProvider) *Usecases {
	return &Usecases{
		GeoUsecase: geo.NewGeoUsecase(geoProvider),
	}
}
