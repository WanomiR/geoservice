package http_v1

import (
	"errors"
	"geoservice/internal/modules/geo/dto"
	"github.com/wanomir/rr"
	"net/http"
)

type GeoServicer interface {
	AddressSearch(input string) ([]dto.Address, error)
	GeoCode(lat, lng string) ([]dto.Address, error)
}

type RequestAddressSearch struct {
	Query string `json:"query" binding:"required" example:"Подкопаевский переулок"`
}

type RequestAddressGeocode struct {
	Lat string `json:"lat" example:"55.753214" binding:"required"`
	Lng string `json:"lng" example:"37.642589" binding:"required"`
}

type GeoController struct {
	geoService GeoServicer
	rr         *rr.ReadResponder
}

func NewGeoController(geoService GeoServicer, responder *rr.ReadResponder) *GeoController {
	return &GeoController{geoService: geoService, rr: responder}
}

// AddressSearch
// @Summary Search by street name
// @Description Returns a list of addresses provided street name
// @Tags address
// @Accept json
// @Produce json
// @Param query body RequestAddressSearch true "street name"
// @Success 200 {object} rr.JSONResponse
// @Failure 400 {object} rr.JSONResponse
// @Router /address/search [post]
func (g *GeoController) AddressSearch(w http.ResponseWriter, r *http.Request) {
	var req RequestAddressSearch
	_ = g.rr.ReadJSON(w, r, &req)

	if req.Query == "" {
		_ = g.rr.WriteJSONError(w, errors.New("query is required"))
		return
	}

	addresses, _ := g.geoService.AddressSearch(req.Query)

	resp := rr.JSONResponse{
		Error:   false,
		Message: "search completed",
		Data:    addresses,
	}

	_ = g.rr.WriteJSON(w, http.StatusOK, resp)
}

// AddressGeocode
// @Summary Search by coordinates
// @Description Returns a list of addresses provided geo coordinates
// @Tags address
// @Accept json
// @Produce json
// @Param query body RequestAddressGeocode true "coordinates"
// @Success 200 {object} rr.JSONResponse
// @Failure 400 {object} rr.JSONResponse
// @Router /address/geocode [post]
func (g *GeoController) AddressGeocode(w http.ResponseWriter, r *http.Request) {
	var req RequestAddressGeocode
	_ = g.rr.ReadJSON(w, r, &req)

	if req.Lat == "" || req.Lng == "" {
		_ = g.rr.WriteJSONError(w, errors.New("both lat and lng are required"))
		return
	}

	addresses, _ := g.geoService.GeoCode(req.Lat, req.Lng)

	resp := rr.JSONResponse{
		Error:   false,
		Message: "search completed",
		Data:    addresses,
	}

	_ = g.rr.WriteJSON(w, http.StatusOK, resp)
}
