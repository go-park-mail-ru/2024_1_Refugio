package email

import "sync"

// EmailMemoryRepository represents the implementation of EmailRepository using an in-memory storage.
type EmailMemoryRepository struct {
	mu     sync.RWMutex
	emails map[uint64]*Email
}

// NewEmailMemoryRepository creates a new instance of EmailMemoryRepository.
func NewEmailMemoryRepository() *EmailMemoryRepository {
	return &EmailMemoryRepository{
		emails: make(map[uint64]*Email),
	}
}

func CreateFakeEmails() *EmailMemoryRepository {
	repo := NewEmailMemoryRepository()
	for i := 0; i < len(FakeEmails); i++ {
		repo.emails[uint64(i+1)] = FakeEmails[i]
	}
	return repo
}

// GetAll returns all emails from the storage.
func (repo *EmailMemoryRepository) GetAll() ([]*Email, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	emails := make([]*Email, 0, len(repo.emails))
	for i := 0; i < len(repo.emails); i++ {
		emails = append(emails, repo.emails[uint64(i+1)])
	}

	return emails, nil
}

// GetByID returns an email based on its unique identifier.
func (repo *EmailMemoryRepository) GetByID(id uint64) (*Email, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	email, found := repo.emails[id]
	if !found {
		return nil, nil
	}

	return email, nil
}

// Add adds a new email to the storage and returns the assigned unique identifier.
func (repo *EmailMemoryRepository) Add(email *Email) (uint64, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	id := uint64(len(repo.emails) + 1)
	email.ID = id
	repo.emails[id] = email

	return id, nil
}

// Update updates the data of an email in the storage based on the provided new email.
func (repo *EmailMemoryRepository) Update(newEmail *Email) (bool, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	existingEmail, found := repo.emails[newEmail.ID]
	if !found {
		return false, nil
	}

	existingEmail.Topic = newEmail.Topic
	existingEmail.Text = newEmail.Text
	existingEmail.PhotoID = newEmail.PhotoID
	existingEmail.ReadStatus = newEmail.ReadStatus
	existingEmail.Mark = newEmail.Mark
	existingEmail.Deleted = newEmail.Deleted
	existingEmail.DateOfDispatch = newEmail.DateOfDispatch
	existingEmail.ReplyToEmailID = newEmail.ReplyToEmailID
	existingEmail.DraftStatus = newEmail.DraftStatus

	return true, nil
}

// Delete removes an email from the storage based on its unique identifier.
func (repo *EmailMemoryRepository) Delete(id uint64) (bool, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	_, found := repo.emails[id]
	if !found {
		return false, nil
	}

	delete(repo.emails, id)
	return true, nil
}
