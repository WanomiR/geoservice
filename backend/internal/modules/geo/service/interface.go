package service

import "proxy/internal/modules/geo/entities"

//go:generate mockgen -source=./interface.go -destination=../controller/mock_service/mock_service.go
type GeoServicer interface {
	AddressSearch(input string) ([]*entities.Address, error)
	GeoCode(lat, lng string) ([]*entities.Address, error)
}
