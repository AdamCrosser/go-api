package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AdamCrosser/go-api/pkg/authentication"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP) // TODO: May be needed depending on if the API is deployed behind a proxy or other gateway
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1/private/", func(r chi.Router) {

		//
		// TODO: This is just for testing the authorization functionality and you shouldn't
		//       hardcode secrets in source code like this for a production application
		//

		apiKeys := map[string]string{
			"ffd7bef1": "admin",
			"98e8d5f4": "reguser",
		}

		r.Use(authentication.Verifier(apiKeys))

		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("{\"authorized\": true}"))
		})
	})

	r.Route("/api/v1/public/", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("{\"alive\": true}"))
		})
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Printf("An unknown error occurred while trying to listen on port 8080, err: %v\n", err)
	}
}
