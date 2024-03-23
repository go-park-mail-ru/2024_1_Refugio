package email

import (
	emailCore "mail/pkg/domain/models"
	"mail/pkg/domain/repository"
)

// EmailUseCase represents the use case for working with emails.
type EmailUseCase struct {
	repo repository.EmailRepository
}

// NewEmailUseCase creates a new instance of EmailUseCase.
func NewEmailUseCase(repo repository.EmailRepository) *EmailUseCase {
	return &EmailUseCase{
		repo: repo,
	}
}

// GetAllEmails returns all emails.
func (uc *EmailUseCase) GetAllEmails() ([]*emailCore.Email, error) {
	return uc.repo.GetAll()
}

// GetEmailByID returns the email by its ID.
func (uc *EmailUseCase) GetEmailByID(id uint64) (*emailCore.Email, error) {
	return uc.repo.GetByID(id)
}

// CreateEmail creates a new email.
func (uc *EmailUseCase) CreateEmail(newEmail *emailCore.Email) (*emailCore.Email, error) {
	return uc.repo.Add(newEmail)
}

// UpdateEmail updates the information of an email.
func (uc *EmailUseCase) UpdateEmail(updatedEmail *emailCore.Email) (bool, error) {
	return uc.repo.Update(updatedEmail)
}

// DeleteEmail deletes the email.
func (uc *EmailUseCase) DeleteEmail(id uint64) (bool, error) {
	return uc.repo.Delete(id)
}
