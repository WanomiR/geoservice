package modules

import (
	auth "geoservice/internal/modules/auth/usecase"
	geo "geoservice/internal/modules/geo/usecase"
)

type Services struct {
	Geo  geo.GeoServicer
	Auth auth.AuthServicer
}

func NewServices(geo geo.GeoServicer, auth auth.AuthServicer) *Services {
	return &Services{Geo: geo, Auth: auth}
}
