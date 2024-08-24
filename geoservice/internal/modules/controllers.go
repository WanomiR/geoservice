package modules

import (
	auth "geoservice/internal/modules/auth/controller/http_v1"
	geo "geoservice/internal/modules/geo/controller/http_v1"
	"github.com/wanomir/rr"
	"net/http"
)

type GeoController interface {
	AddressSearch(w http.ResponseWriter, r *http.Request)
	AddressGeocode(w http.ResponseWriter, r *http.Request)
}

type AuthController interface {
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type Controllers struct {
	Geo  GeoController
	Auth AuthController
}

func NewControllers(services *Services) *Controllers {
	return &Controllers{
		Geo:  geo.NewGeoController(services.Geo, rr.NewReadResponder()),
		Auth: auth.NewAuthController(services.Auth, rr.NewReadResponder()),
	}
}
