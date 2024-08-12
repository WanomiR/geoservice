package controller

import (
	"errors"
	"net/http"
	"proxy/internal/modules/geo/entities"
	"proxy/internal/modules/geo/service"
	"proxy/internal/utils/readresponder"
)

type Geo struct {
	geoService    service.GeoServicer
	readResponder readresponder.ReadResponder
}

func NewGeo(geoService service.GeoServicer, responder readresponder.ReadResponder) *Geo {
	return &Geo{geoService: geoService, readResponder: responder}
}

// AddressSearch
// @Summary Search by street name
// @Security ApiKeyAuth
// @Description Return a list of addresses provided street name
// @Tags address
// @Accept json
// @Produce json
// @Param query body entities.AddressSearch true "street name"
// @Success 200 {object} readresponder.JSONResponse
// @Failure 400 {object} readresponder.JSONResponse
// @Router /api/address/search [post]
func (g *Geo) AddressSearch(w http.ResponseWriter, r *http.Request) {
	var req entities.AddressSearch

	if err := g.readResponder.ReadJSON(w, r, &req); err != nil {
		g.readResponder.WriteJSONError(w, err)
		return
	} else if req.Query == "" {
		g.readResponder.WriteJSONError(w, errors.New("query is required"))
		return
	}

	addresses, _ := g.geoService.AddressSearch(req.Query)

	resp := readresponder.JSONResponse{
		Error:   false,
		Message: "search completed",
		Data:    addresses,
	}

	g.readResponder.WriteJSON(w, http.StatusOK, resp)
}

// AddressGeocode
// @Summary Search by coordinates
// @Security ApiKeyAuth
// @Description Return a list of addresses provided geo coordinates
// @Tags address
// @Accept json
// @Produce json
// @Param query body entities.AddressGeocode true "coordinates"
// @Success 200 {object} readresponder.JSONResponse
// @Failure 400 {object} readresponder.JSONResponse
// @Router /api/address/geocode [post]
func (g *Geo) AddressGeocode(w http.ResponseWriter, r *http.Request) {
	var req entities.AddressGeocode
	if err := g.readResponder.ReadJSON(w, r, &req); err != nil {
		g.readResponder.WriteJSONError(w, err)
		return
	} else if req.Lat == "" || req.Lng == "" {
		g.readResponder.WriteJSONError(w, errors.New("both lat and lng are required"))
		return
	}

	addresses, _ := g.geoService.GeoCode(req.Lat, req.Lng)

	resp := readresponder.JSONResponse{
		Error:   false,
		Message: "search completed",
		Data:    addresses,
	}

	g.readResponder.WriteJSON(w, http.StatusOK, resp)
}
