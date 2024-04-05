package email

import (
	domain "mail/pkg/domain/models"
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
func (uc *EmailUseCase) GetAllEmails(login string, offset, limit int) ([]*domain.Email, error) {
	return uc.repo.GetAll(login, offset, limit)
}

// GetEmailByID returns the email by its ID.
func (uc *EmailUseCase) GetEmailByID(id uint64, login string) (*domain.Email, error) {
	return uc.repo.GetByID(id, login)
}

// CreateEmail creates a new email.
func (uc *EmailUseCase) CreateEmail(newEmail *domain.Email) (int64, *domain.Email, error) {
	return uc.repo.Add(newEmail)
}

func (uc *EmailUseCase) CreateProfileEmail(email_id int64, sender, recipient string) error {
	return uc.repo.AddProfileEmail(email_id, sender, recipient)
}

func (uc *EmailUseCase) CheckRecipientEmail(recipient string) error {
	if er := uc.repo.FindEmail(recipient); er != nil {
		return er
	}
	return nil
}

// UpdateEmail updates the information of an email.
func (uc *EmailUseCase) UpdateEmail(updatedEmail *domain.Email) (bool, error) {
	return uc.repo.Update(updatedEmail)
}

// DeleteEmail deletes the email.
func (uc *EmailUseCase) DeleteEmail(id uint64, login string) (bool, error) {
	return uc.repo.Delete(id, login)
}
