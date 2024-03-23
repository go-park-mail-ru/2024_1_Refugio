package session

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionsManager_Check(t *testing.T) {
	// Тест на отсутствие сессии в запросе
	sm := NewSessionsManager()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = sm.Check(req)
	if err != ErrNoAuth {
		t.Fatalf("Ожидаемая ошибка %v, получено %v", ErrNoAuth, err)
	}

	// Тест на наличие валидной сессии
	userID := uint32(1)
	sess, err := sm.Create(httptest.NewRecorder(), userID)
	if err != nil {
		t.Fatal(err)
	}

	req.AddCookie(&http.Cookie{Name: "session_id", Value: sess.ID})
	gotSess, err := sm.Check(req)
	if err != nil {
		t.Fatal(err)
	}

	if gotSess != sess {
		t.Fatalf("Ожидалась сессия %v, получена %v", sess, gotSess)
	}
}

func TestSessionsManager_Create(t *testing.T) {
	sm := NewSessionsManager()
	userID := uint32(1)

	// Создание новой сессии
	recorder := httptest.NewRecorder()
	sess, err := sm.Create(recorder, userID)
	if err != nil {
		t.Fatal(err)
	}

	// Проверка, что сессия создана и куки установлены
	if len(recorder.Result().Cookies()) == 0 {
		t.Fatal("Куки не установлены")
	}

	// Проверка, что сессия добавлена в менеджер сессий
	sm.mu.RLock()
	_, ok := sm.data[sess.ID]
	sm.mu.RUnlock()

	if !ok {
		t.Fatal("Сессия не добавлена в менеджер сессий")
	}
}

func TestSessionsManager_DestroyCurrent(t *testing.T) {
	sm := NewSessionsManager()
	userID := uint32(1)
	sess, err := sm.Create(httptest.NewRecorder(), userID)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sess.ID})

	// Уничтожение текущей сессии
	err = sm.DestroyCurrent(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Проверка, что сессия удалена из менеджера сессий
	sm.mu.RLock()
	_, ok := sm.data[sess.ID]
	sm.mu.RUnlock()

	if ok {
		t.Fatal("Сессия не удалена из менеджера сессий")
	}
}
