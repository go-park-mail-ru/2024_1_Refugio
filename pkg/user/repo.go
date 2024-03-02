package user

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

// UserMemoryRepository is an in-memory implementation of UserRepository.
type UserMemoryRepository struct {
	mutex sync.RWMutex
	users map[uint32]*User
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// NewUserMemoryRepository creates a new instance of UserMemoryRepository.
func NewInMemoryUserRepository() *UserMemoryRepository {
	defaultUsers := map[uint32]*User{
		1: {ID: 1, Name: "Sergey", Surname: "Fedasov", Login: "sergey@mail.ru", Password: "1234"},
		2: {ID: 2, Name: "Ivan", Surname: "Karpov", Login: "ivan@mail.ru", Password: "1234"},
		3: {ID: 3, Name: "Max", Surname: "Frelich", Login: "max@mail.ru", Password: "love"},
	}

	return &UserMemoryRepository{
		users: defaultUsers,
	}
}

func NewEmptyInMemoryUserRepository() *UserMemoryRepository {
	defaultUsers := map[uint32]*User{}
	return &UserMemoryRepository{
		users: defaultUsers,
	}
}

func ComparingUserObjects(object1, object2 User) bool {
	if object1.ID == object2.ID &&
		object1.Name == object2.Name &&
		object1.Surname == object2.Surname &&
		object1.Login == object2.Login &&
		CheckPasswordHash(object2.Password, object1.Password) {
		return true
	}
	return false
}

// GetAll returns all users from the storage.
func (repo *UserMemoryRepository) GetAll() ([]*User, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	users := make([]*User, 0, len(repo.users))
	for i := 0; i < len(repo.users); i++ {
		users = append(users, repo.users[uint32(i+1)])
	}

	return users, nil
}

// GetByID returns the user by its unique identifier.
func (repo *UserMemoryRepository) GetByID(id uint32) (*User, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	user, exists := repo.users[id]
	if !exists {
		return nil, fmt.Errorf("User with id %d not found", id)
	}

	return user, nil
}

// Add adds a new user to the storage and returns its assigned unique identifier.
func (repo *UserMemoryRepository) Add(user *User) (uint32, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	userID := uint32(len(repo.users) + 1)
	user.ID = userID
	repo.users[userID] = user

	return userID, nil
}

// Update updates the information of a user in the storage based on the provided new user.
func (repo *UserMemoryRepository) Update(newUser *User) (bool, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	_, exists := repo.users[newUser.ID]
	if !exists {
		return false, fmt.Errorf("User with id %d not found", newUser.ID)
	}

	repo.users[newUser.ID] = newUser

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
