package app

import (
	"context"
	"errors"
	"fmt"
	"geoservice/internal/lib/e"
	"geoservice/internal/modules"
	usecaseAuth "geoservice/internal/modules/auth/usecase"
	"geoservice/internal/modules/geo/usecase"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var appInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "geoservice",
	Name:      "info",
	Help:      "App environment info",
}, []string{"version"})

type Config struct {
	host       string
	port       string
	jwtSecret  string
	redisHost  string
	redisPort  string
	apiKey     string
	secretKey  string
	appVersion string
}

type App struct {
	config      Config
	server      *http.Server
	signalChan  chan os.Signal
	services    *modules.Services
	controllers *modules.Controllers
}

func init() {
	prometheus.MustRegister(appInfo)
}

func NewApp() (*App, error) {
	a := &App{}

	if err := a.init(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Start() {
	fmt.Println("Started server on port", a.config.port)
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

func (a *App) init() error {
	if err := a.readConfig(); err != nil {
		return err
	}

	geoService := usecase.NewGeoCacheProxy(
		usecase.NewGeoService(a.config.apiKey, a.config.secretKey),
		fmt.Sprintf("%s:%s", a.config.redisHost, a.config.redisPort),
	)

	authService := usecaseAuth.NewAuthService(
		a.config.host, a.config.host, a.config.jwtSecret, a.config.host,
	)

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

func (a *App) readConfig(envPath ...string) (err error) {
	if len(envPath) > 0 {
		err = godotenv.Load(envPath[0])
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		return e.Wrap("couldn't read .env file", err)
	}

	a.config = Config{
		host:       os.Getenv("HOST"),
		port:       os.Getenv("PORT"),
		jwtSecret:  os.Getenv("JWT_SECRET"),
		redisHost:  os.Getenv("REDIS_HOST"),
		redisPort:  os.Getenv("REDIS_PORT"),
		apiKey:     os.Getenv("DADATA_API_KEY"),
		secretKey:  os.Getenv("DADATA_SECRET_KEY"),
		appVersion: os.Getenv("APP_VERSION"),
	}

	return nil
}
