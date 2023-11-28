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

	// Assigns a unique request identiifer to each request which can help with debugging
	r.Use(middleware.RequestID)

	//
	// TODO: Using the RealIP middleware could be dangerous if the application isn't deployed
	//       behind a CDN or load balancer
	//

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	//
	// Be extra fussy about the Content-Type header as a defense in-depth mechanism. Although, a
	// lot of this only matters if the API is going to be used by browser-based cilents.
	//
	// https://portswigger.net/web-security/graphql/lab-graphql-csrf-via-graphql-api
	// https://directdefense.com/csrf-in-the-age-of-json
	//
	// curl -H "Content-Type: application/x-www-form-urlencoded" -d "test=test "-X POST -v -H "Authorization: Bearer gae_ffd7bef1" http://localhost:8080/api/v1/private/test
	//

	r.Use(middleware.AllowContentType("application/json"))

	//
	// Be extra paranoid about caching of API responses by setting no cache
	// headers in the responses for all API requests
	//

	r.Use(middleware.NoCache)

	// NOTE: A list of available middlewares as part of the go-chi framework
	// https://github.com/go-chi/chi/tree/master/middleware

	// NOTE: A middleware that we could leverage for cross origin resource sharing
	// https://github.com/go-chi/cors

	// NOTE: This may be useful in terms of rate limiting for API endpoints
	// https://github.com/go-chi/httprate
	// https://github.com/go-chi/httprate-redis

	r.Route("/api/v1/private/", func(r chi.Router) {

		//
		// TODO: This is just for testing the authorization functionality and you shouldn't
		//       hardcode secrets in source code like this for a production application
		//

		// The example gae_ prefix is used to make it easier to identify API tokens hardcoded
		// in the source code of the application
		//
		// https://github.blog/2021-04-05-behind-githubs-new-authentication-token-formats/

		apiKeys := map[string]string{
			"gae_ffd7bef1": "admin",
			"gae_98e8d5f4": "reguser",
		}

		// TODO: It's probably better to pass a function to the authentication middleware
		//       so that we can decouple the underlying datastore for tokens from the
		//       implementation of the authorization middleware
		//
		// NOTE: Example below of how Go lets you define function as a first-class
		//       data-type/value
		// x := func(name string) string {
		//	return "Hello " + name
		// }
		// x("Adam")

		//
		// Performs validation of API keys using the Authorization: Bearer $APITOKEN format
		// curl -v -H "Authorization: Bearer gae_ffd7bef1" http://localhost:8080/api/v1/private/test
		//

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

	//
	// TODO: In scenarios where this API would be deployed behind a load balancer or other gateway
	// my recommendation would be to leverage only HTTP/2 for communication with backend systems
	// to reduce the risk of HTTP desync issues, etc. that occur with the ambiguity associated with
	// HTTP/1.1 request parsing
	// https://portswigger.net/research/http-desync-attacks-request-smuggling-reborn
	//

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Printf("An unknown error occurred while trying to listen on port 8080, err: %v\n", err)
	}
}
