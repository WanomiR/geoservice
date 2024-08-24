package app

import (
	"context"
	"errors"
	"fmt"
	cntrl "geoprovider/internal/controller/rpc_v1"
	"geoprovider/internal/usecase"
	rpc_v1 "geoprovider/pkg/geoprovider_rpc_v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wanomir/e"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
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
	host       string
	port       string
	apiKey     string
	secretKey  string
	redisHost  string
	redisPort  string
	appVersion string
}

type App struct {
	config     Config
	server     *grpc.Server
	signalChan chan os.Signal
	controller *cntrl.Controller
}

func NewApp() (*App, error) {
	a := &App{}

	if err := a.init(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Start() {
	//listener, err := net.Listen("tcp", a.config.host+":"+a.config.port)
	listener, err := net.Listen("tcp", "0.0.0.0:"+a.config.port)
	if err != nil {
		log.Fatal(e.Wrap("failed to listen", err))
	}

	fmt.Println("Started grpc server on port", listener.Addr())
	if err = a.server.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
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

	a.server = grpc.NewServer()
	reflection.Register(a.server)
	rpc_v1.RegisterGeoProviderV1Server(a.server, a.controller)

	a.signalChan = make(chan os.Signal, 1)
	signal.Notify(a.signalChan, syscall.SIGINT, syscall.SIGTERM)

	appInfo.With(prometheus.Labels{"version": a.config.appVersion}).Set(1)

	return nil
}

func (a *App) readConfig() error {
	a.config = Config{
		host:       os.Getenv("HOST"),
		port:       os.Getenv("PORT"),
		apiKey:     os.Getenv("DADATA_API_KEY"),
		secretKey:  os.Getenv("DADATA_SECRET_KEY"),
		redisHost:  os.Getenv("REDIS_HOST"),
		redisPort:  os.Getenv("REDIS_PORT"),
		appVersion: os.Getenv("APP_VERSION"),
	}

	if a.config.host == "" || a.config.port == "" || a.config.apiKey == "" {
		return errors.New("env variables not set")
	}

	return nil
}
