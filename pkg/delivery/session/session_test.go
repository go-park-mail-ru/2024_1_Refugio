package session

import (
	"context"
	"testing"
)

func TestSession(t *testing.T) {
	userID := uint32(123)
	sess := NewSession(userID)

	// Тест наличия сессии в контексте
	ctx := ContextWithSession(context.Background(), sess)
	retSess, err := SessionFromContext(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retSess.ID != sess.ID || retSess.UserID != sess.UserID {
		t.Fatalf("Expected session %v, got %v", sess, retSess)
	}

	// Тест отсутствия сессии в пустом контексте
	emptyCtx := context.Background()
	_, err = SessionFromContext(emptyCtx)
	if err != ErrNoAuth {
		t.Fatalf("Expected error %v, got %v", ErrNoAuth, err)
	}
}
