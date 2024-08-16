package app

import (
	"backend/internal/lib/e"
	"backend/internal/modules"
	"backend/internal/modules/geo/usecase"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	host      string
	port      string
	redisHost string
	redisPort string
	apiKey    string
	secretKey string
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

	a.services = modules.NewServices(geoService)
	a.controllers = modules.NewControllers(a.services)

	a.server = &http.Server{
		Addr:         ":" + a.config.port,
		Handler:      a.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	a.signalChan = make(chan os.Signal, 1)
	signal.Notify(a.signalChan, syscall.SIGINT, syscall.SIGTERM)

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
		host:      os.Getenv("HOST"),
		port:      os.Getenv("PORT"),
		redisHost: os.Getenv("REDIS_HOST"),
		redisPort: os.Getenv("REDIS_PORT"),
		apiKey:    os.Getenv("DADATA_API_KEY"),
		secretKey: os.Getenv("DADATA_SECRET_KEY"),
	}

	return nil
}

func (a *App) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	r.Route("/address", func(r chi.Router) {
		r.Post("/search", a.controllers.Geo.AddressSearch)
		r.Post("/geocode", a.controllers.Geo.AddressGeocode)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s:%s/swagger/doc.json", a.config.host, a.config.port)),
	))

	return r
}
