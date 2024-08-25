package rpc_v1

import (
	"geoprovider/internal/entity"
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

func init() {
	prometheus.MustRegister(requestsCount)
	prometheus.MustRegister(requestsLatency)
}

type GeoProvider interface {
	AddressSearch(input string) ([]entity.Address, error)
	GeoCode(lat, lng string) ([]entity.Address, error)
}

type GeoController struct {
	service GeoProvider
}

func NewController(service GeoProvider) *GeoController {
	return &GeoController{service: service}
}

func (c *GeoController) AddressSearch(input string, reply *entity.Addresses) error {
	// for measuring response latency
	start := time.Now()

	addresses, err := c.service.AddressSearch(input)
	if err != nil {
		return e.Wrap("error fetching addresses", err)
	}

	*reply = entity.Addresses{
		Addresses: addresses,
	}

	// measure response latency and total count
	requestsLatency.With(prometheus.Labels{"method": "AddressSearch"}).Observe(time.Since(start).Seconds())
	requestsCount.With(prometheus.Labels{"type": "requests"}).Inc()
	return nil
}

func (c *GeoController) GeoCode(args []string, reply *entity.Addresses) error {
	// for measuring response time
	start := time.Now()

	addresses, err := c.service.GeoCode(args[0], args[1])
	if err != nil {
		return e.Wrap("error fetching addresses", err)
	}

	*reply = entity.Addresses{
		Addresses: addresses,
	}

	// measure response latency and total count
	requestsLatency.With(prometheus.Labels{"method": "GeoCode"}).Observe(time.Since(start).Seconds())
	requestsCount.With(prometheus.Labels{"type": "requests"}).Inc()
	return nil
}
