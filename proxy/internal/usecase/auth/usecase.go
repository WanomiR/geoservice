package auth

type AuthProvider interface {
	Register(email, password, firstName, lastName, age string) (userId int, err error)
	Authorize(email, password string) (token string, err error)
	VerifyToken(token string) (ok bool, err error)
}

type AuthUsecase struct {
	provider AuthProvider
}

func NewAuthUsecase(provider AuthProvider) *AuthUsecase {
	return &AuthUsecase{provider: provider}
}

func (a *AuthUsecase) Register(email, password, firstName, lastName, age string) (userId int, err error) {
	if userId, err = a.provider.Register(email, password, firstName, lastName, age); err != nil {
		return 0, err
	}
	return userId, nil
}

func (a *AuthUsecase) Authorize(email, password string) (token string, err error) {
	if token, err = a.provider.Authorize(email, password); err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthUsecase) VerifyToken(token string) (ok bool, err error) {
	if ok, err = a.provider.VerifyToken(token); err != nil {
		return false, err
	}
	return ok, nil
}
