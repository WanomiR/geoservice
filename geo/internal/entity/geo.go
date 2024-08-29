package entity

import (
	"context"
	"encoding/json"
	"fmt"
	"geo/internal/dto"
	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
	"github.com/wanomir/e"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const endpointBase = "https://suggestions.dadata.ru/suggestions/api/4_1/rs/"

type Geo struct {
	api    *suggest.Api
	apiKey string
}

func NewGeo(apiKey, secretKey string) *Geo {
	endpoint, err := url.Parse(endpointBase)
	if err != nil {
		log.Println(e.Wrap("error parsing base url", err))
		return nil
	}

	creds := &client.Credentials{
		ApiKeyValue:    apiKey,
		SecretKeyValue: secretKey,
	}

	api := &suggest.Api{
		Client: client.NewClient(endpoint, client.WithCredentialProvider(creds)),
	}

	return &Geo{api, apiKey}
}
func (g *Geo) SuggestByAddress(input string) ([]dto.Address, error) {
	var addresses []dto.Address

	suggestions, err := g.api.Address(context.Background(), &suggest.RequestParams{Query: input})
	if err != nil {
		return nil, e.Wrap("error fetching suggestions", err)
	}

	for _, s := range suggestions {
		if s.Data.City == "" || s.Data.Street == "" {
			continue
		}
		addresses = append(addresses, dto.Address{City: s.Data.City, Street: s.Data.Street, House: s.Data.House, Lat: s.Data.GeoLat, Lon: s.Data.GeoLon})
	}

	return addresses, nil
}

func (g *Geo) SuggestByGeoCode(lat, lng string) ([]dto.Address, error) {
	httpClient := &http.Client{}

	var payload = strings.NewReader(fmt.Sprintf(`{"lat": %s, "lon": %s}`, lat, lng))
	req, err := http.NewRequest(http.MethodPost, endpointBase+"geolocate/address", payload)
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
	defer resp.Body.Close()

	var geoCode GeoCode
	_ = json.NewDecoder(resp.Body).Decode(&geoCode)

	var addresses []dto.Address
	for _, r := range geoCode.Suggestions {
		var address dto.Address
		address.City = string(r.Data.City)
		address.Street = string(r.Data.Street)
		address.House = r.Data.House
		address.Lat = r.Data.GeoLat
		address.Lon = r.Data.GeoLon

		addresses = append(addresses, address)
	}

	return addresses, nil
}
