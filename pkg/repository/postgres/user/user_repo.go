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
	"os"
	"time"

	"mail/pkg/repository/converters"

	domain "mail/pkg/domain/models"
	database "mail/pkg/repository/models"
)

// UserRepository represents a repository for managing user data in the database.
type UserRepository struct {
	DB *sqlx.DB
}

// Logger is an instance of a logger used for logging database operations.
var Logger = logger.InitializationEmptyLog()

// NewUserRepository creates a new instance of UserRepository with the given database connection.
func NewUserRepository(db *sqlx.DB) *UserRepository {
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile in user_repo" + "log.txt")
	}
	Logger = logger.InitializationBdLog(f)
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
func (r *UserRepository) GetAll(offset, limit int, requestID string) ([]*domain.User, error) {
	query := "SELECT * FROM profile"

	args := []interface{}{}
	start := time.Now()

	var userModelsDb []database.User
	var err error
	if offset >= 0 && limit > 0 {
		query += " OFFSET $1 LIMIT $2"
		err = r.DB.Select(&userModelsDb, query, offset, limit)
	} else {
		err = r.DB.Select(&userModelsDb, query)
	}
	defer Logger.DbLog(query, requestID, start, &err, args)

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
func (r *UserRepository) GetByID(id uint32, requestID string) (*domain.User, error) {
	query := "SELECT * FROM profile WHERE id = $1"

	args := []interface{}{id}
	start := time.Now()

	var userModelDb database.User
	err := r.DB.Get(&userModelDb, query, id)
	defer Logger.DbLog(query, requestID, start, &err, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}

		return nil, err
	}

	return converters.UserConvertDbInCore(userModelDb), nil
}

// GetUserByLogin returns the user by login.
func (r *UserRepository) GetUserByLogin(login, password, requestID string) (*domain.User, error) {
	query := "SELECT * FROM profile WHERE login = $1"

	args := []interface{}{login}
	start := time.Now()

	var userModelDb database.User
	err := r.DB.Get(&userModelDb, query, login)
	defer Logger.DbLog(query, requestID, start, &err, args)

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
func (r *UserRepository) Add(userModelCore *domain.User, requestID string) (*domain.User, error) {
	query := `
		INSERT INTO profile (login, password, firstname, surname, patronymic, gender, birthday, registration_date, avatar_id, phone_number, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	userModelDb := converters.UserConvertCoreInDb(*userModelCore)
	args := []interface{}{userModelDb.Login, userModelDb.Password, userModelDb.FirstName, userModelDb.Surname, userModelDb.Patronymic, userModelDb.Gender, userModelDb.Birthday, time.Now(), userModelDb.AvatarID, userModelDb.PhoneNumber, userModelDb.Description}
	start := time.Now()

	password, status := HashPassword(userModelDb.Password)
	if !status {
		return &domain.User{}, fmt.Errorf("user with login %s fail", userModelDb.Login)
	}

	userModelDb.Password = password
	_, err := r.DB.Exec(query, userModelDb.Login, userModelDb.Password, userModelDb.FirstName, userModelDb.Surname, userModelDb.Patronymic, userModelDb.Gender, userModelDb.Birthday, time.Now(), userModelDb.AvatarID, userModelDb.PhoneNumber, userModelDb.Description)
	defer Logger.DbLog(query, requestID, start, &err, args)

	if err != nil {
		return &domain.User{}, fmt.Errorf("user with login %s fail", userModelDb.Login)
	}

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
	defer Logger.DbLog(query, requestID, start, &err, args)

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
func (r *UserRepository) Delete(id uint32, requestID string) (bool, error) {
	query := "DELETE FROM profile WHERE id = $1"

	args := []interface{}{id}
	start := time.Now()

	result, err := r.DB.Exec(query, id)
	defer Logger.DbLog(query, requestID, start, &err, args)

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
