package memes

import (
	"fmt"

	"github.com/brandonbraner/maas/config"
	"github.com/brandonbraner/maas/internal/geolocation/google"
)

type MemeGenerator interface {
	Generate(req MemeRequest) (MemeResponse, error)
}

func NewMemeGenerator(aiPermission bool) (*MemeGenerator, error) {

	geoservice, err := google.NewGeoLocationService(config.AppConfig.GOOGLE_GEOCODE_API_KEY)

	if err != nil {
		return nil, err
	}

	var generator MemeGenerator

	switch {
	case aiPermission:
		generator = &AITextMemeGenerator{
			GeoService: geoservice,
		}
	case !aiPermission:
		generator = &TextMemeGenerator{
			GeoService: geoservice,
		}
	default:
		return nil, fmt.Errorf("invalid permission state")
	}

	return &generator, nil
}

type TextMemeGenerator struct {
	GeoService *google.GeoLocationService
}

func (g *TextMemeGenerator) Generate(req MemeRequest) (MemeResponse, error) {

	locationinfo, err := g.GeoService.GetLocationInfo(req.Lat, req.Lng)

	if err != nil {
		return MemeResponse{}, err
	}

	return MemeResponse{
		Text:     fmt.Sprintf("This is a text meme about %s based at the location %s", req.Query, locationinfo.Address),
		Location: fmt.Sprintf("Location %s derived from lat/lng %f/%f", locationinfo.Address, req.Lat, req.Lng),
	}, nil
}

type AITextMemeGenerator struct {
	GeoService *google.GeoLocationService
}

func (g *AITextMemeGenerator) Generate(req MemeRequest) (MemeResponse, error) {
	locationinfo, err := g.GeoService.GetLocationInfo(req.Lat, req.Lng)

	if err != nil {
		return MemeResponse{}, err
	}

	return MemeResponse{
		Text:     fmt.Sprintf("This is a AI meme about %s based at the location %s", req.Query, locationinfo.Address),
		Location: fmt.Sprintf("Location %s derived from lat/lng %f/%f", locationinfo.Address, req.Lat, req.Lng),
	}, nil
}
