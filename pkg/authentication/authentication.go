package authentication

import (
	"fmt"
	"net/http"
)

func Verifier(apiKeys map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header

			authHeaderVal := headers.Get("Authorization")

			if authHeaderVal == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Authorization header should include an API key"))
				return
			}

			apiKey, found := apiKeys[authHeaderVal]

			if !found {
				fmt.Printf("User attempted to authenticate with an invalid API key, key: %s\n", apiKey)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("An invalid API key was specified in the Authorization header"))
				return
			}

			fmt.Printf("User successfully authenticated with API key, key: %s\n", apiKey)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(hfn)
	}
}
