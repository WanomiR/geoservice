package app

import (
	"go.uber.org/zap"
	"net/http"
)

func (a *App) ZapLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		a.logger.Info("incoming request",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.String("addr", r.RemoteAddr),
		)
	})
}

func (a *App) RequireAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := a.control.VerifyRequest(w, r); err != nil {
			a.logger.Error("could not verify token", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Authorization required"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
