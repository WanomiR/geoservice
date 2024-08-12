package controller

import (
	"net/http"
)

type GeoServicer interface {
	AddressSearch(w http.ResponseWriter, r *http.Request)
	AddressGeocode(w http.ResponseWriter, r *http.Request)
}
