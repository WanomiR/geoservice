package app

import (
	"context"
	"errors"
	"fmt"
	"geoservice/internal/modules"
	"geoservice/internal/modules/auth/dto"
	"geoservice/internal/modules/auth/infrastructure/repository"
	usecaseAuth "geoservice/internal/modules/auth/usecase"
	"geoservice/internal/modules/geo/infrastructure/geoprovider"
	"geoservice/internal/modules/geo/usecase"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var prometheusNamespace = "geoservice"

var appInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: prometheusNamespace,
	Name:      "info",
	Help:      "App environment info",
}, []string{"version"})

func init() {
	prometheus.MustRegister(appInfo)
}

type Config struct {
	host             string
	port             string
	jwtSecret        string
	providerName     string
	providerHost     string
	providerPort     string
	providerProtocol string
	appVersion       string
}

type App struct {
	config      Config
	server      *http.Server
	signalChan  chan os.Signal
	services    *modules.Services
	controllers *modules.Controllers
}

func NewApp() (*App, error) {
	a := &App{}

	if err := a.init(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Start() {
	fmt.Println("Started http server on port", a.config.port)
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func (a *App) Shutdown() {
	<-a.signalChan

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	<-ctx.Done()

	fmt.Println("Shutting down server gracefully")
}

func (a *App) init() (err error) {
	if err = a.readConfig(); err != nil {
		return err
	}

	geoProvider, err := geoprovider.NewGeoProvider(
		a.config.providerHost,
		a.config.providerPort,
		a.config.providerName,
		a.config.providerProtocol,
	)
	if err != nil {
		return err
	}

	geoService := usecase.NewGeoService(geoProvider)

	dbRepo := repository.NewMapDBRepo(dto.User{Email: "john.doe@gmail.com", Password: "password"})
	authService := usecaseAuth.NewAuthService(a.config.host, a.config.host, a.config.jwtSecret, a.config.host, dbRepo)

	a.services = modules.NewServices(geoService, authService)
	a.controllers = modules.NewControllers(a.services)

	a.server = &http.Server{
		Addr:         ":" + a.config.port,
		Handler:      a.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 60 * time.Second, // for profiling
	}

	a.signalChan = make(chan os.Signal, 1)
	signal.Notify(a.signalChan, syscall.SIGINT, syscall.SIGTERM)

	appInfo.With(prometheus.Labels{"version": a.config.appVersion}).Set(1)

	return nil
}

func (a *App) readConfig() error {
	a.config = Config{
		host:             os.Getenv("HOST"),
		port:             os.Getenv("PORT"),
		jwtSecret:        os.Getenv("JWT_SECRET"),
		providerName:     os.Getenv("GEOPROVIDER_NAME"),
		providerHost:     os.Getenv("GEOPROVIDER_HOST"),
		providerPort:     os.Getenv("GEOPROVIDER_PORT"),
		providerProtocol: os.Getenv("GEOPROVIDER_PROTOCOL"),
		appVersion:       os.Getenv("APP_VERSION"),
	}

	if a.config.host == "" || a.config.port == "" || a.config.jwtSecret == "" || a.config.providerHost == "" || a.config.providerPort == "" {
		return errors.New("env variables not set")
	}

	return nil
}
