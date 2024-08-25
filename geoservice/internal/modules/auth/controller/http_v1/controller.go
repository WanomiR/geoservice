package http_v1

import (
	"github.com/wanomir/e"
	"github.com/wanomir/rr"
	"net/http"
)

type Auther interface {
	Register(email, password string) error
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

// Register godoc
// @Summary Creates new user
// @Tags auth
// @Produce json
// @Param email formData string true "New user email"
// @Param password formData string true "New user password"
// @Success 201 {object} rr.JSONResponse
// @Failure 400 {object} rr.JSONResponse
// @Router /auth/register [post]
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		_ = c.rr.WriteJSONError(w, e.Wrap("could not parse form data", err))
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if err := c.authService.Register(email, password); err != nil {
		_ = c.rr.WriteJSONError(w, e.Wrap("could not register user", err))
		return
	}

	resp := rr.JSONResponse{Message: "user registered"}
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
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		_ = c.rr.WriteJSONError(w, e.Wrap("could not parse form data", err))
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

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
