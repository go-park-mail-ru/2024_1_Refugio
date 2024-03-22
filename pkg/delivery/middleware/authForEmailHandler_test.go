package middleware

/*
import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	// Создаем фейковый хендлер, который всегда возвращает OK
	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Создаем фейковый запрос без куки
	reqWithoutCookie := httptest.NewRequest("GET", "/", nil)
	recWithoutCookie := httptest.NewRecorder()

	// Используем мидлвар с фейковым хендлером
	AuthMiddleware(fakeHandler).ServeHTTP(recWithoutCookie, reqWithoutCookie)

	// Проверяем, что получен редирект на /login
	assert.Equal(t, http.StatusFound, recWithoutCookie.Code)
	assert.Equal(t, "/login", recWithoutCookie.Header().Get("Location"))

	// Создаем фейковый запрос с некорректной кукой
	reqWithInvalidCookie := httptest.NewRequest("GET", "/", nil)
	reqWithInvalidCookie.AddCookie(&http.Cookie{Name: "session_id", Value: ""})
	recWithInvalidCookie := httptest.NewRecorder()

	// Используем мидлвар с фейковым хендлером
	AuthMiddleware(fakeHandler).ServeHTTP(recWithInvalidCookie, reqWithInvalidCookie)

	// Проверяем, что получен редирект на /login
	assert.Equal(t, http.StatusFound, recWithInvalidCookie.Code)
	assert.Equal(t, "/login", recWithInvalidCookie.Header().Get("Location"))

	// Создаем фейковый запрос с валидной кукой
	reqWithValidCookie := httptest.NewRequest("GET", "/", nil)
	reqWithValidCookie.AddCookie(&http.Cookie{Name: "session_id", Value: "valid_session_id"})
	recWithValidCookie := httptest.NewRecorder()

	// Используем мидлвар с фейковым хендлером
	AuthMiddleware(fakeHandler).ServeHTTP(recWithValidCookie, reqWithValidCookie)

	// Проверяем, что код ответа не изменился
	assert.Equal(t, http.StatusOK, recWithValidCookie.Code)
}
*/
