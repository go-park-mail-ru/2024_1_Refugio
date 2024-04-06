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

// GetAllEmails returns all emails incoming.
func (uc *EmailUseCase) GetAllEmailsIncoming(login, requestID string, offset, limit int) ([]*domain.Email, error) {
	return uc.repo.GetAllIncoming(login, requestID, offset, limit)
}

// GetAllEmails returns all emails sent.
func (uc *EmailUseCase) GetAllEmailsSent(login, requestID string, offset, limit int) ([]*domain.Email, error) {
	return uc.repo.GetAllSent(login, requestID, offset, limit)
}

// GetEmailByID returns the email by its ID.
func (uc *EmailUseCase) GetEmailByID(id uint64, login, requestID string) (*domain.Email, error) {
	return uc.repo.GetByID(id, login, requestID)
}

// CreateEmail creates a new email.
func (uc *EmailUseCase) CreateEmail(newEmail *domain.Email, requestID string) (int64, *domain.Email, error) {
	return uc.repo.Add(newEmail, requestID)
}

func (uc *EmailUseCase) CreateProfileEmail(email_id int64, sender, recipient, requestID string) error {
	return uc.repo.AddProfileEmail(email_id, sender, recipient, requestID)
}

func (uc *EmailUseCase) CheckRecipientEmail(recipient, requestID string) error {
	if er := uc.repo.FindEmail(recipient, requestID); er != nil {
		return er
	}
	return nil
}

// UpdateEmail updates the information of an email.
func (uc *EmailUseCase) UpdateEmail(updatedEmail *domain.Email, requestID string) (bool, error) {
	return uc.repo.Update(updatedEmail, requestID)
}

// DeleteEmail deletes the email.
func (uc *EmailUseCase) DeleteEmail(id uint64, login, requestID string) (bool, error) {
	return uc.repo.Delete(id, login, requestID)
}
