package usecase

import (
	"backend/internal/entity"
	mock_usecase "backend/internal/mocks/mock_geoservice"
	"github.com/golang/mock/gomock"
	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"testing"
)

const (
	mockAddress = "SomeAddress"
	mockLat     = "55.12505"
	mockLng     = "52.025821"
)

func TestGeoCacheProxy_AddressSearch(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := NewMockService(controller)

	// run test with uncached data
	conn := redigomock.NewConn()
	cmd := conn.GenericCommand("SETEX").Expect([]*entity.Addresses{})
	gp := &GeoCacheProxy{mockService, &redis.Pool{Dial: func() (redis.Conn, error) {
		return conn, nil
	}}}

	t.Run("uncached", func(t *testing.T) {
		_, _ = gp.AddressSearch(mockAddress)

		if conn.Stats(cmd) != 1 {
			t.Errorf("SETEX wat not called")
		}
	})

	// run test with cached data
	conn.Clear()
	cmd = conn.Command("GET", mockAddress).Expect("")
	gp.redis = &redis.Pool{Dial: func() (redis.Conn, error) { return conn, nil }}
	t.Run("cached", func(t *testing.T) {
		_, _ = gp.AddressSearch(mockAddress)

		if conn.Stats(cmd) != 1 {
			t.Errorf("GET wat not called")
		}
	})
}

func TestGeoCacheProxy_GeoCode(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := NewMockService(controller)

	// run test with uncached data
	conn := redigomock.NewConn()
	cmd := conn.GenericCommand("SETEX").Expect([]*entity.Addresses{})
	gp := &GeoCacheProxy{mockService, &redis.Pool{Dial: func() (redis.Conn, error) {
		return conn, nil
	}}}

	t.Run("uncached", func(t *testing.T) {
		_, _ = gp.GeoCode(mockLat, mockLng)

		if conn.Stats(cmd) != 1 {
			t.Errorf("SETEX wat not called")
		}
	})

	// run test with cached data
	conn.Clear()
	cmd = conn.Command("GET", mockLat+mockLng).Expect("")
	gp.redis = &redis.Pool{Dial: func() (redis.Conn, error) { return conn, nil }}
	t.Run("cached", func(t *testing.T) {
		_, _ = gp.GeoCode(mockLat, mockLng)

		if conn.Stats(cmd) != 1 {
			t.Errorf("GET wat not called")
		}
	})
}

func NewMockService(controller *gomock.Controller) *mock_usecase.MockGeoServicer {
	mockService := mock_usecase.NewMockGeoServicer(controller)

	mockService.EXPECT().AddressSearch(gomock.Any()).Return([]*entity.Address{}, nil).AnyTimes()
	mockService.EXPECT().GeoCode(gomock.Any(), gomock.Any()).Return([]*entity.Address{}, nil).AnyTimes()

	return mockService
}
