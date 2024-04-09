package user

import (
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"mail/pkg/domain/logger"
	"math/rand"
	"time"

	"mail/pkg/repository/converters"

	domain "mail/pkg/domain/models"
	database "mail/pkg/repository/models"
)

var Logger = logger.InitializationBdLog()

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// PasswordHasher represents the password hashing function.
type PasswordHasher func(password string) (string, bool)

// HashPassword takes a plaintext password as input and returns its bcrypt hash.
var HashPassword PasswordHasher = func(password string) (string, bool) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", false
	}

	return string(bytes), true
}

// CheckPasswordHash compares a password with a hash and returns true if they match, otherwise false.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

// ComparingUserObjects compares two user objects by comparing.
// If all fields match, the function returns true, otherwise false.
func ComparingUserObjects(object1, object2 domain.User) bool {
	userDb1 := converters.UserConvertCoreInDb(object1)
	userDb2 := converters.UserConvertCoreInDb(object2)

	if userDb1.ID == userDb2.ID &&
		userDb1.FirstName == userDb2.FirstName &&
		userDb1.Surname == userDb2.Surname &&
		userDb1.Patronymic == userDb2.Patronymic &&
		userDb1.Gender == userDb2.Gender &&
		userDb1.Birthday == userDb2.Birthday &&
		userDb1.Login == userDb2.Login &&
		CheckPasswordHash(userDb2.Password, userDb1.Password) {
		return true
	}

	return false
}

// RandomIDGenerator represents a function that generates random uint32 identifiers.
type RandomIDGenerator func() uint32

// GenerateRandomID generates a random uint32 identifier.
var GenerateRandomID RandomIDGenerator = func() uint32 {
	randBytes := make([]byte, 4)
	_, err := rand.Read(randBytes)
	if err != nil {
		// В случае ошибки вернуть ноль, но это можно обработать в вашем приложении по-разному.
		return 0
	}

	randID := binary.BigEndian.Uint32(randBytes)
	return randID / 100
}

// GetAll returns all users from the storage.
func (r *UserRepository) GetAll(offset, limit int, requestID string) ([]*domain.User, error) {
	query := "SELECT * FROM profile"
	args := []interface{}{}

	var userModelsDb []database.User
	var err error
	start := time.Now()
	if offset >= 0 && limit > 0 {
		query += " OFFSET $1 LIMIT $2"
		err = r.DB.Select(&userModelsDb, query, offset, limit)
	} else {
		err = r.DB.Select(&userModelsDb, query)
	}

	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return nil, err
	}

	usersCore := make([]*domain.User, 0, len(userModelsDb))
	for _, userModelDb := range userModelsDb {
		usersCore = append(usersCore, converters.UserConvertDbInCore(userModelDb))
	}

	Logger.DbLog(query, requestID, 200, start, nil, args)
	return usersCore, nil
}

// GetByID returns the user by its unique identifier.
func (r *UserRepository) GetByID(id uint32, requestID string) (*domain.User, error) {
	query := "SELECT * FROM profile WHERE id = $1"
	args := []interface{}{id}

	var userModelDb database.User
	start := time.Now()
	err := r.DB.Get(&userModelDb, query, id)
	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err
	}

	Logger.DbLog(query, requestID, 200, start, nil, args)
	return converters.UserConvertDbInCore(userModelDb), nil
}

// GetUserByLogin returns the user by login.
func (r *UserRepository) GetUserByLogin(login, password, requestID string) (*domain.User, error) {
	query := "SELECT * FROM profile WHERE login = $1"
	args := []interface{}{login}

	var userModelDb database.User
	start := time.Now()
	err := r.DB.Get(&userModelDb, query, login)
	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with login %s not found", login)
		}
		return nil, err
	}

	if CheckPasswordHash(password, userModelDb.Password) {
		Logger.DbLog(query, requestID, 200, start, nil, args)
		return converters.UserConvertDbInCore(userModelDb), nil
	} else {
		Logger.DbLog(query, requestID, 500, start, fmt.Errorf("user with login %s not found", login), args)
		return nil, fmt.Errorf("user with login %s not found", login)
	}
}

// Add adds a new user to the storage and returns its assigned unique identifier.
func (r *UserRepository) Add(userModelCore *domain.User, requestID string) (*domain.User, error) {
	query := `
		INSERT INTO profile (login, password_hash, firstname, surname, patronymic, gender, birthday, registration_date, avatar_id, phone_number, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	userModelDb := converters.UserConvertCoreInDb(*userModelCore)
	args := []interface{}{userModelDb.Login, userModelDb.Password, userModelDb.FirstName, userModelDb.Surname, userModelDb.Patronymic, userModelDb.Gender, userModelDb.Birthday, time.Now(), userModelDb.AvatarID, userModelDb.PhoneNumber, userModelDb.Description}
	password, status := HashPassword(userModelDb.Password)
	if !status {
		Logger.DbLog(query, requestID, 500, time.Now(), fmt.Errorf("user with login %s fail", userModelDb.Login), args)
		return &domain.User{}, fmt.Errorf("user with login %s fail", userModelDb.Login)
	}
	userModelDb.Password = password
	start := time.Now()
	_, err := r.DB.Exec(query, userModelDb.Login, userModelDb.Password, userModelDb.FirstName, userModelDb.Surname, userModelDb.Patronymic, userModelDb.Gender, userModelDb.Birthday, time.Now(), userModelDb.AvatarID, userModelDb.PhoneNumber, userModelDb.Description)
	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return &domain.User{}, fmt.Errorf("user with login %s fail", userModelDb.Login)
	}

	Logger.DbLog(query, requestID, 200, start, nil, args)
	return userModelCore, nil
}

// Update updates the information of a user in the storage based on the provided new user.
func (r *UserRepository) Update(newUserCore *domain.User, requestID string) (bool, error) {
	newUserDb := converters.UserConvertCoreInDb(*newUserCore)

	query := `
        UPDATE profile
        SET
            firstname = $1,
            surname = $2,
            patronymic = $3,
            gender = $4,
            birthday = $5,
            avatar_id = $6,
            phone_number = $7,
            description = $8
        WHERE
            id = $9
    `
	args := []interface{}{
		newUserDb.FirstName,
		newUserDb.Surname,
		newUserDb.Patronymic,
		newUserDb.Gender,
		newUserDb.Birthday,
		newUserDb.AvatarID,
		newUserDb.PhoneNumber,
		newUserDb.Description,
		newUserDb.ID,
	}

	start := time.Now()
	result, err := r.DB.Exec(
		query,
		newUserDb.FirstName,
		newUserDb.Surname,
		newUserDb.Patronymic,
		newUserDb.Gender,
		newUserDb.Birthday,
		newUserDb.AvatarID,
		newUserDb.PhoneNumber,
		newUserDb.Description,
		newUserDb.ID,
	)
	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return false, fmt.Errorf("failed to update user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return false, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		Logger.DbLog(query, requestID, 500, start, fmt.Errorf("user with id %d not found", newUserDb.ID), args)
		return false, fmt.Errorf("user with id %d not found", newUserDb.ID)
	}

	Logger.DbLog(query, requestID, 200, start, nil, args)
	return true, nil
}

// Delete removes the user from the storage by its unique identifier.
func (r *UserRepository) Delete(id uint32, requestID string) (bool, error) {
	query := "DELETE FROM profile WHERE id = $1"
	args := []interface{}{id}

	start := time.Now()
	result, err := r.DB.Exec(query, id)
	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return false, fmt.Errorf("failed to delete user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		Logger.DbLog(query, requestID, 500, start, err, args)
		return false, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		Logger.DbLog(query, requestID, 500, start, fmt.Errorf("user with id %d not found", id), args)
		return false, fmt.Errorf("user with id %d not found", id)
	}

	Logger.DbLog(query, requestID, 200, start, nil, args)
	return true, nil
}
