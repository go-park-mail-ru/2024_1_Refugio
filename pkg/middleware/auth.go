package middleware

import (
	"net/http"
)

// AuthMiddleware is a middleware to check user authentication using cookies.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")

		// If no cookie or error in retrieving, redirect to /login.
		if err != nil || cookie.Value == "" {
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusFound)
			return
		}

		// Perform additional checks if needed, for example, validate the token.

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
