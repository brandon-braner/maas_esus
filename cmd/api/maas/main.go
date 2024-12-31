package main

import (
	"log"
	"net/http"

	"github.com/brandonbraner/maas/pkg/http/middleware"
	"github.com/urfave/negroni"
)

func main() {
	mux := http.NewServeMux()
	n := negroni.Classic()
	n.Use(middleware.AuthMiddleware{})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	n.UseHandler(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: n,
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
