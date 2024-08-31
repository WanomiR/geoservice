package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wanomir/e"
	"github.com/wanomir/rr"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"proxy/internal/controller/http_v1"
	grpc_v1 "proxy/internal/infrastructure/api_clients/geo_v1"
	"proxy/internal/usecase"
	"proxy/pkg/logger"
	"runtime/debug"
	"time"
)

const (
	exitStatusOk     = 0
	exitStatusFailed = 1
)

var appInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "proxy",
	Name:      "info",
	Help:      "App environment info",
}, []string{"version"})

func init() {
	prometheus.MustRegister(appInfo)
}

type App struct {
	config  *Config
	logger  *zap.Logger
	errChan chan error
	server  *http.Server
	control *http_v1.Controller
}

func NewApp() (*App, error) {
	a := &App{}

	if err := a.init(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() (exitCode int) {
	defer a.recoverFromPanic(&exitCode)
	var err error

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	go a.listenAndServe()
	defer a.serverShutdown(ctx)

	select {
	case err = <-a.errChan:
		a.logger.Error(e.Wrap("fatal error, service shutdown", err).Error())
	case <-ctx.Done():
		a.logger.Info("service shutdown")
	}

	return exitStatusOk
}

func (a *App) init() (err error) {
	defer func() { err = e.WrapIfErr("error initializing app", err) }()

	if err = a.readConfig(); err != nil {
		return err
	}

	a.logger = logger.NewLogger(a.config.Log.Level)
	a.errChan = make(chan error)

	usecases := usecase.NewUsecases(
		grpc_v1.NewGeoProvider(a.config.Geo.Host, a.config.Geo.Port),
	)
	a.control = http_v1.NewController(usecases, rr.NewReadResponder())

	a.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", a.config.Host, a.config.Port),
		Handler:      a.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	appInfo.With(prometheus.Labels{"version": a.config.AppVersion}).Set(1)

	return nil
}

func (a *App) readConfig() (err error) {
	defer func() { err = e.WrapIfErr("error reading config", err) }()

	a.config = new(Config)
	if err = cleanenv.ReadEnv(a.config); err != nil {
		return err
	}

	return nil
}

func (a *App) recoverFromPanic(exitCode *int) {
	if panicErr := recover(); panicErr != nil {
		a.logger.Error(
			fmt.Sprintf("recover after panic: %v, stack trace: %s", panicErr, string(debug.Stack())),
		)
		*exitCode = exitStatusFailed
	}
}

func (a *App) listenAndServe() {
	a.logger.Info("starting server on port " + a.config.Port)
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.errChan <- err
	}
}

func (a *App) serverShutdown(ctx context.Context) {
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error(e.Wrap("error attempting http server shutdown", err).Error())
	}
	a.logger.Info("http server shutdown complete")
}
