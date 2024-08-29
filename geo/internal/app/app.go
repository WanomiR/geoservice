package app

import (
	"errors"
	"fmt"
	"geo/internal/controller/grpc_v1"
	"geo/internal/usecase"
	pb "geo/pkg/geo_v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wanomir/e"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var appInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "geo",
	Name:      "info",
	Help:      "App environment info",
}, []string{"version"})

func init() {
	prometheus.MustRegister(appInfo)
}

type Config struct {
	host        string
	port        string
	metricsPort string
	apiKey      string
	secretKey   string
	redisHost   string
	redisPort   string
	appVersion  string
}

type App struct {
	config     Config
	server     *grpc.Server
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
	listener, err := net.Listen("tcp", a.config.host+":"+a.config.port)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("started grpc server on port", a.config.port)
	if err = a.server.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		log.Fatal(err)
	}
}

func (a *App) Shutdown() {
	<-a.signalChan

	a.server.GracefulStop()

	log.Println("shutting down gracefully...")
}

func (a *App) ServeMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", a.config.host, a.config.metricsPort), nil)) // serve metrics on a separate port
}

func (a *App) init() (err error) {
	defer func() { err = e.WrapIfErr("error initializing app", err) }()

	if err = a.readConfig(); err != nil {
		return err
	}

	service := usecase.NewGeoCacheProxy(
		usecase.NewGeoService(a.config.apiKey, a.config.secretKey),
		fmt.Sprintf("%s:%s", a.config.redisHost, a.config.redisPort), // redis address
	)

	controller := grpc_v1.NewGeoController(service)
	a.server = grpc.NewServer()

	reflection.Register(a.server)
	pb.RegisterGeoProviderV1Server(a.server, controller)

	a.signalChan = make(chan os.Signal, 1)
	signal.Notify(a.signalChan, syscall.SIGINT, syscall.SIGTERM)

	appInfo.With(prometheus.Labels{"version": a.config.appVersion}).Set(1)

	return nil
}

func (a *App) readConfig() error {
	a.config = Config{
		host:        os.Getenv("HOST"),
		port:        os.Getenv("PORT"),
		metricsPort: os.Getenv("METRICS_PORT"),
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
