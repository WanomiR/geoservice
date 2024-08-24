package http_v1

import (
	eauth "geoservice/internal/modules/auth/entity"
	"github.com/wanomir/e"
	"github.com/wanomir/rr"
	"net/http"
)

type Auther interface {
	Register(user eauth.User) error
	Authorize(email string, password string) (string, *http.Cookie, error)
	ResetCookie() *http.Cookie
	RequireAuthorization(next http.Handler) http.Handler
}

type AuthController struct {
	authService Auther
	rr          *rr.ReadResponder
}

func NewAuthController(authService Auther, readResponder *rr.ReadResponder) *AuthController {
	return &AuthController{
		authService: authService,
		rr:          readResponder,
	}
}

// Login godoc
// @Summary Logs user into the system
// @Tags auth
// @Produce json
// @Param email query string true "Email for login"
// @Param password query string true "Password for login"
// @Success 200 {object} rr.JSONResponse
// @Failure 401 {object} rr.JSONResponse
// @Router /auth/login [get]
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	email := query.Get("email")
	password := query.Get("password")

	token, cookie, err := c.authService.Authorize(email, password)
	if err != nil {
		_ = c.rr.WriteJSONError(w, e.Wrap("couldn't authorize user", err), 401)
		return
	}

	resp := rr.JSONResponse{Error: false, Message: "user authorized", Data: token}

	http.SetCookie(w, cookie)
	_ = c.rr.WriteJSON(w, 200, resp)

}

// Logout godoc
// @Summary Logs out current user
// @Tags auth
// @Produce json
// @Success 200 {object} rr.JSONResponse
// @Router /auth/logout [get]
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, c.authService.ResetCookie())

	resp := rr.JSONResponse{Error: false, Message: "user logged out"}
	_ = c.rr.WriteJSON(w, 200, resp)
}
