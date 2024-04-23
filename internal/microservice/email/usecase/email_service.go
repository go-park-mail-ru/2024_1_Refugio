package usecase

import (
	"context"
	repository "mail/internal/microservice/email/interface"
	domain "mail/internal/microservice/models/domain_models"
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
func (uc *EmailUseCase) GetAllEmailsIncoming(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	return uc.repo.GetAllIncoming(login, offset, limit, ctx)
}

// GetAllEmails returns all emails sent.
func (uc *EmailUseCase) GetAllEmailsSent(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	return uc.repo.GetAllSent(login, offset, limit, ctx)
}

// GetEmailByID returns the email by its ID.
func (uc *EmailUseCase) GetEmailByID(id uint64, login string, ctx context.Context) (*domain.Email, error) {
	return uc.repo.GetByID(id, login, ctx)
}

// CreateEmail creates a new email.
func (uc *EmailUseCase) CreateEmail(newEmail *domain.Email, ctx context.Context) (uint64, *domain.Email, error) {
	return uc.repo.Add(newEmail, ctx)
}

// CreateProfileEmail creates a new profile_email
func (uc *EmailUseCase) CreateProfileEmail(email_id uint64, sender, recipient string, ctx context.Context) error {
	return uc.repo.AddProfileEmail(email_id, sender, recipient, ctx)
}

// CheckRecipientEmail checking recipient email
func (uc *EmailUseCase) CheckRecipientEmail(recipient string, ctx context.Context) error {
	if er := uc.repo.FindEmail(recipient, ctx); er != nil {
		return er
	}
	return nil
}

// UpdateEmail updates the information of an email.
func (uc *EmailUseCase) UpdateEmail(updatedEmail *domain.Email, ctx context.Context) (bool, error) {
	return uc.repo.Update(updatedEmail, ctx)
}

// DeleteEmail deletes the email.
func (uc *EmailUseCase) DeleteEmail(id uint64, login string, ctx context.Context) (bool, error) {
	return uc.repo.Delete(id, login, ctx)
}
