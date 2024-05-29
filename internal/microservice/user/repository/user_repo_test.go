package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mail/internal/pkg/logger"
	"mail/internal/pkg/utils/constants"

	domain "mail/internal/microservice/models/domain_models"
)

func GetCTX() context.Context {
	ctx := context.WithValue(context.Background(), interface{}(string(constants.LoggerKey)), logger.InitializationBdLog(nil))
	ctx2 := context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), []string{"testID"})

	return ctx2
}

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

	ctx := GetCTX()

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

		users, err := repo.GetAll(-1, -1, ctx)
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

		users, err := repo.GetAll(1, 2, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT * FROM profile").WillReturnError(sql.ErrNoRows)

		users, err := repo.GetAll(-1, -1, ctx)
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

	ctx := GetCTX()

	t.Run("UserExists", func(t *testing.T) {
		expectedUser := &domain.User{
			ID:          1,
			Login:       "testuser",
			FirstName:   "John",
			Surname:     "Doe",
			Patronymic:  "",
			Gender:      "male",
			Birthday:    time.Now(),
			AvatarID:    "1",
			PhoneNumber: "123456789",
			Description: "Test user",
		}

		rows := sqlmock.NewRows([]string{"id", "login", "firstname", "surname", "patronymic", "gender", "birthday", "avatar", "phone_number", "description"}).
			AddRow(expectedUser.ID, expectedUser.Login, expectedUser.FirstName, expectedUser.Surname, expectedUser.Patronymic, expectedUser.Gender, expectedUser.Birthday, expectedUser.AvatarID, expectedUser.PhoneNumber, expectedUser.Description)

		mock.ExpectQuery(`SELECT p\.id, p\.login, p\.firstname, p\.surname, p\.patronymic, p\.gender, p\.birthday, f\.file_id AS avatar, p\.phone_number, p\.description FROM profile p LEFT JOIN file f ON p\.avatar_id = f\.id WHERE p\.id = \$1`).WithArgs(expectedUser.ID).WillReturnRows(rows)

		user, err := repo.GetByID(expectedUser.ID, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT p\.id, p\.login, p\.firstname, p\.surname, p\.patronymic, p\.gender, p\.birthday, f\.file_id AS avatar, p\.phone_number, p\.description FROM profile p LEFT JOIN file f ON p\.avatar_id = f\.id WHERE p\.id = \$1`).WithArgs(1).WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByID(1, ctx)
		assert.Nil(t, user)
		assert.Error(t, err)
		expectedErrorMessage := fmt.Sprintf("user with id %d not found", 1)
		assert.EqualError(t, err, expectedErrorMessage)
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

	ctx := GetCTX()

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
			WithArgs(newUser.FirstName, newUser.Surname, newUser.Patronymic, newUser.Gender, newUser.Birthday, newUser.PhoneNumber, newUser.Description, newUser.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		updated, err := repo.Update(newUser, ctx)

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
			WithArgs(newUser.FirstName, newUser.Surname, newUser.Patronymic, newUser.Gender, newUser.Birthday, newUser.PhoneNumber, newUser.Description, newUser.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		updated, err := repo.Update(newUser, ctx)

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
			WithArgs(newUser.FirstName, newUser.Surname, newUser.Patronymic, newUser.Gender, newUser.Birthday, newUser.PhoneNumber, newUser.Description, newUser.ID).
			WillReturnError(fmt.Errorf("database error"))

		updated, err := repo.Update(newUser, ctx)

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

	ctx := GetCTX()

	t.Run("UserDeletedSuccessfully", func(t *testing.T) {
		userID := uint32(1)

		mock.ExpectExec(`DELETE FROM profile`).WithArgs(userID).WillReturnResult(sqlmock.NewResult(0, 1))

		deleted, err := repo.Delete(userID, ctx)

		assert.NoError(t, err)
		assert.True(t, deleted)
	})

	t.Run("UserDeleteFailedNoRowsAffected", func(t *testing.T) {
		userID := uint32(2)

		mock.ExpectExec(`DELETE FROM profile`).WithArgs(userID).WillReturnResult(sqlmock.NewResult(0, 0))

		deleted, err := repo.Delete(userID, ctx)

		assert.Error(t, err)
		assert.False(t, deleted)
	})

	t.Run("UserDeleteFailedDBError", func(t *testing.T) {
		userID := uint32(3)

		mock.ExpectExec(`DELETE FROM profile`).WithArgs(userID).WillReturnError(fmt.Errorf("database error"))

		deleted, err := repo.Delete(userID, ctx)

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

	login := "testuser"
	password := "password"

	ctx := GetCTX()

	t.Run("UserFound", func(t *testing.T) {
		expectedUser := &domain.User{
			ID:        1,
			Login:     login,
			FirstName: "John",
			Surname:   "Doe",
			Birthday:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		}

		rows := sqlmock.NewRows([]string{"id", "login", "password_hash", "firstname", "surname", "patronymic", "gender", "birthday", "phone_number", "description"}).
			AddRow(expectedUser.ID, expectedUser.Login, "hashed_password", expectedUser.FirstName, expectedUser.Surname, "", "", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), "", "")

		mock.ExpectQuery(`SELECT p.id, p.login, p.password_hash, p.firstname, p.surname, p.patronymic, p.gender, p.birthday, p.phone_number, p.description FROM profile p WHERE login = \$1`).WithArgs(login).WillReturnRows(rows)

		user, err := repo.GetUserByLogin(login, password, ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT p.id, p.login, p.password_hash, p.firstname, p.surname, p.patronymic, p.gender, p.birthday, p.phone_number, p.description FROM profile p WHERE login = \$1`).WithArgs(login).WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByLogin(login, password, ctx)

		assert.Error(t, err)
		assert.Nil(t, user)
		expectedErrorMessage := fmt.Sprintf("user with login %s not found", login)
		assert.EqualError(t, err, expectedErrorMessage)
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		mockCheckPassword := func(password, hash string) bool {
			return false
		}

		originalCheckPassword := CheckPasswordHash
		defer func() { CheckPasswordHash = originalCheckPassword }()
		CheckPasswordHash = mockCheckPassword

		rows := sqlmock.NewRows([]string{"id", "login", "password_hash", "firstname", "surname", "patronymic", "gender", "birthday", "phone_number", "description"}).
			AddRow(1, "login", "hashed_password", "John", "Petr", "", "", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), "", "")

		mock.ExpectQuery(`SELECT p.id, p.login, p.password_hash, p.firstname, p.surname, p.patronymic, p.gender, p.birthday, p.phone_number, p.description FROM profile p WHERE login = \$1`).WithArgs(login).WillReturnRows(rows)

		user, err := repo.GetUserByLogin(login, password, ctx)

		assert.Error(t, err)
		assert.Nil(t, user)
		expectedErrorMessage := fmt.Sprintf("user with login %s not found", login)
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

	ctx := GetCTX()

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
			VKId:        1,
		}

		mock.ExpectExec(`INSERT INTO profile`).
			WithArgs(user.Login, user.Password, user.FirstName, user.Surname, user.Patronymic, user.Gender, user.Birthday, sqlmock.AnyArg(), user.PhoneNumber, user.Description, user.VKId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		createUser, err := repo.Add(user, ctx)
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
			WithArgs(user.Login, user.Password, user.FirstName, user.Surname, user.Patronymic, user.Gender, user.Birthday, sqlmock.AnyArg(), user.PhoneNumber, user.Description).
			WillReturnError(fmt.Errorf("failed to insert user"))

		var userFail *domain.User
		userRes, err := repo.Add(user, ctx)
		assert.Error(t, err)
		assert.Equal(t, userFail, userRes)
	})
}

func TestAddAvatar(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := UserRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()

	userID := uint32(1)
	fileID := "avatar123"
	fileType := "image/jpeg"

	t.Run("AvatarAddedSuccessfully", func(t *testing.T) {
		mock.ExpectExec(`UPDATE file SET file_id = \$1 FROM profile WHERE file.id = profile.avatar_id AND profile.id = \$2`).
			WithArgs(fileID, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		isAdded, err := repo.AddAvatar(userID, fileID, fileType, ctx)
		assert.NoError(t, err)
		assert.True(t, isAdded)
	})

	t.Run("FailedToAddAvatar", func(t *testing.T) {
		mock.ExpectExec(`UPDATE file SET file_id = \$1 FROM profile WHERE file.id = profile.avatar_id AND profile.id = \$2`).
			WithArgs(fileID, userID).
			WillReturnError(fmt.Errorf("failed to update avatar"))

		isAdded, err := repo.AddAvatar(userID, fileID, fileType, ctx)
		assert.Error(t, err)
		assert.False(t, isAdded)
	})
}

func TestDeleteAvatarByUserID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := UserRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()

	userID := uint32(1)

	t.Run("AvatarDeletedSuccessfully", func(t *testing.T) {
		mock.ExpectExec(`UPDATE file SET file_id = '' FROM profile WHERE file.id = profile.avatar_id AND profile.id = \$1`).
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteAvatarByUserID(userID, ctx)
		assert.NoError(t, err)
	})

	t.Run("FailedToDeleteAvatar", func(t *testing.T) {
		mock.ExpectExec(`UPDATE file SET file_id = '' FROM profile WHERE file.id = profile.avatar_id AND profile.id = \$1`).
			WithArgs(userID).
			WillReturnError(errors.New("failed to delete avatar"))

		err := repo.DeleteAvatarByUserID(userID, ctx)
		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("failed to delete user avatar: %v", "failed to delete avatar"))
	})
}

func TestInitAvatar(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := UserRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()

	userID := uint32(1)
	fileID := "avatar123"
	fileType := "image/jpeg"

	t.Run("AvatarInitializedSuccessfully", func(t *testing.T) {
		expectedQuery := `WITH inserted_file AS \( INSERT INTO file \(file_id, file_type\) VALUES \(\$1, \$2\) RETURNING id \) UPDATE profile SET avatar_id = \(SELECT id FROM inserted_file\) WHERE id = \$3`
		mock.ExpectExec(expectedQuery).
			WithArgs(fileID, fileType, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		success, err := repo.InitAvatar(userID, fileID, fileType, ctx)
		assert.NoError(t, err)
		assert.True(t, success)
	})

	t.Run("FailedToInitializeAvatar", func(t *testing.T) {
		expectedQuery := `WITH inserted_file AS \( INSERT INTO file \(file_id, file_type\) VALUES \(\$1, \$2\) RETURNING id \) UPDATE profile SET avatar_id = \(SELECT id FROM inserted_file\) WHERE id = \$3`
		mock.ExpectExec(expectedQuery).
			WithArgs(fileID, fileType, userID).
			WillReturnError(errors.New("failed to initialize avatar"))

		success, err := repo.InitAvatar(userID, fileID, fileType, ctx)
		assert.Error(t, err)
		assert.False(t, success)
		assert.EqualError(t, err, fmt.Sprintf("failed to add user avatar: %v", "failed to initialize avatar"))
	})
}
