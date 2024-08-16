package v1

import (
	"backend/internal/lib/e"
	"backend/internal/lib/rr"
	"backend/internal/modules/auth/usecase"
	"net/http"
)

type Controller interface {
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type AuthController struct {
	authService usecase.AuthServicer
	rr          rr.ReadResponder
}

func NewAuthController(authService usecase.AuthServicer, readResponder rr.ReadResponder) *AuthController {
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
