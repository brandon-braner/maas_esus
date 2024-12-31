package v1

import (
	"net/http"

	"github.com/brandonbraner/maas/internal/memes"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/meme", memes.MemeGeneraterHandler)
	mux.HandleFunc("GET /v1/tokens", memes.MemeTokenHandler)
}
