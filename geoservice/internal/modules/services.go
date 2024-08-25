package modules

import (
	egeo "geoservice/internal/modules/geo/entity"
	"net/http"
)

type Auther interface {
	Register(email, password string) error
	Authorize(email, password string) (string, *http.Cookie, error)
	ResetCookie() *http.Cookie
	RequireAuthorization(next http.Handler) http.Handler
}

type GeoProvider interface {
	AddressSearch(input string) ([]egeo.Address, error)
	GeoCode(lat, lng string) ([]egeo.Address, error)
}

type Services struct {
	Geo  GeoProvider
	Auth Auther
}

func NewServices(geo GeoProvider, auth Auther) *Services {
	return &Services{Geo: geo, Auth: auth}
}
