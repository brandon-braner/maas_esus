package google

import (
	"github.com/brandonbraner/maas/internal/geolocation"
	"googlemaps.github.io/maps"
)

func NewGeoLocationService(apiKey string) (*GeoLocationService, error) {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GeoLocationService{client: client}, nil
}

func (g *GeoLocationService) GetLocationInfo(lat, lng float64) (*geolocation.LocationInfo, error) {
	// r := &maps.GeocodingRequest{
	// 	LatLng: &maps.LatLng{
	// 		Lat: lat,
	// 		Lng: lng,
	// 	},
	// }

	// results, err := g.client.ReverseGeocode(context.Background(), r)
	// if err != nil {
	// 	return nil, err
	// }

	// if len(results) == 0 {
	// 	return &geolocation.LocationInfo{}, nil
	// }

	// we are just returning the full address and going to let the llm handle the rest
	info := &geolocation.LocationInfo{}
	info.Address = "Your address based on your %d lat, %d lng" //google is slowing down the ability to hit 100 rs

	return info, nil
}
