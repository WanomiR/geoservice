package app

import (
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wanomir/e"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var appInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "geo",
	Name:      "info",
	Help:      "App environment info",
}, []string{"version"})

func init() {
	prometheus.MustRegister(appInfo)
}

//package config
//
//import (
//"github.com/ilyakaznacheev/cleanenv"
//)
//
//type (
//	Config struct {
//		GRPC
//		Log
//		PG
//	}
//
//	GRPC struct {
//		Port int `env-required:"true" env:"GRPC_PORT"`
//	}
//
//	Log struct {
//		Level string `env-required:"true" env:"LOG_LEVEL"`
//	}
//
//	PG struct {
//		Host          string `env-required:"true" env:"PG_HOST"`
//		User          string `env-required:"true" env:"PG_USER"`
//		Password      string `env-required:"true" env:"PG_PASSWORD"`
//		UserAdmin     string `env-required:"true" env:"PG_USER_ADMIN"`
//		PasswordAdmin string `env-required:"true" env:"PG_PASSWORD_ADMIN"`
//		Database      string `env-required:"true" env:"PG_DATABASE"`
//	}
//)
//
//func NewConfig() (*Config, error) {
//	cfg := &Config{}
//
//	err := cleanenv.ReadEnv(cfg)
//	if err != nil {
//		return nil, err
//	}
//
//	return cfg, nil
//}

type Config struct {
	host       string
	port       string
	appVersion string
}

type App struct {
	config     Config
	server     *http.Server
	signalChan chan os.Signal
	// controller
}

func NewApp() (*App, error) {
	a := &App{}

	if err := a.init(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Start() {
	fmt.Println("started server on port", a.config.port)
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func (a *App) Shutdown() {
	<-a.signalChan

	a.server.Close()

	log.Println("shutting down gracefully...")
}

//func (a *App) ServeMetrics() {
//	http.Handle("/metrics", promhttp.Handler())
//	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", a.config.host, a.config.metricsPort), nil)) // serve metrics on a separate port
//}

func (a *App) init() (err error) {
	defer func() { err = e.WrapIfErr("error initializing app", err) }()

	if err = a.readConfig(); err != nil {
		return err
	}

	a.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", a.config.host, a.config.port),
		Handler:      a.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	a.signalChan = make(chan os.Signal, 1)
	signal.Notify(a.signalChan, syscall.SIGINT, syscall.SIGTERM)

	appInfo.With(prometheus.Labels{"version": a.config.appVersion}).Set(1)

	return nil
}

func (a *App) readConfig() error {
	a.config = Config{
		host:       os.Getenv("HOST"),
		port:       os.Getenv("PORT"),
		appVersion: os.Getenv("APP_VERSION"),
	}

	if a.config.host == "" || a.config.port == "" {
		return errors.New("env variables not set")
	}

	return nil
}
