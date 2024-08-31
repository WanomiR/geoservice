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
