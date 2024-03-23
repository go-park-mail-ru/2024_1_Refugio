package email

import (
	"fmt"
	"mail/pkg/repository/converters"
	"sync"

	emailCore "mail/pkg/domain/models"
	"mail/pkg/repository/models"
)

// EmailMemoryRepository represents the implementation of EmailRepository using an in-memory storage.
type EmailMemoryRepository struct {
	mu     sync.RWMutex
	emails map[uint64]*models.Email
}

// NewEmailMemoryRepository creates a new instance of EmailMemoryRepository.
func NewEmailMemoryRepository() *EmailMemoryRepository {
	fakeEmails := FakeEmails

	return &EmailMemoryRepository{
		emails: fakeEmails,
	}
}

// NewEmptyInMemoryEmailRepository creates a new email repository in memory with an empty default email list.
func NewEmptyInMemoryEmailRepository() *EmailMemoryRepository {
	defaultEmails := map[uint64]*models.Email{}

	return &EmailMemoryRepository{
		emails: defaultEmails,
	}
}

func CreateFakeEmails() *EmailMemoryRepository {
	repo := NewEmptyInMemoryEmailRepository()

	for i := 1; i-1 < len(FakeEmails); i++ {
		repo.emails[uint64(i)] = FakeEmails[uint64(i)]
	}

	return repo
}

// GetAll returns all emails from the storage.
func (repository *EmailMemoryRepository) GetAll() ([]*emailCore.Email, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	emails := make([]*emailCore.Email, 0, len(repository.emails))
	for _, email := range repository.emails {
		emails = append(emails, converters.EmailConvertDbInCore(*email))
	}

	return emails, nil
}

// GetByID returns an email based on its unique identifier.
func (repository *EmailMemoryRepository) GetByID(id uint64) (*emailCore.Email, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	email, found := repository.emails[id]
	if !found {
		return nil, fmt.Errorf("Email with id %d not found", id)
	}

	return converters.EmailConvertDbInCore(*email), nil
}

// Add adds a new email to the storage and returns the assigned unique identifier.
func (repository *EmailMemoryRepository) Add(email *emailCore.Email) (*emailCore.Email, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	emailDb := converters.EmailConvertCoreInDb(*email)

	id := uint64(len(repository.emails) + 1)
	emailDb.ID = id
	repository.emails[id] = emailDb

	return converters.EmailConvertDbInCore(*repository.emails[id]), nil
}

// Update updates the data of an email in the storage based on the provided new email.
func (repository *EmailMemoryRepository) Update(newEmail *emailCore.Email) (bool, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	emailDb := converters.EmailConvertCoreInDb(*newEmail)

	existingEmail, found := repository.emails[emailDb.ID]
	if !found {
		return false, fmt.Errorf("Email with id %d not found", emailDb.ID)
	}

	existingEmail.Topic = emailDb.Topic
	existingEmail.Text = emailDb.Text
	existingEmail.PhotoID = emailDb.PhotoID
	existingEmail.ReadStatus = emailDb.ReadStatus
	existingEmail.Mark = emailDb.Mark
	existingEmail.Deleted = emailDb.Deleted
	existingEmail.DateOfDispatch = emailDb.DateOfDispatch
	existingEmail.ReplyToEmailID = emailDb.ReplyToEmailID
	existingEmail.DraftStatus = emailDb.DraftStatus

	return true, nil
}

// Delete removes an email from the storage based on its unique identifier.
func (repository *EmailMemoryRepository) Delete(id uint64) (bool, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	_, found := repository.emails[id]
	if !found {
		return false, fmt.Errorf("Email with id %d not found", id)
	}

	delete(repository.emails, id)

	return true, nil
}
