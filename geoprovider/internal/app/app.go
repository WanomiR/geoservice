package app

import (
	"context"
	"errors"
	"fmt"
	cntrl "geoprovider/internal/controller/rpc_v1"
	"geoprovider/internal/usecase"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wanomir/e"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var appInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "geoprovider",
	Name:      "info",
	Help:      "App environment info",
}, []string{"version"})

type Config struct {
	host        string
	port        string
	serviceName string
	apiKey      string
	secretKey   string
	redisHost   string
	redisPort   string
	appVersion  string
}

type App struct {
	config     Config
	server     *rpc.Server
	signalChan chan os.Signal
	controller *cntrl.GeoController
}

func NewApp() (*App, error) {
	a := &App{}

	if err := a.init(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Start() {
	listener, err := net.Listen("tcp", ":"+a.config.port)
	if err != nil {
		log.Fatal(e.Wrap("failed to listen", err))
	}

	fmt.Println("Started rpc server on port", a.config.port)
	for {
		a.server.Accept(listener)
	}

}

func (a *App) Shutdown() {
	<-a.signalChan

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	<-ctx.Done()

	fmt.Println("Shutting down server gracefully")
}

func (a *App) ServeMetrics() {
	panic("not implemented")
}

func (a *App) init() error {
	if err := a.readConfig(); err != nil {
		return err
	}

	service := usecase.NewGeoCacheProxy(
		usecase.NewGeoService(a.config.apiKey, a.config.secretKey),
		fmt.Sprintf("%s:%s", a.config.redisHost, a.config.redisPort),
	)

	a.controller = cntrl.NewController(service)

	a.server = rpc.NewServer()
	if err := a.server.RegisterName("GeoProvider", a.controller); err != nil {
		return e.Wrap("error registering rpc GeoProvider", err)
	}

	a.signalChan = make(chan os.Signal, 1)
	signal.Notify(a.signalChan, syscall.SIGINT, syscall.SIGTERM)

	appInfo.With(prometheus.Labels{"version": a.config.appVersion}).Set(1)

	return nil
}

func (a *App) readConfig() error {
	a.config = Config{
		host:        os.Getenv("HOST"),
		port:        os.Getenv("PORT"),
		serviceName: os.Getenv("SERVICE_NAME"),
		apiKey:      os.Getenv("DADATA_API_KEY"),
		secretKey:   os.Getenv("DADATA_SECRET_KEY"),
		redisHost:   os.Getenv("REDIS_HOST"),
		redisPort:   os.Getenv("REDIS_PORT"),
		appVersion:  os.Getenv("APP_VERSION"),
	}

	if a.config.host == "" || a.config.port == "" || a.config.apiKey == "" {
		return errors.New("env variables not set")
	}

	return nil
}
