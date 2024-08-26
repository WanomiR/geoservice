package jsonrpc_v1

import (
	"geoprovider/internal/dto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wanomir/e"
	"time"
)

var requestsCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "geoprovider",
	Name:      "total_requests_count",
	Help:      "Total number of requests received.",
}, []string{"type"})

var requestsLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "geoprovider",
	Name:      "response_latency_seconds",
	Help:      "Histogram of response times in seconds.",
	Buckets:   []float64{0.001, 0.002, 0.004, 0.008, 0.016, 0.032, 0.064, 0.128, 0.256, 0.512, 1.024, 2.048},
}, []string{"method"})

type GeoProvider interface {
	AddressSearch(input string) ([]dto.Address, error)
	GeoCode(lat, lng string) ([]dto.Address, error)
}

type AddressSearchArgs struct {
	Address string
}

type GeoCodeArgs struct {
	Lat string
	Lng string
}

type GeoController struct {
	service GeoProvider
}

func NewGeoController(service GeoProvider) *GeoController {
	registerMetrics(requestsCount, requestsLatency)
	return &GeoController{service: service}
}

func (c *GeoController) AddressSearch(args *AddressSearchArgs, reply *dto.Addresses) error {

	// for measuring response latency
	start := time.Now()

	addresses, err := c.service.AddressSearch(args.Address)
	if err != nil {
		return e.Wrap("error fetching addresses", err)
	}

	*reply = dto.Addresses{
		Addresses: addresses,
	}

	// measure response latency and total count
	requestsLatency.With(prometheus.Labels{"method": "AddressSearch"}).Observe(time.Since(start).Seconds())
	requestsCount.With(prometheus.Labels{"type": "requests"}).Inc()
	return nil
}

func (c *GeoController) GeoCode(args *GeoCodeArgs, reply *dto.Addresses) error {
	// for measuring response time
	start := time.Now()

	addresses, err := c.service.GeoCode(args.Lat, args.Lng)
	if err != nil {
		return e.Wrap("error fetching addresses", err)
	}

	*reply = dto.Addresses{
		Addresses: addresses,
	}

	// measure response latency and total count
	requestsLatency.With(prometheus.Labels{"method": "GeoCode"}).Observe(time.Since(start).Seconds())
	requestsCount.With(prometheus.Labels{"type": "requests"}).Inc()
	return nil
}

func registerMetrics(metrics ...prometheus.Collector) {
	for _, metric := range metrics {
		prometheus.MustRegister(metric)
	}
}
