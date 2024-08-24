package modules

import (
	auth "geoservice/internal/modules/auth/controller/http/v1"
	geo "geoservice/internal/modules/geo/controller/http/v1"
	"github.com/wanomir/rr"
)

type Controllers struct {
	Geo  geo.Controller
	Auth auth.Controller
}

func NewControllers(services *Services) *Controllers {
	return &Controllers{
		Geo:  geo.NewGeoController(services.Geo, rr.NewReadResponder()),
		Auth: auth.NewAuthController(services.Auth, rr.NewReadResponder()),
	}
}
