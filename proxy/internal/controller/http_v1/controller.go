package http_v1

import (
	"errors"
	"github.com/wanomir/rr"
	"net/http"
	"proxy/internal/dto"
)

type SuperUseCase interface {
	AddressSearch(query string) ([]dto.Address, error)
	GeoCode(lat, lng string) ([]dto.Address, error)
}

type RequestAddressSearch struct {
	Query string `json:"query" binding:"required" example:"Подкопаевский переулок"`
}

type RequestAddressGeocode struct {
	Lat string `json:"lat" example:"55.753214" binding:"required"`
	Lng string `json:"lng" example:"37.642589" binding:"required"`
}

type Controller struct {
	usecase SuperUseCase
	rr      *rr.ReadResponder
}

func NewController(usecase SuperUseCase, readResponder *rr.ReadResponder) *Controller {
	return &Controller{
		usecase: usecase,
		rr:      readResponder,
	}
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
// @Router /api/address/search [post]
func (c *Controller) AddressSearch(w http.ResponseWriter, r *http.Request) {
	var req RequestAddressSearch
	_ = c.rr.ReadJSON(w, r, &req)

	if req.Query == "" {
		_ = c.rr.WriteJSONError(w, errors.New("query is required"))
		return
	}

	addresses, _ := c.usecase.AddressSearch(req.Query)

	resp := rr.JSONResponse{
		Error:   false,
		Message: "search completed",
		Data:    addresses,
	}

	_ = c.rr.WriteJSON(w, http.StatusOK, resp)
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
// @Router /api/address/geocode [post]
func (c *Controller) AddressGeocode(w http.ResponseWriter, r *http.Request) {
	var req RequestAddressGeocode
	_ = c.rr.ReadJSON(w, r, &req)

	if req.Lat == "" || req.Lng == "" {
		_ = c.rr.WriteJSONError(w, errors.New("both lat and lng are required"))
		return
	}

	addresses, _ := c.usecase.GeoCode(req.Lat, req.Lng)

	resp := rr.JSONResponse{
		Error:   false,
		Message: "search completed",
		Data:    addresses,
	}

	_ = c.rr.WriteJSON(w, http.StatusOK, resp)
}
