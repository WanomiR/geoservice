package usecase

import (
	"backend/internal/modules/geo/entity"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var reqTimeHist = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "api_request_duration_seconds",
	Help:    "API request duration in seconds",
	Buckets: []float64{0.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
})

func init() {
	prometheus.MustRegister(reqTimeHist)
}

//go:generate mockgen -source=./geoservice.go -destination=../../../mocks/mock_geoservice/mock_geoservice.go
type GeoServicer interface {
	AddressSearch(input string) ([]*entity.Address, error)
	GeoCode(lat, lng string) ([]*entity.Address, error)
}

type GeoService struct {
	api       *suggest.Api
	apiKey    string
	secretKey string
}

func NewGeoService(apiKey, secretKey string) *GeoService {
	endpointUrl, err := url.Parse("https://suggestions.dadata.ru/suggestions/api/4_1/rs/")
	if err != nil {
		return nil
	}

	credentials := client.Credentials{
		ApiKeyValue:    apiKey,
		SecretKeyValue: secretKey,
	}

	api := suggest.Api{
		Client: client.NewClient(endpointUrl, client.WithCredentialProvider(&credentials)),
	}

	return &GeoService{
		api:       &api,
		apiKey:    apiKey,
		secretKey: secretKey,
	}
}

func (g *GeoService) AddressSearch(input string) ([]*entity.Address, error) {
	// for measuring request duration
	start := time.Now()

	var res []*entity.Address

	rawRes, err := g.api.Address(context.Background(), &suggest.RequestParams{Query: input})
	if err != nil {
		return nil, err
	}

	for _, r := range rawRes {
		if r.Data.City == "" || r.Data.Street == "" {
			continue
		}
		res = append(res, &entity.Address{City: r.Data.City, Street: r.Data.Street, House: r.Data.House, Lat: r.Data.GeoLat, Lon: r.Data.GeoLon})
	}

	// register duration
	reqTimeHist.Observe(time.Since(start).Seconds())

	return res, nil
}

func (g *GeoService) GeoCode(lat, lng string) ([]*entity.Address, error) {
	// for measuring request duration
	start := time.Now()

	httpClient := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`{"lat": %s, "lon": %s}`, lat, lng))
	req, err := http.NewRequest("POST", "https://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address", data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", g.apiKey))

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var geoCode entity.GeoCode

	_ = json.NewDecoder(resp.Body).Decode(&geoCode)

	var res []*entity.Address
	for _, r := range geoCode.Suggestions {
		var address entity.Address
		address.City = string(r.Data.City)
		address.Street = string(r.Data.Street)
		address.House = r.Data.House
		address.Lat = r.Data.GeoLat
		address.Lon = r.Data.GeoLon

		res = append(res, &address)
	}

	// register duration
	reqTimeHist.Observe(time.Since(start).Seconds())

	return res, nil
}
