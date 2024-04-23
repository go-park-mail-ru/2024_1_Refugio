package repository

import (
	"context"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	domain "mail/internal/microservice/models/domain_models"
	converters "mail/internal/microservice/models/repository_converters"
	database "mail/internal/microservice/models/repository_models"
	"mail/internal/pkg/logger"
	"math/rand"
	"time"
)

var requestIDContextKey interface{} = "requestID"

// UserRepository represents a repository for managing user data in the database.
type UserRepository struct {
	DB *sqlx.DB
}

// NewUserRepository creates a new instance of UserRepository with the given database connection.
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

// CheckPassword compares a password with a hash and returns true if they match, otherwise false.
type CheckPassword func(password, hash string) bool

// CheckPasswordHash compares a password with a hash and returns true if they match, otherwise false.
var CheckPasswordHash CheckPassword = func(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
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
func (r *UserRepository) GetAll(offset, limit int, ctx context.Context) ([]*domain.User, error) {
	query := "SELECT * FROM profile"

	var userModelsDb []database.User
	var err error
	if offset >= 0 && limit > 0 {
		query += " OFFSET $1 LIMIT $2"
		err = r.DB.Select(&userModelsDb, query, offset, limit)
	} else {
		err = r.DB.Select(&userModelsDb, query)
	}

	args := []interface{}{}
	start := time.Now()
	requestIDValue := ctx.Value(requestIDContextKey)
	requestIDValue = logger.GetRequestIDString(requestIDValue)
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, requestIDValue.(string), start, &err, args)

	if err != nil {
		return nil, err
	}

	usersCore := make([]*domain.User, 0, len(userModelsDb))
	for _, userModelDb := range userModelsDb {
		usersCore = append(usersCore, converters.UserConvertDbInCore(userModelDb))
	}

	return usersCore, nil
}

// GetByID returns the user by its unique identifier.
func (r *UserRepository) GetByID(id uint32, ctx context.Context) (*domain.User, error) {
	query := "SELECT * FROM profile WHERE id = $1"

	var userModelDb database.User
	err := r.DB.Get(&userModelDb, query, id)

	args := []interface{}{id}
	start := time.Now()
	requestIDValue := ctx.Value(requestIDContextKey)
	requestIDValue = logger.GetRequestIDString(requestIDValue)
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, requestIDValue.(string), start, &err, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}

		return nil, err
	}

	return converters.UserConvertDbInCore(userModelDb), nil
}

// GetUserByLogin returns the user by login.
func (r *UserRepository) GetUserByLogin(login, password string, ctx context.Context) (*domain.User, error) {
	query := "SELECT * FROM profile WHERE login = $1"

	var userModelDb database.User
	err := r.DB.Get(&userModelDb, query, login)

	args := []interface{}{login}
	start := time.Now()
	requestIDValue := ctx.Value(requestIDContextKey)
	requestIDValue = logger.GetRequestIDString(requestIDValue)
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, requestIDValue.(string), start, &err, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with login %s not found", login)
		}

		return nil, err
	}

	if CheckPasswordHash(password, userModelDb.Password) {
		return converters.UserConvertDbInCore(userModelDb), nil
	} else {
		err = fmt.Errorf("user with login %s not found", login)
		return nil, err
	}
}

// Add adds a new user to the storage and returns its assigned unique identifier.
func (r *UserRepository) Add(userModelCore *domain.User, ctx context.Context) (*domain.User, error) {
	query := `
		INSERT INTO profile (login, password, firstname, surname, patronymic, gender, birthday, registration_date, avatar_id, phone_number, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	userModelDb := converters.UserConvertCoreInDb(*userModelCore)

	password, status := HashPassword(userModelDb.Password)
	if !status {
		return &domain.User{}, fmt.Errorf("user with login %s fail", userModelDb.Login)
	}

	userModelDb.Password = password
	_, err := r.DB.Exec(query, userModelDb.Login, userModelDb.Password, userModelDb.FirstName, userModelDb.Surname, userModelDb.Patronymic, userModelDb.Gender, userModelDb.Birthday, time.Now(), userModelDb.AvatarID, userModelDb.PhoneNumber, userModelDb.Description)

	args := []interface{}{userModelDb.Login, userModelDb.FirstName, userModelDb.Surname, userModelDb.Patronymic, userModelDb.Gender, userModelDb.Birthday, time.Now(), userModelDb.AvatarID, userModelDb.PhoneNumber, userModelDb.Description}
	start := time.Now()
	requestIDValue := ctx.Value(requestIDContextKey)
	requestIDValue = logger.GetRequestIDString(requestIDValue)
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, requestIDValue.(string), start, &err, args)

	if err != nil {
		return &domain.User{}, fmt.Errorf("user with login %s fail", userModelDb.Login)
	}

	return userModelCore, nil
}

// Update updates the information of a user in the storage based on the provided new user.
func (r *UserRepository) Update(newUserCore *domain.User, ctx context.Context) (bool, error) {
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
	requestIDValue := ctx.Value(requestIDContextKey)
	requestIDValue = logger.GetRequestIDString(requestIDValue)
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, requestIDValue.(string), start, &err, args)

	if err != nil {
		return false, fmt.Errorf("failed to update user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		err = fmt.Errorf("user with id %d not found", newUserDb.ID)
		return false, err
	}

	return true, nil
}

// Delete removes the user from the storage by its unique identifier.
func (r *UserRepository) Delete(id uint32, ctx context.Context) (bool, error) {
	query := "DELETE FROM profile WHERE id = $1"

	result, err := r.DB.Exec(query, id)

	args := []interface{}{id}
	start := time.Now()
	requestIDValue := ctx.Value(requestIDContextKey)
	requestIDValue = logger.GetRequestIDString(requestIDValue)
	defer ctx.Value("logger").(*logger.LogrusLogger).DbLog(query, requestIDValue.(string), start, &err, args)

	if err != nil {
		return false, fmt.Errorf("failed to delete user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		err = fmt.Errorf("user with id %d not found", id)
		return false, err
	}

	return true, nil
}