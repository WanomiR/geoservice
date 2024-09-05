package app

import (
	"auth/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wanomir/e"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
)

const (
	exitStatusOk     = 0
	exitStatusFailed = 1
)

var appInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "auth",
	Name:      "info",
	Help:      "App environment info",
}, []string{"version"})

func init() {
	prometheus.MustRegister(appInfo)
}

type App struct {
	config *Config

	ctx     context.Context
	errChan chan error

	logger *zap.Logger
	server *grpc.Server
}

func NewApp() (*App, error) {
	a := &App{}

	if err := a.init(); err != nil {
		return nil, e.Wrap("failed to init app", err)
	}

	return a, nil
}

func (a *App) Run() (exitCode int) {
	defer a.recoverFromPanic(&exitCode)
	var err error

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	go a.listenAndServe()
	defer a.serverShutdown()

	select {
	case err = <-a.errChan:
		a.logger.Error(e.Wrap("fatal error, service shutdown", err).Error())
	case <-ctx.Done():
		a.logger.Info("service shutdown")
	}

	return exitStatusOk
}

func (a *App) ServeMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", a.config.Host, a.config.MetricsPort), nil)) // serve metrics on a separate port
}

func (a *App) serverShutdown() {
	a.server.GracefulStop()
	a.logger.Info("grpcs server shutdown")
}

func (a *App) init() (err error) {
	if err = a.readConfig(); err != nil {
		return e.Wrap("failed to read config", err)
	}

	a.logger = logger.NewLogger(a.config.Log.Level)
	a.errChan = make(chan error)

	a.server = grpc.NewServer()
	reflection.Register(a.server)
	//pb.RegisterGeoProviderV1Server(a.server, controller)

	//service := usecase.NewGeoCacheProxy(
	//	usecase.NewGeoService(a.config.apiKey, a.config.secretKey),
	//	fmt.Sprintf("%s:%s", a.config.redisHost, a.config.redisPort), // redis address
	//)
	//
	//controller := grpc_v1.NewGeoController(service)
	//a.server = grpc.NewServer()
	//
	//reflection.Register(a.server)
	//pb.RegisterGeoProviderV1Server(a.server, controller)

	appInfo.With(prometheus.Labels{"version": a.config.AppVersion}).Set(1)

	return nil
}

func (a *App) readConfig() error {
	a.config = new(Config)
	if err := cleanenv.ReadEnv(a.config); err != nil {
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
	listener, err := net.Listen("tcp", a.config.Host+":"+a.config.Port)
	if err != nil {
		a.logger.Error("failed to listen", zap.Error(err))
		a.errChan <- err
		return
	}

	a.logger.Info("started grpc server", zap.String("address", a.config.Host+":"+a.config.Port))
	if err = a.server.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		a.logger.Error("fatal error, service shutdown", zap.Error(err))
		a.errChan <- err
		return
	}
}
