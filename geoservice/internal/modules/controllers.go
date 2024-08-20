package modules

import (
	"geoservice/internal/lib/rr"
	auth "geoservice/internal/modules/auth/controller/http/v1"
	geo "geoservice/internal/modules/geo/controller/http/v1"
)

type Controllers struct {
	Geo  geo.Controller
	Auth auth.Controller
}

func NewControllers(services *Services) *Controllers {
	return &Controllers{
		Geo:  geo.NewGeoController(services.Geo, rr.NewReadRespond()),
		Auth: auth.NewAuthController(services.Auth, rr.NewReadRespond()),
	}
}
