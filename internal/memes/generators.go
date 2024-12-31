package memes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/brandonbraner/maas/config"
	"github.com/brandonbraner/maas/internal/ai"
	"github.com/brandonbraner/maas/internal/geolocation"
	"github.com/brandonbraner/maas/internal/geolocation/google"
	"github.com/go-redis/redis/v8"
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

	// Initialize Redis client
	opt, err := redis.ParseURL(config.AppConfig.REDIS_URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}
	redisClient := redis.NewClient(opt)

	var generator MemeGenerator

	switch {
	case aiPermission:
		aiService := ai.NewOpenAIMemeService()
		generator = &AITextMemeGenerator{
			GeoService: geoservice,
			AIService:  aiService,
			Redis:      redisClient,
		}
	case !aiPermission:
		generator = &TextMemeGenerator{
			GeoService: geoservice,
			Redis:      redisClient,
		}
	default:
		return nil, fmt.Errorf("invalid permission state")
	}

	return generator, nil
}

type TextMemeGenerator struct {
	GeoService *google.GeoLocationService
	Redis      *redis.Client
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

	// Generate cache key from lat/lng
	cacheKey := fmt.Sprintf("location:%f:%f", req.Lat, req.Lng)

	// Check Redis cache first
	var locationinfo *geolocation.LocationInfo
	cachedData, err := g.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		_, locationCacheSpan := tracer.Start(ctx, "get-location-info-cache")
		// Cache hit
		if err := json.Unmarshal([]byte(cachedData), &locationinfo); err == nil {
			return MemeResponse{
				Text:     fmt.Sprintf("This is a text meme about %s based at the location %s", req.Query, locationinfo.Address),
				Location: fmt.Sprintf("Location %s derived from lat/lng %f/%f", locationinfo.Address, req.Lat, req.Lng),
			}, nil
		}
		locationCacheSpan.End()
	}

	// Cache miss - fetch from GeoService
	_, locationSpan := tracer.Start(ctx, "get-location-info")
	locationinfo, err = g.GeoService.GetLocationInfo(req.Lat, req.Lng)
	if err != nil {
		span.RecordError(err)
		return MemeResponse{}, err
	}
	locationSpan.End()

	// Cache the result
	locationData, err := json.Marshal(locationinfo)
	if err == nil {
		g.Redis.Set(ctx, cacheKey, locationData, time.Duration(config.AppConfig.REDIS_CACHE_TTL)*time.Second)
	}

	return MemeResponse{
		Text:     fmt.Sprintf("This is a text meme about %s based at the location %s", req.Query, locationinfo.Address),
		Location: fmt.Sprintf("Location %s derived from lat/lng %f/%f", locationinfo.Address, req.Lat, req.Lng),
	}, nil
}

type AITextMemeGenerator struct {
	GeoService *google.GeoLocationService
	AIService  *ai.OpenAIMemeService
	Redis      *redis.Client
}

func (g *AITextMemeGenerator) Generate(req MemeRequest) (MemeResponse, error) {
	ctx := req.Context
	if ctx == nil {
		ctx = context.Background()
	}

	// Generate cache key from lat/lng
	cacheKey := fmt.Sprintf("location:%f:%f", req.Lat, req.Lng)

	// Check Redis cache first
	var locationinfo *geolocation.LocationInfo
	cachedData, err := g.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit
		if err := json.Unmarshal([]byte(cachedData), &locationinfo); err == nil {
			// Continue with AI meme generation using cached location
			return g.generateMemeFromLocation(locationinfo, req)
		}
	}

	// Cache miss - fetch from GeoService
	locationinfo, err = g.GeoService.GetLocationInfo(req.Lat, req.Lng)
	if err != nil {
		return MemeResponse{}, err
	}

	// Cache the result
	locationData, err := json.Marshal(locationinfo)
	if err == nil {
		g.Redis.Set(ctx, cacheKey, locationData, time.Duration(config.AppConfig.REDIS_CACHE_TTL)*time.Second)
	}

	return g.generateMemeFromLocation(locationinfo, req)
}

func (g *AITextMemeGenerator) generateMemeFromLocation(locationinfo *geolocation.LocationInfo, req MemeRequest) (MemeResponse, error) {
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

	userPrompt := ""
	if locationinfo != nil {
		locationPrompt := fmt.Sprintf("The full address of the location is: %s. "+
			"Please take the city and state or equivalent from it and generate the meme for that location.", locationinfo.Address)
		userPrompt = locationPrompt + "\n" + req.Query
	} else {
		userPrompt = req.Query
	}

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
