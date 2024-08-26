package app

import (
	"context"
	"errors"
	"fmt"
	"geoprovider/internal/usecase"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wanomir/e"
	"log"
	"net/http"
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

func init() {
	prometheus.MustRegister(appInfo)
}

type Config struct {
	host        string
	port        string
	serviceName string
	rpcProtocol string
	apiKey      string
	secretKey   string
	redisHost   string
	redisPort   string
	appVersion  string
}

type App struct {
	config     Config
	server     Server
	signalChan chan os.Signal
}

func NewApp() (*App, error) {
	a := &App{}

	if err := a.init(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Start() {
	fmt.Println("Started "+a.config.rpcProtocol+" server on port", a.config.port)

	log.Fatal(a.server.ListenAndServe())
}

func (a *App) Shutdown() {
	<-a.signalChan

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	<-ctx.Done()

	fmt.Println("Shutting down server gracefully")
}

func (a *App) ServeMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":7778", nil))
}

func (a *App) init() (err error) {
	defer func() { err = e.WrapIfErr("error initializing app", err) }()

	if err = a.readConfig(); err != nil {
		return err
	}

	if a.server, err = a.createServer(a.config.rpcProtocol, a.config.serviceName, a.config.port); err != nil {
		return err
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
		rpcProtocol: os.Getenv("RPC_PROTOCOL"),
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

func (a *App) createServer(rpcProtocol, serviceName, port string) (Server, error) {
	geoUsecase := usecase.NewGeoCacheProxy(
		usecase.NewGeoService(a.config.apiKey, a.config.secretKey),
		fmt.Sprintf("%s:%s", a.config.redisHost, a.config.redisPort),
	)
	serverConfig := ServerConfig{serviceName, port}

	switch rpcProtocol {
	case "rpc":
		return NewRpcServer(geoUsecase, serverConfig), nil
	case "json-rpc":
		return NewJsonRpcServer(geoUsecase, serverConfig), nil
	default:
		return nil, errors.New("invalid protocol")
	}
}
