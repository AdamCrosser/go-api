package authentication

import (
	"fmt"
	"net/http"
	"strings"
)

func Verifier(validAPIKeys map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header

			authHeaderString := headers.Get("Authorization")

			if authHeaderString == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Authorization header should include an API key"))
				return
			}

			parts := strings.Split(authHeaderString, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Invalid Authorization header format. It should be in the format: Bearer $APITOKEN"))
				return
			}

			reqAPIKey := parts[1]
			_, found := validAPIKeys[reqAPIKey]
			if !found {
				fmt.Printf("User attempted to authenticate with an invalid API key, key: %s\n", reqAPIKey)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("An invalid API key was specified in the Authorization header"))
				return
			}

			fmt.Printf("User successfully authenticated with API key, key: %s\n", reqAPIKey)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(hfn)
	}
}
