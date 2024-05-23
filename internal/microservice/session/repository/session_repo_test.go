package repository

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"mail/internal/pkg/logger"
	"mail/internal/pkg/utils/constants"

	domain "mail/internal/microservice/models/domain_models"
)

func GetCTX() context.Context {
	ctx := context.WithValue(context.Background(), interface{}(string(constants.LoggerKey)), logger.InitializationBdLog(nil))
	ctx2 := context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), []string{"testID"})

	return ctx2
}

func TestCreateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGenerateRandomID := func() string {
		return "10101010"
	}

	originalGenerateRandomID := SessionGenerateRandomID
	defer func() { SessionGenerateRandomID = originalGenerateRandomID }()
	SessionGenerateRandomID = mockGenerateRandomID

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := SessionRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()

	t.Run("Success", func(t *testing.T) {
		ID := "10101010"
		userID := uint32(100)
		device := "mobile"
		lifeTime := 3600
		csrfToken := "10101010"

		mock.ExpectExec(`INSERT INTO session`).WithArgs(ID, userID, sqlmock.AnyArg(), device, lifeTime, csrfToken).WillReturnResult(sqlmock.NewResult(1, 1))

		sessionID, err := repo.CreateSession(userID, device, lifeTime, ctx)

		assert.NoError(t, err)
		assert.Equal(t, ID, sessionID)
	})

	t.Run("Error", func(t *testing.T) {
		ID := "10101010"
		userID := uint32(101)
		device := "web"
		lifeTime := 7200
		csrfToken := "10101010"

		mock.ExpectExec(`INSERT INTO session`).WithArgs(ID, userID, sqlmock.AnyArg(), device, lifeTime, csrfToken).WillReturnError(fmt.Errorf("failed to insert"))

		sessionID, err := repo.CreateSession(userID, device, lifeTime, ctx)

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

	ctx := GetCTX()

	t.Run("Success", func(t *testing.T) {
		sessionID := "10101010"
		expectedSession := &domain.Session{
			ID:           sessionID,
			UserID:       100,
			Device:       "mobile",
			LifeTime:     3600,
			CreationDate: time.Now(),
			CsrfToken:    "10101010",
		}

		rows := sqlmock.NewRows([]string{"id", "profile_id", "device", "creation_date", "life_time", "csrf_token"}).
			AddRow(sessionID, 100, "mobile", time.Now(), 3600, "10101010")

		query := `SELECT \* FROM session WHERE id = \$1`
		mock.ExpectQuery(query).WithArgs(sessionID).WillReturnRows(rows)

		session, err := repo.GetSessionByID(sessionID, ctx)

		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, expectedSession.ID, session.ID)
	})

	t.Run("SessionNotFound", func(t *testing.T) {
		sessionID := "10101010"

		query := `SELECT \* FROM session WHERE id = \$1`
		mock.ExpectQuery(query).WithArgs(sessionID).WillReturnError(sql.ErrNoRows)

		session, err := repo.GetSessionByID(sessionID, ctx)

		assert.Error(t, err)
		assert.Nil(t, session)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetLoginBySessionID(t *testing.T) {
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

	ctx := GetCTX()

	t.Run("Success", func(t *testing.T) {
		sessionID := "10101010"
		expectedLogin := "test_user"

		rows := sqlmock.NewRows([]string{"login"}).AddRow(expectedLogin)

		query := `SELECT login FROM profile JOIN session ON session.profile_id = profile.id WHERE session.id = \$1`
		mock.ExpectQuery(query).WithArgs(sessionID).WillReturnRows(rows)

		login, err := repo.GetLoginBySessionID(sessionID, ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedLogin, login)
	})

	t.Run("Error", func(t *testing.T) {
		sessionID := "10101010"
		expectedError := fmt.Errorf("failed to get session")

		query := `SELECT login FROM profile JOIN session ON session.profile_id = profile.id WHERE session.id = \$1`
		mock.ExpectQuery(query).WithArgs(sessionID).WillReturnError(expectedError)

		login, err := repo.GetLoginBySessionID(sessionID, ctx)

		assert.Error(t, err)
		assert.Empty(t, login)
	})

	t.Run("EmptyResult", func(t *testing.T) {
		sessionID := "10101010"

		rows := sqlmock.NewRows([]string{"login"})

		query := `SELECT login FROM profile JOIN session ON session.profile_id = profile.id WHERE session.id = \$1`
		mock.ExpectQuery(query).WithArgs(sessionID).WillReturnRows(rows)

		login, err := repo.GetLoginBySessionID(sessionID, ctx)

		assert.Error(t, err)
		assert.Empty(t, login)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		sessionID := "10101010"
		expectedError := fmt.Errorf("failed to get session: database error")

		query := `SELECT login FROM profile JOIN session ON session.profile_id = profile.id WHERE session.id = \$1`
		mock.ExpectQuery(query).WithArgs(sessionID).WillReturnError(expectedError)

		login, err := repo.GetLoginBySessionID(sessionID, ctx)

		assert.Error(t, err)
		assert.Empty(t, login)
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

	ctx := GetCTX()

	t.Run("Success", func(t *testing.T) {
		sessionID := "10101010"

		query := `DELETE FROM session WHERE id = \$1`
		mock.ExpectExec(query).WithArgs(sessionID).WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteSessionByID(sessionID, ctx)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		sessionID := "10101010"

		query := `DELETE FROM session WHERE id = \$1`
		mock.ExpectExec(query).WithArgs(sessionID).WillReturnError(fmt.Errorf("failed to delete session"))

		err := repo.DeleteSessionByID(sessionID, ctx)

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

	ctx := GetCTX()

	t.Run("Success", func(t *testing.T) {
		queryPattern := regexp.QuoteMeta(`DELETE FROM session WHERE creation_date + life_time * interval '1 second' < now()`)
		mock.ExpectExec(queryPattern).WillReturnResult(sqlmock.NewResult(0, 3))

		err := repo.DeleteExpiredSessions(ctx)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		queryPattern := regexp.QuoteMeta(`DELETE FROM session WHERE creation_date + life_time * interval '1 second' < now()`)
		mock.ExpectExec(queryPattern).WillReturnError(fmt.Errorf("failed to delete expired sessions"))

		err := repo.DeleteExpiredSessions(ctx)

		assert.Error(t, err)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetProfileIDBySessionID(t *testing.T) {
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

	ctx := GetCTX()

	t.Run("Success", func(t *testing.T) {
		sessionID := "10101010"
		query := `SELECT profile_id FROM session WHERE session.id = \$1`
		rows := sqlmock.NewRows([]string{"profile_id"}).AddRow(42)
		mock.ExpectQuery(query).WithArgs(sessionID).WillReturnRows(rows)

		profileID, err := repo.GetProfileIDBySessionID(sessionID, ctx)

		assert.NoError(t, err)
		assert.Equal(t, uint32(42), profileID)
	})

	t.Run("Error", func(t *testing.T) {
		sessionID := "10101010"
		query := `SELECT profile_id FROM session WHERE session.id = \$1`
		mock.ExpectQuery(query).WithArgs(sessionID).WillReturnError(fmt.Errorf("failed to retrieve profile ID"))

		_, err := repo.GetProfileIDBySessionID(sessionID, ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get user id")
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
