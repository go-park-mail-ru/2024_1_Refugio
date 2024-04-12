package repository

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	domain "mail/internal/models/domain_models"
	"testing"
	"time"
)

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := UserRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	t.Run("NoOffsetAndLimit", func(t *testing.T) {
		expectedUsers := []*domain.User{
			{ID: 1, FirstName: "User 1"},
			{ID: 2, FirstName: "User 2"},
			{ID: 3, FirstName: "User 3"},
		}

		rows := sqlmock.NewRows([]string{"id", "firstname"}).
			AddRow(1, "User 1").
			AddRow(2, "User 2").
			AddRow(3, "User 3")

		mock.ExpectQuery(`SELECT \* FROM profile`).WillReturnRows(rows)

		users, err := repo.GetAll(-1, -1, "requestID")
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
	})

	t.Run("WithOffsetAndLimit", func(t *testing.T) {
		expectedUsers := []*domain.User{
			{ID: 2, FirstName: "User 2"},
			{ID: 3, FirstName: "User 3"},
		}

		rows := sqlmock.NewRows([]string{"id", "firstname"}).
			AddRow(2, "User 2").
			AddRow(3, "User 3")

		mock.ExpectQuery(`SELECT \* FROM profile OFFSET \$1 LIMIT \$2`).WithArgs(1, 2).WillReturnRows(rows)

		users, err := repo.GetAll(1, 2, "requestID")
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT * FROM profile").WillReturnError(sql.ErrNoRows)

		users, err := repo.GetAll(-1, -1, "requestID")
		assert.Error(t, err)
		assert.Nil(t, users)
	})
}

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := UserRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	t.Run("UserExists", func(t *testing.T) {
		expectedUser := &domain.User{
			ID:        1,
			FirstName: "John",
			Surname:   "Doe",
		}

		rows := sqlmock.NewRows([]string{"id", "firstname", "surname"}).
			AddRow(expectedUser.ID, expectedUser.FirstName, expectedUser.Surname)

		mock.ExpectQuery(`SELECT \* FROM profile WHERE id = \$1`).WithArgs(expectedUser.ID).WillReturnRows(rows)

		user, err := repo.GetByID(expectedUser.ID, "requestID")
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM profile WHERE id = \$1`).WithArgs(1).WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByID(1, "requestID")
		assert.Nil(t, user)
		assert.Error(t, err)
		expectedErrorMessage := fmt.Sprintf("user with id %d not found", 1)
		assert.EqualError(t, err, expectedErrorMessage)
	})
}

func TestAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHashPassword := func(password string) (string, bool) {
		return "1234", true
	}

	originalHashPassword := HashPassword
	defer func() { HashPassword = originalHashPassword }()
	HashPassword = mockHashPassword

	mockRandomIDGenerator := func() uint32 {
		return 1
	}

	originalRandomIDGenerator := GenerateRandomID
	defer func() { GenerateRandomID = originalRandomIDGenerator }()
	GenerateRandomID = mockRandomIDGenerator

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := UserRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	t.Run("UserAddedSuccessfully", func(t *testing.T) {
		user := &domain.User{
			Login:       "testuser",
			Password:    "1234",
			FirstName:   "John",
			Surname:     "Doe",
			Patronymic:  "Doe",
			Gender:      "male",
			Birthday:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			AvatarID:    "avatar123",
			PhoneNumber: "123456789",
			Description: "Тестовый пользователь",
		}

		mock.ExpectExec(`INSERT INTO profile`).
			WithArgs(user.Login, user.Password, user.FirstName, user.Surname, user.Patronymic, user.Gender, user.Birthday, sqlmock.AnyArg(), user.AvatarID, user.PhoneNumber, user.Description).
			WillReturnResult(sqlmock.NewResult(1, 1))

		createUser, err := repo.Add(user, "requestID")
		assert.NoError(t, err)
		assert.Equal(t, user, createUser)
	})

	t.Run("UserAddFailed", func(t *testing.T) {
		user := &domain.User{
			Login:       "testuser",
			Password:    "password",
			FirstName:   "John",
			Surname:     "Doe",
			Patronymic:  "Doe",
			Gender:      "male",
			Birthday:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			AvatarID:    "avatar123",
			PhoneNumber: "123456789",
			Description: "Test user",
		}

		mock.ExpectExec(`INSERT INTO profile`).
			WithArgs(user.Login, user.Password, user.FirstName, user.Surname, user.Patronymic, user.Gender, user.Birthday, sqlmock.AnyArg(), user.AvatarID, user.PhoneNumber, user.Description).
			WillReturnError(fmt.Errorf("failed to insert user"))

		userRes, err := repo.Add(user, "requestID")
		assert.Error(t, err)
		assert.Equal(t, &domain.User{}, userRes)
	})
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := UserRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	t.Run("UserUpdatedSuccessfully", func(t *testing.T) {
		newUser := &domain.User{
			ID:          1,
			FirstName:   "John",
			Surname:     "Doe",
			Patronymic:  "Doe",
			Gender:      "male",
			Birthday:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			AvatarID:    "new_avatar123",
			PhoneNumber: "987654321",
			Description: "Updated user",
		}

		mock.ExpectExec(`UPDATE profile`).
			WithArgs(newUser.FirstName, newUser.Surname, newUser.Patronymic, newUser.Gender, newUser.Birthday, newUser.AvatarID, newUser.PhoneNumber, newUser.Description, newUser.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		updated, err := repo.Update(newUser, "requestID")

		assert.NoError(t, err)
		assert.True(t, updated)
	})

	t.Run("UserUpdateFailedNoRowsAffected", func(t *testing.T) {
		newUser := &domain.User{
			ID:          2,
			FirstName:   "John",
			Surname:     "Doe",
			Patronymic:  "Doe",
			Gender:      "male",
			Birthday:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			AvatarID:    "new_avatar123",
			PhoneNumber: "987654321",
			Description: "Updated user",
		}

		mock.ExpectExec(`UPDATE profile`).
			WithArgs(newUser.FirstName, newUser.Surname, newUser.Patronymic, newUser.Gender, newUser.Birthday, newUser.AvatarID, newUser.PhoneNumber, newUser.Description, newUser.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		updated, err := repo.Update(newUser, "requestID")

		assert.Error(t, err)
		assert.False(t, updated)
	})

	t.Run("UserUpdateFailedDBError", func(t *testing.T) {
		newUser := &domain.User{
			ID:          3,
			FirstName:   "John",
			Surname:     "Doe",
			Patronymic:  "Doe",
			Gender:      "male",
			Birthday:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			AvatarID:    "new_avatar123",
			PhoneNumber: "987654321",
			Description: "Updated user",
		}

		mock.ExpectExec(`UPDATE profile`).
			WithArgs(newUser.FirstName, newUser.Surname, newUser.Patronymic, newUser.Gender, newUser.Birthday, newUser.AvatarID, newUser.PhoneNumber, newUser.Description, newUser.ID).
			WillReturnError(fmt.Errorf("database error"))

		updated, err := repo.Update(newUser, "requestID")

		assert.Error(t, err)
		assert.False(t, updated)
	})
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := UserRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	t.Run("UserDeletedSuccessfully", func(t *testing.T) {
		userID := uint32(1)

		mock.ExpectExec(`DELETE FROM profile`).WithArgs(userID).WillReturnResult(sqlmock.NewResult(0, 1))

		deleted, err := repo.Delete(userID, "requestID")

		assert.NoError(t, err)
		assert.True(t, deleted)
	})

	t.Run("UserDeleteFailedNoRowsAffected", func(t *testing.T) {
		userID := uint32(2)

		mock.ExpectExec(`DELETE FROM profile`).WithArgs(userID).WillReturnResult(sqlmock.NewResult(0, 0))

		deleted, err := repo.Delete(userID, "requestID")

		assert.Error(t, err)
		assert.False(t, deleted)
	})

	t.Run("UserDeleteFailedDBError", func(t *testing.T) {
		userID := uint32(3)

		mock.ExpectExec(`DELETE FROM profile`).WithArgs(userID).WillReturnError(fmt.Errorf("database error"))

		deleted, err := repo.Delete(userID, "requestID")

		assert.Error(t, err)
		assert.False(t, deleted)
	})
}

func TestGetUserByLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mockCheckPassword := func(password, hash string) bool {
		return true
	}

	originalCheckPassword := CheckPasswordHash
	defer func() { CheckPasswordHash = originalCheckPassword }()
	CheckPasswordHash = mockCheckPassword

	repo := UserRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	t.Run("UserFound", func(t *testing.T) {
		login := "testuser"
		password := "password"
		requestID := "requestID"

		expectedUser := &domain.User{
			ID:        1,
			Login:     login,
			FirstName: "John",
			Surname:   "Doe",
		}

		rows := sqlmock.NewRows([]string{"id", "login", "firstname", "surname", "password"}).
			AddRow(expectedUser.ID, expectedUser.Login, expectedUser.FirstName, expectedUser.Surname, "hashed_password")

		mock.ExpectQuery(`SELECT \* FROM profile WHERE login = \$1`).WithArgs(login).WillReturnRows(rows)

		user, err := repo.GetUserByLogin(login, password, requestID)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		login := "testuser"
		password := "password"
		requestID := "requestID"

		mock.ExpectQuery(`SELECT \* FROM profile WHERE login = \$1`).WithArgs(login).WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByLogin(login, password, requestID)

		assert.Error(t, err)
		assert.Nil(t, user)
		expectedErrorMessage := fmt.Sprintf("user with login %s not found", login)
		assert.EqualError(t, err, expectedErrorMessage)
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		login := "testuser"
		password := "invalid_password"
		requestID := "requestID"

		mockCheckPassword := func(password, hash string) bool {
			return false
		}

		originalCheckPassword := CheckPasswordHash
		defer func() { CheckPasswordHash = originalCheckPassword }()
		CheckPasswordHash = mockCheckPassword

		rows := sqlmock.NewRows([]string{"id", "login", "firstname", "surname", "password"}).
			AddRow(1, login, "John", "Doe", "hashed_password")

		mock.ExpectQuery(`SELECT \* FROM profile WHERE login = \$1`).WithArgs(login).WillReturnRows(rows)

		user, err := repo.GetUserByLogin(login, password, requestID)

		assert.Error(t, err)
		assert.Nil(t, user)
		expectedErrorMessage := fmt.Sprintf("user with login %s not found", login)
		assert.EqualError(t, err, expectedErrorMessage)
	})
}
