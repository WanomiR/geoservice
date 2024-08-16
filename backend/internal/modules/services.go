package modules

import "backend/internal/modules/geo/usecase"

type Services struct {
	Geo usecase.GeoServicer
}

func NewServices(geo usecase.GeoServicer) *Services {
	return &Services{Geo: geo}
}
