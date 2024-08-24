package usecase

import (
	"encoding/json"
	"geoprovider/internal/entity"
	"github.com/gomodule/redigo/redis"
	"github.com/wanomir/e"
	"log"
)

type RedisPooler interface {
	Dial() (redis.Conn, error)
}

type GeoProvider interface {
	AddressSearch(input string) ([]entity.Address, error)
	GeoCode(lat, lng string) ([]entity.Address, error)
}

type GeoCacheProxy struct {
	geo   GeoProvider
	redis *redis.Pool
}

func NewGeoCacheProxy(geoservice GeoProvider, redisAddress string) *GeoCacheProxy {
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

func (p *GeoCacheProxy) AddressSearch(input string) (addresses []entity.Address, err error) {
	conn := p.redis.Get()
	defer conn.Close()

	cachedData, err := redis.Bytes(conn.Do("GET", input))
	if err == nil {
		err = json.Unmarshal(cachedData, &addresses)
		log.Println("loaded addresses from cache")

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

func (p *GeoCacheProxy) GeoCode(lat, lng string) (addresses []entity.Address, err error) {
	conn := p.redis.Get()
	defer conn.Close()

	key := lat + lng

	cachedData, err := redis.Bytes(conn.Do("GET", key))
	if err == nil {
		err = json.Unmarshal(cachedData, &addresses)
		log.Println("loaded addresses from cache")

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
