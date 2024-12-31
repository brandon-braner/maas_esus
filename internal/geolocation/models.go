package geolocation

type LocationInfo struct {
	Address string
}

type GeoLocationInterface interface {
	GetLocationInfo(lat, lon float64) (*LocationInfo, error)
}
