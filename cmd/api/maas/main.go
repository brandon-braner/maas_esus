package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "github.com/brandonbraner/maas/internal/api/v1"
	"github.com/brandonbraner/maas/config"
	"github.com/brandonbraner/maas/pkg/http/middleware"
	"github.com/urfave/negroni"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	// Initialize tracing
	shutdown := config.InitTracer()
	defer shutdown(context.Background())

	// Set up signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	n := negroni.Classic()
	n.Use(middleware.AuthMiddleware{})

	// Wrap the default handler with OpenTelemetry instrumentation
	mux.Handle("/", otelhttp.NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		}),
		"root",
	))

	n.UseHandler(mux)
	v1.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: n,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited properly")
}
