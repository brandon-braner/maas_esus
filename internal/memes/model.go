package memes

import (
	"context"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type MemeRequest struct {
	Lat     float64   `json:"lat" validate:"latitude"` // Latitude
	Lng     float64   `json:"lng" validate:"longitude` // Longitude
	Query   string    `json:"query"`                   // Search query
	Context context.Context `json:"-"`                 // Context for tracing
}

// MemeResponse represents the response for a generated meme
type MemeResponse struct {
	Text     string `json:"text"`     // Meme text
	Location string `json:"location"` // Location derived from lat/lng
}

func NewMemeRequest(lat float64, lng float64, query string) (MemeRequest, error) {
	request := &MemeRequest{
		Lat:   lat,
		Lng:   lng,
		Query: query,
	}

	err := validate.Struct(request)
	if err != nil {
		return MemeRequest{}, err
	}

	return *request, nil
}
