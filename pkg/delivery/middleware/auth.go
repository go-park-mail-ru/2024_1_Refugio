package middleware

import (
	"mail/pkg/delivery"
	"mail/pkg/delivery/session"
	"net/http"
)

// AuthMiddleware is a middleware to check user authentication using cookies.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := session.GlobalSeaaionManager.Check(r)
		if err != nil {
			delivery.HandleError(w, http.StatusUnauthorized, "Not Authorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}
