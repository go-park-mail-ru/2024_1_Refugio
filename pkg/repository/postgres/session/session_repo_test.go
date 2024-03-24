package session

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	domain "mail/pkg/domain/models"
	"regexp"
	"testing"
	"time"
	// domain "mail/pkg/domain/models"
)

func TestCreateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := SessionRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	t.Run("Success", func(t *testing.T) {
		ID := uint32(1)
		userID := uint32(100)
		device := "mobile"
		lifeTime := 3600

		mock.ExpectExec(`INSERT INTO sessions`).WithArgs(ID, userID, device, sqlmock.AnyArg(), lifeTime).WillReturnResult(sqlmock.NewResult(1, 1))

		sessionID, err := repo.CreateSession(ID, userID, device, lifeTime)

		assert.NoError(t, err)
		assert.Equal(t, ID, sessionID)
	})

	t.Run("Error", func(t *testing.T) {
		ID := uint32(2)
		userID := uint32(101)
		device := "web"
		lifeTime := 7200

		mock.ExpectExec(`INSERT INTO sessions`).WithArgs(ID, userID, device, sqlmock.AnyArg(), lifeTime).WillReturnError(fmt.Errorf("failed to insert"))

		sessionID, err := repo.CreateSession(ID, userID, device, lifeTime)

		assert.Error(t, err)
		assert.Zero(t, sessionID)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetSessionByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := SessionRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	t.Run("Success", func(t *testing.T) {
		sessionID := uint32(1)
		expectedSession := &domain.Session{
			ID:           sessionID,
			UserID:       100,
			Device:       "mobile",
			LifeTime:     3600,
			CreationDate: time.Now(),
		}

		rows := sqlmock.NewRows([]string{"id", "user_id", "device", "creation_date", "life_time"}).
			AddRow(sessionID, 100, "mobile", time.Now(), 3600)

		query := `SELECT \* FROM sessions WHERE id = \$1`
		mock.ExpectQuery(query).WithArgs(sessionID).WillReturnRows(rows)

		session, err := repo.GetSessionByID(sessionID)

		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, expectedSession.ID, session.ID)
		// Добавьте дополнительные проверки по мере необходимости
	})

	t.Run("SessionNotFound", func(t *testing.T) {
		sessionID := uint32(2)

		query := `SELECT \* FROM sessions WHERE id = \$1`
		mock.ExpectQuery(query).WithArgs(sessionID).WillReturnError(sql.ErrNoRows)

		session, err := repo.GetSessionByID(sessionID)

		assert.Error(t, err)
		assert.Nil(t, session)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteSessionByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := SessionRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	t.Run("Success", func(t *testing.T) {
		sessionID := uint32(1)

		query := `DELETE FROM sessions WHERE id = \$1`
		mock.ExpectExec(query).WithArgs(sessionID).WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteSessionByID(sessionID)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		sessionID := uint32(2)

		query := `DELETE FROM sessions WHERE id = \$1`
		mock.ExpectExec(query).WithArgs(sessionID).WillReturnError(fmt.Errorf("failed to delete session"))

		err := repo.DeleteSessionByID(sessionID)

		assert.Error(t, err)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteExpiredSessions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := SessionRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	t.Run("Success", func(t *testing.T) {
		queryPattern := regexp.QuoteMeta(`DELETE FROM sessions WHERE creation_date + life_time * interval '1 second' < now()`)
		mock.ExpectExec(queryPattern).WillReturnResult(sqlmock.NewResult(0, 3)) // Assuming 3 expired sessions were deleted

		err := repo.DeleteExpiredSessions()

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		queryPattern := regexp.QuoteMeta(`DELETE FROM sessions WHERE creation_date + life_time * interval '1 second' < now()`)
		mock.ExpectExec(queryPattern).WillReturnError(fmt.Errorf("failed to delete expired sessions"))

		err := repo.DeleteExpiredSessions()

		assert.Error(t, err)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
