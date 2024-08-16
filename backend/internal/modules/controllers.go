package modules

import (
	"backend/internal/lib/rr"
	auth "backend/internal/modules/auth/controller/http/v1"
	geo "backend/internal/modules/geo/controller/http/v1"
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
