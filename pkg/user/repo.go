package user

import (
	"errors"
	"sync"
)

// UserMemoryRepository is an in-memory implementation of UserRepository.
type UserMemoryRepository struct {
	users map[uint32]*User
	mutex sync.RWMutex
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

// GetAll returns all users from the storage.
func (repo *UserMemoryRepository) GetAll() ([]*User, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	users := make([]*User, 0, len(repo.users))
	for _, user := range repo.users {
		users = append(users, user)
	}

	return users, nil
}

// GetByID returns the user by its unique identifier.
func (repo *UserMemoryRepository) GetByID(id uint32) (*User, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	user, exists := repo.users[id]
	if !exists {
		return nil, errors.New("user not found")
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
		return false, errors.New("user not found")
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
		return false, errors.New("user not found")
	}

	delete(repo.users, id)
	return true, nil
}
