package entity

type GeoCode struct {
	Suggestions []Suggestion `json:"suggestions"`
}

type Suggestion struct {
	Value             string `json:"value"`
	UnrestrictedValue string `json:"unrestricted_value"`
	Data              Data   `json:"data"`
}

type Data struct {
	PostalCode           string          `json:"postal_code"`
	Country              Country         `json:"country"`
	CountryISOCode       CountryISOCode  `json:"country_iso_code"`
	FederalDistrict      FederalDistrict `json:"federal_district"`
	RegionFiasID         string          `json:"region_fias_id"`
	RegionKladrID        string          `json:"region_kladr_id"`
	RegionISOCode        RegionISOCode   `json:"region_iso_code"`
	RegionWithType       WithType        `json:"region_with_type"`
	RegionType           Type            `json:"region_type"`
	RegionTypeFull       TypeFull        `json:"region_type_full"`
	Region               City            `json:"region"`
	AreaFiasID           interface{}     `json:"area_fias_id"`
	AreaKladrID          interface{}     `json:"area_kladr_id"`
	AreaWithType         interface{}     `json:"area_with_type"`
	AreaType             interface{}     `json:"area_type"`
	AreaTypeFull         interface{}     `json:"area_type_full"`
	Area                 interface{}     `json:"area"`
	CityFiasID           string          `json:"city_fias_id"`
	CityKladrID          string          `json:"city_kladr_id"`
	CityWithType         WithType        `json:"city_with_type"`
	CityType             Type            `json:"city_type"`
	CityTypeFull         TypeFull        `json:"city_type_full"`
	City                 City            `json:"city"`
	CityArea             interface{}     `json:"city_area"`
	CityDistrictFiasID   interface{}     `json:"city_district_fias_id"`
	CityDistrictKladrID  interface{}     `json:"city_district_kladr_id"`
	CityDistrictWithType interface{}     `json:"city_district_with_type"`
	CityDistrictType     interface{}     `json:"city_district_type"`
	CityDistrictTypeFull interface{}     `json:"city_district_type_full"`
	CityDistrict         interface{}     `json:"city_district"`
	SettlementFiasID     interface{}     `json:"settlement_fias_id"`
	SettlementKladrID    interface{}     `json:"settlement_kladr_id"`
	SettlementWithType   interface{}     `json:"settlement_with_type"`
	SettlementType       interface{}     `json:"settlement_type"`
	SettlementTypeFull   interface{}     `json:"settlement_type_full"`
	Settlement           interface{}     `json:"settlement"`
	StreetFiasID         string          `json:"street_fias_id"`
	StreetKladrID        string          `json:"street_kladr_id"`
	StreetWithType       StreetWithType  `json:"street_with_type"`
	StreetType           StreetType      `json:"street_type"`
	StreetTypeFull       StreetTypeFull  `json:"street_type_full"`
	Street               Street          `json:"street"`
	SteadFiasID          interface{}     `json:"stead_fias_id"`
	SteadCadnum          interface{}     `json:"stead_cadnum"`
	SteadType            interface{}     `json:"stead_type"`
	SteadTypeFull        interface{}     `json:"stead_type_full"`
	Stead                interface{}     `json:"stead"`
	HouseFiasID          *string         `json:"house_fias_id"`
	HouseKladrID         *string         `json:"house_kladr_id"`
	HouseCadnum          interface{}     `json:"house_cadnum"`
	HouseType            HouseType       `json:"house_type"`
	HouseTypeFull        HouseTypeFull   `json:"house_type_full"`
	House                string          `json:"house"`
	BlockType            *string         `json:"block_type"`
	BlockTypeFull        *string         `json:"block_type_full"`
	Block                *string         `json:"block"`
	Entrance             interface{}     `json:"entrance"`
	Floor                interface{}     `json:"floor"`
	FlatFiasID           interface{}     `json:"flat_fias_id"`
	FlatCadnum           interface{}     `json:"flat_cadnum"`
	FlatType             interface{}     `json:"flat_type"`
	FlatTypeFull         interface{}     `json:"flat_type_full"`
	Flat                 interface{}     `json:"flat"`
	FlatArea             interface{}     `json:"flat_area"`
	SquareMeterPrice     interface{}     `json:"square_meter_price"`
	FlatPrice            interface{}     `json:"flat_price"`
	RoomFiasID           interface{}     `json:"room_fias_id"`
	RoomCadnum           interface{}     `json:"room_cadnum"`
	RoomType             interface{}     `json:"room_type"`
	RoomTypeFull         interface{}     `json:"room_type_full"`
	Room                 interface{}     `json:"room"`
	PostalBox            interface{}     `json:"postal_box"`
	FiasID               string          `json:"fias_id"`
	FiasCode             interface{}     `json:"fias_code"`
	FiasLevel            string          `json:"fias_level"`
	FiasActualityState   string          `json:"fias_actuality_state"`
	KladrID              string          `json:"kladr_id"`
	GeonameID            string          `json:"geoname_id"`
	CapitalMarker        string          `json:"capital_marker"`
	Okato                string          `json:"okato"`
	Oktmo                string          `json:"oktmo"`
	TaxOffice            string          `json:"tax_office"`
	TaxOfficeLegal       string          `json:"tax_office_legal"`
	Timezone             interface{}     `json:"timezone"`
	GeoLat               string          `json:"geo_lat"`
	GeoLon               string          `json:"geo_lon"`
	BeltwayHit           interface{}     `json:"beltway_hit"`
	BeltwayDistance      interface{}     `json:"beltway_distance"`
	Metro                interface{}     `json:"metro"`
	Divisions            interface{}     `json:"divisions"`
	QcGeo                string          `json:"qc_geo"`
	QcComplete           interface{}     `json:"qc_complete"`
	QcHouse              interface{}     `json:"qc_house"`
	HistoryValues        []string        `json:"history_values"`
	UnparsedParts        interface{}     `json:"unparsed_parts"`
	Source               interface{}     `json:"source"`
	Qc                   interface{}     `json:"qc"`
}

type City string

const (
	СанктПетербург City = "Санкт-Петербург"
)

type Type string

const (
	Г Type = "г"
)

type TypeFull string

const (
	Город TypeFull = "город"
)

type WithType string

const (
	ГСанктПетербург WithType = "г Санкт-Петербург"
)

type Country string

const (
	Россия Country = "Россия"
)

type CountryISOCode string

const (
	Ru CountryISOCode = "RU"
)

type FederalDistrict string

const (
	СевероЗападный FederalDistrict = "Северо-Западный"
)

type HouseType string

const (
	Д HouseType = "д"
)

type HouseTypeFull string

const (
	Дом HouseTypeFull = "дом"
)

type RegionISOCode string

const (
	RuSpe RegionISOCode = "RU-SPE"
)

type Street string

const (
	Казанская        Street = "Казанская"
	КаналаГрибоедова Street = "Канала Грибоедова"
	Фонарный         Street = "Фонарный"
)

type StreetType string

const (
	Наб StreetType = "наб"
	Пер StreetType = "пер"
	Ул  StreetType = "ул"
)

type StreetTypeFull string

const (
	Набережная StreetTypeFull = "набережная"
	Переулок   StreetTypeFull = "переулок"
	Улица      StreetTypeFull = "улица"
)

type StreetWithType string

const (
	НабКаналаГрибоедова StreetWithType = "наб Канала Грибоедова"
	УлКазанская         StreetWithType = "ул Казанская"
	ФонарныйПер         StreetWithType = "Фонарный пер"
)
