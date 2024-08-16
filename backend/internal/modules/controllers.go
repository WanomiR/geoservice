package modules

import (
	"backend/internal/lib/rr"
	geo "backend/internal/modules/geo/controller/http/v1"
)

type Controllers struct {
	Geo geo.Controller
}

func NewControllers(services *Services) *Controllers {
	return &Controllers{
		Geo: geo.NewGeoController(services.Geo, rr.NewReadRespond()),
	}
}
