package v1

import (
	"net/http"

	"github.com/brandonbraner/maas/internal/memes"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /v1/meme", otelhttp.NewHandler(
		http.HandlerFunc(memes.MemeGeneraterHandler),
		"v1/meme",
	))
	mux.Handle("GET /v1/tokens", otelhttp.NewHandler(
		http.HandlerFunc(memes.MemeTokenHandler),
		"v1/tokens",
	))
}
