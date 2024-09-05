package http_v1

import (
	"errors"
	"fmt"
	"github.com/wanomir/e"
	"github.com/wanomir/rr"
	"go.uber.org/zap"
	"net/http"
	"proxy/internal/dto"
	"strings"
)

type SuperUsecase interface {
	AddressSearch(query string) (addresses []dto.Address, err error)
	GeoCode(lat, lng string) (addresses []dto.Address, err error)
	Register(email, password, firstName, lastName, age string) (userId int, err error)
	Authorize(email, password string) (token string, cookie *http.Cookie, err error)
	VerifyToken(token string) (ok bool, err error)
}

type RequestAddressSearch struct {
	Query string `json:"query" binding:"required" example:"Подкопаевский переулок"`
}

type RequestAddressGeocode struct {
	Lat string `json:"lat" example:"55.753214" binding:"required"`
	Lng string `json:"lng" example:"37.642589" binding:"required"`
}

type Controller struct {
	usecase SuperUsecase
	rr      *rr.ReadResponder
	logger  *zap.Logger
}

func NewController(usecase SuperUsecase, readResponder *rr.ReadResponder, logger *zap.Logger) *Controller {
	return &Controller{
		usecase: usecase,
		rr:      readResponder,
		logger:  logger,
	}
}

// AddressSearch
// @Summary Returns a list of addresses provided street name
// @Security ApiKeyAuth
// @Tags address
// @Accept json
// @Produce json
// @Param query body RequestAddressSearch true "street name"
// @Success 200 {object} rr.JSONResponse
// @Failure 400 {object} rr.JSONResponse
// @Router /address/search [post]
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
// @Summary Returns a list of addresses provided geo coordinates
// @Security ApiKeyAuth
// @Tags address
// @Accept json
// @Produce json
// @Param query body RequestAddressGeocode true "coordinates"
// @Success 200 {object} rr.JSONResponse
// @Failure 400 {object} rr.JSONResponse
// @Router /address/geocode [post]
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

// Register godoc
// @Summary Creates new user
// @Tags auth
// @Produce json
// @Param email formData string true "New user email"
// @Param password formData string true "New user password"
// @Param firstName formData string false "User first name"
// @Param lastName formData string false "User last name"
// @Param age formData int false "User age in years"
// @Success 201 {object} rr.JSONResponse
// @Failure 400 {object} rr.JSONResponse
// @Router /auth/register [post]
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		_ = c.rr.WriteJSONError(w, e.Wrap("could not parse form data", err))
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	age := r.FormValue("age")

	userId, err := c.usecase.Register(email, password, firstName, lastName, age)
	if err != nil {
		_ = c.rr.WriteJSONError(w, e.Wrap("could not register user", err))
		return
	}

	resp := rr.JSONResponse{Message: fmt.Sprintf("new user registered, id: %d", userId)}
	_ = c.rr.WriteJSON(w, 201, resp)
}

// Login godoc
// @Summary Logs user into the system
// @Tags auth
// @Produce json
// @Param email formData string true "Email for login (john.doe@gmail.com)"
// @Param password formData string true "Password for login (password)"
// @Success 200 {object} rr.JSONResponse
// @Failure 400,401 {object} rr.JSONResponse
// @Router /auth/login [post]
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		_ = c.rr.WriteJSONError(w, e.Wrap("could not parse form data", err))
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	token, cookie, err := c.usecase.Authorize(email, password)
	if err != nil {
		_ = c.rr.WriteJSONError(w, e.Wrap("couldn't authorize user", err), 401)
		return
	}

	resp := rr.JSONResponse{Error: false, Message: "user authorized", Data: token}

	http.SetCookie(w, cookie)
	_ = c.rr.WriteJSON(w, 200, resp)

}

func (c *Controller) VerifyRequest(w http.ResponseWriter, r *http.Request) error {
	token, err := c.getTokenFromHeader(w, r)
	if err != nil {
		return err
	}

	if ok, err := c.usecase.VerifyToken(token); !ok || err != nil {
		return errors.New("invalid token")
	}
	return nil
}

func (c *Controller) getTokenFromHeader(w http.ResponseWriter, r *http.Request) (string, error) {
	w.Header().Add("Vary", "Authorization")

	// get authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	// split the header on spaces
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", errors.New("invalid authorization header")
	}

	// check to see if we have the word "Bearer"
	if parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	token := parts[1]

	return token, nil
}
