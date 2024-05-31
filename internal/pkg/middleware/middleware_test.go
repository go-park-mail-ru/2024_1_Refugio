package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicMiddleware(t *testing.T) {
	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	t.Run("WithoutCookie", func(t *testing.T) {
		reqWithoutCookie := httptest.NewRequest("GET", "/api/v1/auth/verify-auth", nil)
		recWithoutCookie := httptest.NewRecorder()
		PanicMiddleware(fakeHandler).ServeHTTP(recWithoutCookie, reqWithoutCookie)
		assert.Equal(t, http.StatusOK, recWithoutCookie.Code)
	})

}

func TestAuthMiddleware(t *testing.T) {
	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("WithoutCookie", func(t *testing.T) {
		reqWithoutCookie := httptest.NewRequest("GET", "/api/v1/auth/verify-auth", nil)
		recWithoutCookie := httptest.NewRecorder()
		AuthMiddleware(fakeHandler).ServeHTTP(recWithoutCookie, reqWithoutCookie)
		assert.Equal(t, http.StatusUnauthorized, recWithoutCookie.Code)
	})

	t.Run("WithInvalidCookie", func(t *testing.T) {
		reqWithInvalidCookie := httptest.NewRequest("GET", "/api/v1/auth/verify-auth", nil)
		reqWithInvalidCookie.AddCookie(&http.Cookie{Name: "session_id", Value: ""})
		recWithInvalidCookie := httptest.NewRecorder()
		AuthMiddleware(fakeHandler).ServeHTTP(recWithInvalidCookie, reqWithInvalidCookie)
		assert.Equal(t, http.StatusUnauthorized, recWithInvalidCookie.Code)
	})
}
