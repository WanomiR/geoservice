package grpc_v1

import (
	"context"
	"geo/internal/dto"
	pb "geo/pkg/geo_v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wanomir/e"
	"time"
)

var requestsCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "geo",
	Name:      "total_requests_count",
	Help:      "Total number of requests received.",
}, []string{"type"})

var requestsLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "geo",
	Name:      "response_latency_seconds",
	Help:      "Histogram of response times in seconds.",
	Buckets:   []float64{0.001, 0.002, 0.004, 0.008, 0.016, 0.032, 0.064, 0.128, 0.256, 0.512, 1.024, 2.048},
}, []string{"method"})

type GeoProvider interface {
	AddressSearch(query string) ([]dto.Address, error)
	GeoCode(lat, lng string) ([]dto.Address, error)
}

type GeoController struct {
	pb.UnimplementedGeoProviderV1Server // safety fallback for non-implemented controllers
	service                             GeoProvider
}

func NewGeoController(service GeoProvider) *GeoController {
	registerMetrics(requestsCount, requestsLatency)
	return &GeoController{service: service}
}

func (c *GeoController) AddressSearch(_ context.Context, in *pb.AddressRequest) (*pb.AddressesResponse, error) {
	// for measuring response latency
	start := time.Now()

	addresses, err := c.service.AddressSearch(in.GetQuery())
	if err != nil {
		return nil, e.Wrap("error fetching addresses", err)
	}

	response := &pb.AddressesResponse{
		Addresses: make([]*pb.AddressResponse, 0),
	}

	for _, address := range addresses {
		addressResponse := &pb.AddressResponse{
			City:   address.City,
			Street: address.Street,
			House:  address.House,
			Lat:    address.Lat,
			Lon:    address.Lon,
		}
		response.Addresses = append(response.Addresses, addressResponse)
	}

	// measure response latency and total count
	requestsLatency.With(prometheus.Labels{"method": "AddressSearch"}).Observe(time.Since(start).Seconds())
	requestsCount.With(prometheus.Labels{"type": "requests"}).Inc()

	return response, nil
}

func (c *GeoController) GeoCode(_ context.Context, in *pb.GeoRequest) (*pb.AddressesResponse, error) {
	// for measuring response time
	start := time.Now()

	addresses, err := c.service.GeoCode(in.GetLat(), in.GetLng())
	if err != nil {
		return nil, e.Wrap("error fetching address", err)
	}

	response := &pb.AddressesResponse{
		Addresses: make([]*pb.AddressResponse, 0),
	}

	for _, address := range addresses {
		addressResponse := &pb.AddressResponse{
			City:   address.City,
			Street: address.Street,
			House:  address.House,
			Lat:    address.Lat,
			Lon:    address.Lon,
		}
		response.Addresses = append(response.Addresses, addressResponse)
	}

	// measure response latency and total count
	requestsLatency.With(prometheus.Labels{"method": "GeoCode"}).Observe(time.Since(start).Seconds())
	requestsCount.With(prometheus.Labels{"type": "requests"}).Inc()

	return response, nil
}

func registerMetrics(metrics ...prometheus.Collector) {
	for _, metric := range metrics {
		prometheus.MustRegister(metric)
	}
}
