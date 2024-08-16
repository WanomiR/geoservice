package usecase

import (
	"backend/internal/lib/e"
	"backend/internal/modules/geo/entity"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)

var cacheTimeHist = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "cache_access_duration_seconds",
	Help:    "Cached data access duration in seconds",
	Buckets: []float64{0.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
})

func init() {
	prometheus.MustRegister(cacheTimeHist)
}

//go:generate mockgen -source=./proxy.go -destination=../../../mocks/mock_redis/mock_redis.go
type RedisPooler interface {
	Dial() (redis.Conn, error)
}

type GeoCacheProxy struct {
	geo   GeoServicer
	redis *redis.Pool
}

func NewGeoCacheProxy(geoservice GeoServicer, redisAddress string) *GeoCacheProxy {
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisAddress)
		},
	}
	return &GeoCacheProxy{
		geo:   geoservice,
		redis: redisPool,
	}
}

func (p *GeoCacheProxy) AddressSearch(input string) (addresses []*entity.Address, err error) {
	// for measuring data access duration
	start := time.Now()

	conn := p.redis.Get()
	defer conn.Close()

	cachedData, err := redis.Bytes(conn.Do("GET", input))
	if err == nil {
		err = json.Unmarshal(cachedData, &addresses)
		log.Println("loaded addresses from cache")

		// register duration
		cacheTimeHist.Observe(time.Since(start).Seconds())

		return addresses, err
	}

	if addresses, err = p.geo.AddressSearch(input); err != nil {
		return nil, e.Wrap("error fetching addresses", err)
	}

	serialized, _ := json.Marshal(addresses)
	if _, err = conn.Do("SETEX", input, 60*60*24, serialized); err != nil {
		log.Println("error caching addresses", err)
	}
	log.Println("cached addresses")

	return addresses, nil
}

func (p *GeoCacheProxy) GeoCode(lat, lng string) (addresses []*entity.Address, err error) {
	// for measuring data access duration
	start := time.Now()

	conn := p.redis.Get()
	defer conn.Close()

	key := lat + lng

	cachedData, err := redis.Bytes(conn.Do("GET", key))
	if err == nil {
		err = json.Unmarshal(cachedData, &addresses)
		log.Println("loaded addresses from cache")

		// register duration
		cacheTimeHist.Observe(time.Since(start).Seconds())

		return addresses, err
	}

	if addresses, err = p.geo.GeoCode(lat, lng); err != nil {
		return nil, e.Wrap("error fetching addresses", err)
	}

	serialized, err := json.Marshal(addresses)
	if _, err = conn.Do("SETEX", key, 60*60*24, serialized); err != nil {
		log.Println("error caching addresses", err)
	}
	log.Println("cached addresses")

	return addresses, nil
}
