package user

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	userCore "mail/pkg/domain/models"
	"mail/pkg/repository/converters"
	"mail/pkg/repository/models"
	"sync"
)

// UserMemoryRepository is an in-memory implementation of UserRepository.
type UserMemoryRepository struct {
	mutex sync.RWMutex
	users map[uint32]*models.User
}

// NewUserMemoryRepository creates a new instance of UserMemoryRepository.
func NewInMemoryUserRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		users: FakeUsers,
	}
}

// NewEmptyInMemoryUserRepository creates a new user repository in memory with an empty default user list.
func NewEmptyInMemoryUserRepository() *UserMemoryRepository {
	defaultUsers := map[uint32]*models.User{}
	return &UserMemoryRepository{
		users: defaultUsers,
	}
}

// HashPassword takes a plaintext password as input and returns its bcrypt hash.
func HashPassword(password string) (string, bool) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", false
	}

	hashedPassword := string(bytes)
	fmt.Println("Hashed Password:", hashedPassword)

	return string(bytes), true
}

// CheckPasswordHash compares a password with a hash and returns true if they match, otherwise false.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ComparingUserObjects compares two user objects by comparing their IDs, names, surnames, logins, and password hashes.
// If all fields match, the function returns true, otherwise false.
func ComparingUserObjects(object1, object2 userCore.User) bool {
	userDb1 := converters.UserConvertCoreInDb(object1)
	userDb2 := converters.UserConvertCoreInDb(object2)

	if userDb1.ID == userDb2.ID &&
		userDb1.FirstName == userDb2.FirstName &&
		userDb1.Surname == userDb2.Surname &&
		userDb1.Login == userDb2.Login &&
		CheckPasswordHash(userDb2.Password, userDb1.Password) {
		return true
	}

	return false
}

// GetAll returns all users from the storage.
func (repo *UserMemoryRepository) GetAll(offset, limit int) ([]*userCore.User, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	users := make([]*userCore.User, 0, len(repo.users))
	for i := 0; i < len(repo.users); i++ {
		users = append(users, converters.UserConvertDbInCore(*repo.users[uint32(i+1)]))
	}

	return users, nil
}

// GetByID returns the user by its unique identifier.
func (repo *UserMemoryRepository) GetByID(id uint32) (*userCore.User, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	user, exists := repo.users[id]
	if !exists {
		return nil, fmt.Errorf("User with id %d not found", id)
	}

	return converters.UserConvertDbInCore(*user), nil
}

// GetUserByLogin returns the user by login.
func (repo *UserMemoryRepository) GetUserByLogin(login string, password string) (*userCore.User, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	for _, u := range repo.users {
		if u.Login == login {
			if CheckPasswordHash(password, u.Password) {
				return converters.UserConvertDbInCore(*u), nil
			} else {
				return nil, fmt.Errorf("User with the username %s was not found", login)
			}
		}
	}

	return nil, fmt.Errorf("User with the username %s was not found", login)
}

// Add adds a new user to the storage and returns its assigned unique identifier.
func (repo *UserMemoryRepository) Add(user *userCore.User) (uint32, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	userDb := converters.UserConvertCoreInDb(*user)

	userID := uint32(len(repo.users) + 1)
	userDb.ID = userID
	var err bool
	userDb.Password, err = HashPassword(userDb.Password)
	if err == false {
		return userID, fmt.Errorf("Operation failed")
	}
	repo.users[userID] = userDb

	return userID, nil
}

// Update updates the information of a user in the storage based on the provided new user.
func (repo *UserMemoryRepository) Update(newUser *userCore.User) (bool, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	newUserDb := converters.UserConvertCoreInDb(*newUser)

	_, exists := repo.users[newUserDb.ID]
	if !exists {
		return false, fmt.Errorf("User with id %d not found", newUserDb.ID)
	}

	repo.users[newUserDb.ID] = newUserDb

	return true, nil
}

// Delete removes the user from the storage by its unique identifier.
func (repo *UserMemoryRepository) Delete(id uint32) (bool, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	_, exists := repo.users[id]
	if !exists {
		return false, fmt.Errorf("User with id %d not found", id)
	}

	delete(repo.users, id)

	return true, nil
}
