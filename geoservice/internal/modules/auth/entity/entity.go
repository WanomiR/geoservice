package entity

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type TokensPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	jwt.RegisteredClaims
}

type Auth struct {
	Issuer       string
	Audience     string
	Secret       string
	TokenExpiry  time.Duration
	CookieDomain string
	CookiePath   string
	CookieName   string
}

func NewAuth(issuer, audience, secret, cookieDomain string) *Auth {
	return &Auth{
		Issuer:       issuer,
		Audience:     audience,
		Secret:       secret,
		TokenExpiry:  15 * time.Minute,
		CookieDomain: cookieDomain,
		CookiePath:   "/",
		CookieName:   "__Host-refresh_token",
	}
}

func (a *Auth) ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[\w-.]+@([\w-]+\.)+[\w-]{2,4}$`)
	if !re.MatchString(email) {
		return errors.New("invalid email")
	}
	return nil
}

func (a *Auth) ValidatePassword(password string) error {
	re := regexp.MustCompile(`.{3,}`)
	if !re.MatchString(password) {
		return errors.New("invalid password: must be at least 3 characters long")
	}
	return nil
}

func (a *Auth) EncryptPassword(password string) (string, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}

func (a *Auth) VerifyPassword(password string, encryptedPassword string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (a *Auth) GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// set token claims
	claims := token.Claims.(jwt.MapClaims)

	claims["name"] = username
	claims["sub"] = username
	claims["aud"] = a.Audience
	claims["iss"] = a.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"

	// set the expiry time
	claims["exp"] = time.Now().UTC().Add(a.TokenExpiry).Unix()

	// create signed token
	signedToken, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil

}

func (a *Auth) CreateCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     a.CookieName,
		Path:     a.CookiePath,
		Value:    token,
		Expires:  time.Now().UTC().Add(a.TokenExpiry),
		MaxAge:   int(a.TokenExpiry.Seconds()),
		SameSite: http.SameSiteLaxMode,
		Domain:   a.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}
}

func (a *Auth) CreateExpiredCookie() *http.Cookie {
	return &http.Cookie{
		Name:     a.CookieName,
		Path:     a.CookiePath,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Domain:   a.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}
}

func (a *Auth) VerifyRequest(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	// try to get token from cookie
	token, err := a.getCookieValue(r)
	if err != nil {
		// check header if no cookie found
		if token, err = a.getTokenFromHeader(w, r); err != nil {
			return "", nil, err
		}
	}

	// prepare function for token parsing
	parseClaims := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.Secret), nil
	}

	// parse token into claims
	claims := new(Claims)

	if _, err = jwt.ParseWithClaims(token, claims, parseClaims); err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("expired token")
		}
		return "", nil, err
	}

	// check issuer
	if claims.Issuer != a.Issuer {
		return "", nil, errors.New("invalid token issuer")
	}

	return token, claims, nil
}

func (a *Auth) RequireAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _, err := a.VerifyRequest(w, r)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Authorization required"))

		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (a *Auth) getCookieValue(r *http.Request) (string, error) {
	cookie, err := r.Cookie(a.CookieName)
	if err != nil || cookie.Value == "" {
		return "", errors.New("invalid cookie")
	}
	return cookie.Value, nil
}

func (a *Auth) getTokenFromHeader(w http.ResponseWriter, r *http.Request) (string, error) {
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
