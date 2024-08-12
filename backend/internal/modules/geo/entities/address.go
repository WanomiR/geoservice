package entities

type Address struct {
	City   string `json:"city"`
	Street string `json:"street"`
	House  string `json:"house"`
	Lat    string `json:"lat"`
	Lon    string `json:"lon"`
}

type Addresses struct {
	Addresses []*Address `json:"addresses"`
}

type AddressSearch struct {
	Query string `json:"query" binding:"required" example:"Подкопаевский переулок"`
}

type AddressGeocode struct {
	Lat string `json:"lat" example:"55.753214" binding:"required"`
	Lng string `json:"lng" example:"37.642589" binding:"required"`
}
