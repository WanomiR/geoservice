package usecase

import (
	"geoprovider/internal/dto"
	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"testing"
)

const (
	mockAddress = "SomeAddress"
	mockLat     = "55.12505"
	mockLng     = "52.025821"
)

var cacheProxy *GeoCacheProxy

func init() {
	cacheProxy = NewGeoCacheProxy(NewMockGeoService(), "")
}

func TestGeoCacheProxy_AddressSearch(t *testing.T) {
	// run test with uncached data
	conn := redigomock.NewConn()
	cmd := conn.GenericCommand("SETEX").Expect([]dto.Address{})
	cacheProxy.redis = &redis.Pool{Dial: func() (redis.Conn, error) {
		return conn, nil
	}}

	t.Run("uncached", func(t *testing.T) {
		_, _ = cacheProxy.AddressSearch(mockAddress)

		if conn.Stats(cmd) != 1 {
			t.Errorf("SETEX wat not called")
		}
	})

	// run test with cached data
	conn.Clear()
	cmd = conn.Command("GET", mockAddress).Expect("")
	cacheProxy.redis = &redis.Pool{Dial: func() (redis.Conn, error) { return conn, nil }}
	t.Run("cached", func(t *testing.T) {
		_, _ = cacheProxy.AddressSearch(mockAddress)

		if conn.Stats(cmd) != 1 {
			t.Errorf("GET wat not called")
		}
	})
}

func TestGeoCacheProxy_GeoCode(t *testing.T) {
	// run test with uncached data
	conn := redigomock.NewConn()
	cmd := conn.GenericCommand("SETEX").Expect([]dto.Address{})
	cacheProxy.redis = &redis.Pool{Dial: func() (redis.Conn, error) {
		return conn, nil
	}}

	t.Run("uncached", func(t *testing.T) {
		_, _ = cacheProxy.GeoCode(mockLat, mockLng)

		if conn.Stats(cmd) != 1 {
			t.Errorf("SETEX wat not called")
		}
	})

	// run test with cached data
	conn.Clear()
	cmd = conn.Command("GET", mockLat+mockLng).Expect("")
	cacheProxy.redis = &redis.Pool{Dial: func() (redis.Conn, error) { return conn, nil }}
	t.Run("cached", func(t *testing.T) {
		_, _ = cacheProxy.GeoCode(mockLat, mockLng)

		if conn.Stats(cmd) != 1 {
			t.Errorf("GET wat not called")
		}
	})
}

type MockGeoService struct{}

func NewMockGeoService() *MockGeoService {
	return &MockGeoService{}
}

func (m *MockGeoService) AddressSearch(_ string) ([]dto.Address, error) {
	return []dto.Address{}, nil
}

func (m *MockGeoService) GeoCode(_, _ string) ([]dto.Address, error) {
	return []dto.Address{}, nil
}
