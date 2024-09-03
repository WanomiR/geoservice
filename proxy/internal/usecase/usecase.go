package usecase

import (
	"proxy/internal/usecase/auth"
	"proxy/internal/usecase/geo"
)

type Usecases struct {
	*geo.GeoUsecase
	*auth.AuthUsecase
}

func NewUsecases(geoProvider geo.GeoProvider, authProvider auth.AuthProvider) *Usecases {
	return &Usecases{
		GeoUsecase:  geo.NewGeoUsecase(geoProvider),
		AuthUsecase: auth.NewAuthUsecase(authProvider),
	}
}
