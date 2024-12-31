package memes

import (
	"context"
	"fmt"

	"github.com/brandonbraner/maas/config"
	"github.com/brandonbraner/maas/internal/ai"
	"github.com/brandonbraner/maas/internal/geolocation/google"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type MemeGenerator interface {
	Generate(req MemeRequest) (MemeResponse, error)
}

func NewMemeGenerator(aiPermission bool) (MemeGenerator, error) {

	geoservice, err := google.NewGeoLocationService(config.AppConfig.GOOGLE_GEOCODE_API_KEY)

	if err != nil {
		return nil, err
	}

	var generator MemeGenerator

	switch {
	case aiPermission:
		aiService := ai.NewOpenAIMemeService()
		generator = &AITextMemeGenerator{
			GeoService: geoservice,
			AIService:  aiService,
		}
	case !aiPermission:
		generator = &TextMemeGenerator{
			GeoService: geoservice,
		}
	default:
		return nil, fmt.Errorf("invalid permission state")
	}

	return generator, nil
}

type TextMemeGenerator struct {
	GeoService *google.GeoLocationService
}

func (g *TextMemeGenerator) Generate(req MemeRequest) (MemeResponse, error) {
	ctx := req.Context
	if ctx == nil {
		ctx = context.Background()
	}
	tracer := otel.Tracer("memes")

	_, span := tracer.Start(ctx, "text-meme-generator")
	span.SetAttributes(
		attribute.String("meme.query", req.Query),
		attribute.Float64("location.lat", req.Lat),
		attribute.Float64("location.lng", req.Lng),
	)
	defer span.End()

	_, locationSpan := tracer.Start(ctx, "get-location-info")

	locationinfo, err := g.GeoService.GetLocationInfo(req.Lat, req.Lng)

	if err != nil {
		span.RecordError(err)
		return MemeResponse{}, err
	}
	locationSpan.End()

	return MemeResponse{
		Text:     fmt.Sprintf("This is a text meme about %s based at the location %s", req.Query, locationinfo.Address),
		Location: fmt.Sprintf("Location %s derived from lat/lng %f/%f", locationinfo.Address, req.Lat, req.Lng),
	}, nil
}

type AITextMemeGenerator struct {
	GeoService *google.GeoLocationService
	AIService  *ai.OpenAIMemeService
}

func (g *AITextMemeGenerator) Generate(req MemeRequest) (MemeResponse, error) {
	locationinfo, err := g.GeoService.GetLocationInfo(req.Lat, req.Lng)

	if err != nil {
		return MemeResponse{}, err
	}

	systemPrompt := `
	You are a "Meme-as-a-Service" generator. Your task is to create original and humorous memes based on user requests.
	
	**Capabilities:**
	*   **Text Generation:** You can generate witty, humorous, and relevant text captions for memes.
	*   **Concept Combination:** You can creatively combine user-provided concepts or topics to generate unexpected and funny meme ideas.
	*   **Adaptability:** You can adapt your meme generation style based on user instructions (e.g., "sarcastic," "absurdist," "wholesome").
	*   **Location** If give a location you must include that location in the generation of the meme
	
	**Response Format:**
	Return the generated meme in the following JSON format:
	{
	  "text": "Your generated meme text here",
	  "location": "Location of the meme (if provided)"
	}
	`

	var locationPrompt string

	if err == nil {
		locationPrompt = fmt.Sprintf("The full address of the location is: %s. "+
			"Please take the city and state or equivelent from it and generate the meme for that location.", locationinfo.Address)
	}

	userPrompt := locationPrompt + "\n" + req.Query

	prompts := ai.MemePrompt{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
	}

	meme, _ := g.AIService.GenerateTextMeme(&prompts)

	return MemeResponse{
		Text:     meme,
		Location: fmt.Sprintf("Location %s derived from lat/lng %f/%f", locationinfo.Address, req.Lat, req.Lng),
	}, nil
}
