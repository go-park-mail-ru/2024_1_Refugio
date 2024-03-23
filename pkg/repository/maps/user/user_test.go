package user

import (
	"github.com/stretchr/testify/assert"
	userCore "mail/pkg/domain/models"
	"testing"
)

func TestGetAll(t *testing.T) {
	// Создаем репозиторий в памяти
	repo := NewInMemoryUserRepository()

	// Получаем всех пользователей из репозитория
	users, err := repo.GetAll()

	// Проверяем, что ошибка не возникла
	assert.NoError(t, err)
	// Проверяем, что список пользователей не пустой
	assert.NotEmpty(t, users)
}

func TestAddUser(t *testing.T) {
	// Создаем репозиторий в памяти
	repo := NewEmptyInMemoryUserRepository()

	// Создаем нового пользователя
	newUser := &userCore.User{
		Name:     "John",
		Surname:  "Doe",
		Login:    "john_doe",
		Password: "secret123",
	}

	// Добавляем пользователя в репозиторий
	userID, err := repo.Add(newUser)

	// Проверяем, что ошибка не возникла
	assert.NoError(t, err)
	// Проверяем, что ID пользователя присвоен
	assert.NotZero(t, userID)

	// Получаем пользователя из репозитория
	addedUser, _ := repo.GetByID(userID)

	// Проверяем, что полученный пользователь соответствует созданному
	assert.True(t, newUser.Name == addedUser.Name &&
		newUser.Surname == addedUser.Surname &&
		newUser.Login == addedUser.Login &&
		newUser.AvatarID == addedUser.AvatarID)
}

func TestUpdateUser(t *testing.T) {
	// Создаем репозиторий в памяти
	repo := NewEmptyInMemoryUserRepository()

	// Создаем нового пользователя
	newUser := &userCore.User{
		Name:     "John",
		Surname:  "Doe",
		Login:    "john_doe",
		Password: "secret123",
	}

	// Добавляем пользователя в репозиторий
	userID, _ := repo.Add(newUser)

	// Меняем данные пользователя
	newUserData := &userCore.User{
		ID:       userID,
		Name:     "Jane",
		Surname:  "Doe",
		Login:    "jane_doe",
		Password: "newsecret456",
	}

	// Обновляем пользователя в репозитории
	updated, err := repo.Update(newUserData)

	// Проверяем, что ошибка не возникла
	assert.NoError(t, err)
	// Проверяем, что пользователь был обновлен
	assert.True(t, updated)

	// Получаем обновленного пользователя из репозитория
	updatedUser, _ := repo.GetByID(userID)

	// Проверяем, что данные пользователя обновлены
	assert.True(t, newUserData.Name == updatedUser.Name &&
		newUserData.Surname == updatedUser.Surname &&
		newUserData.Login == updatedUser.Login &&
		newUserData.AvatarID == updatedUser.AvatarID)
}

func TestDeleteUser(t *testing.T) {
	// Создаем репозиторий в памяти
	repo := NewEmptyInMemoryUserRepository()

	// Создаем нового пользователя
	newUser := &userCore.User{
		Name:     "John",
		Surname:  "Doe",
		Login:    "john_doe",
		Password: "secret123",
	}

	// Добавляем пользователя в репозиторий
	userID, _ := repo.Add(newUser)

	// Удаляем пользователя из репозитория
	deleted, err := repo.Delete(userID)

	// Проверяем, что ошибка не возникла
	assert.NoError(t, err)
	// Проверяем, что пользователь был удален
	assert.True(t, deleted)

	// Пытаемся получить удаленного пользователя
	deletedUser, _ := repo.GetByID(userID)

	// Проверяем, что удаленного пользователя нет
	assert.Nil(t, deletedUser)
}

func TestGetUserByLogin(t *testing.T) {
	// Создаем репозиторий в памяти
	repo := NewInMemoryUserRepository()

	// Создаем нового пользователя
	newUser := &userCore.User{
		Name:     "John",
		Surname:  "Doe",
		Login:    "john_doe",
		Password: "1234",
	}

	// Добавляем пользователя в репозиторий
	_, _ = repo.Add(newUser)

	// Тестируем успешный поиск пользователя
	foundUser, err := repo.GetUserByLogin("john_doe", "1234")
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, newUser.Login, foundUser.Login)

	// Тестируем поиск пользователя с неверным паролем
	_, err = repo.GetUserByLogin(newUser.Login, "wrong_password")
	assert.Error(t, err)

	// Тестируем поиск несуществующего пользователя
	_, err = repo.GetUserByLogin("nonexistent_user", "some_password")
	assert.Error(t, err)
}
