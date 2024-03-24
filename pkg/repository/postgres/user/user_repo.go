package user

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"time"

	"mail/pkg/repository/converters"

	domain "mail/pkg/domain/models"
	database "mail/pkg/repository/models"
)

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

/* // HashPassword takes a plaintext password as input and returns its bcrypt hash.
func HashPassword(password string) (string, bool) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", false
	}

	return string(bytes), true
}*/

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

// GetAll returns all users from the storage.
func (r *UserRepository) GetAll(offset, limit int) ([]*domain.User, error) {
	query := "SELECT * FROM users"

	var userModelsDb []database.User
	var err error
	if offset >= 0 && limit > 0 {
		query += " OFFSET $1 LIMIT $2"
		err = r.DB.Select(&userModelsDb, query, offset, limit)
	} else {
		err = r.DB.Select(&userModelsDb, query)
	}

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
func (r *UserRepository) GetByID(id uint32) (*domain.User, error) {
	query := "SELECT * FROM users WHERE id = $1"

	var userModelDb database.User
	err := r.DB.Get(&userModelDb, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err
	}

	return converters.UserConvertDbInCore(userModelDb), nil
}

// GetUserByLogin returns the user by login.
func (r *UserRepository) GetUserByLogin(login string, password string) (*domain.User, error) {
	query := "SELECT * FROM users WHERE login = $1"

	var userModelDb database.User
	err := r.DB.Get(&userModelDb, query, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with login %s not found", login)
		}
		return nil, err
	}

	if CheckPasswordHash(password, userModelDb.Password) {
		return converters.UserConvertDbInCore(userModelDb), nil
	} else {
		return nil, fmt.Errorf("user with login %s not found", login)
	}
}

// Add adds a new user to the storage and returns its assigned unique identifier.
func (r *UserRepository) Add(userModelCore *domain.User) (uint32, error) {
	query := `
		INSERT INTO users (login, password, firstname, surname, patronymic, gender, birthday, registration_date, avatar_id, phone_number, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`

	userModelDb := converters.UserConvertCoreInDb(*userModelCore)
	password, status := HashPassword(userModelDb.Password)
	if !status {
		return 0, fmt.Errorf("user with login %s fail", userModelDb.Login)
	}
	userModelDb.Password = password
	var userID uint32
	err := r.DB.QueryRow(query, userModelDb.Login, userModelDb.Password, userModelDb.FirstName, userModelDb.Surname, userModelDb.Patronymic, userModelDb.Gender, userModelDb.Birthday, time.Now(), userModelDb.AvatarID, userModelDb.PhoneNumber, userModelDb.Description).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// Update updates the information of a user in the storage based on the provided new user.
func (r *UserRepository) Update(newUserCore *domain.User) (bool, error) {
	newUserDb := converters.UserConvertCoreInDb(*newUserCore)

	query := `
        UPDATE users
        SET
            firstname = $1,
            surname = $2,
            patronymic = $3,
            gender = $4,
            birthday = $5,
            avatarid = $6,
            phonenumber = $7,
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
	if err != nil {
		return false, fmt.Errorf("failed to update user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return false, fmt.Errorf("user with id %d not found", newUserDb.ID)
	}

	return true, nil
}

// Delete removes the user from the storage by its unique identifier.
func (r *UserRepository) Delete(id uint32) (bool, error) {
	query := "DELETE FROM users WHERE id = $1"

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return false, fmt.Errorf("failed to delete user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return false, fmt.Errorf("user with id %d not found", id)
	}

	return true, nil
}
