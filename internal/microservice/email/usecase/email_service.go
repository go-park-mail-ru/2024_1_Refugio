package usecase

import (
	"context"
	"fmt"

	"mail/internal/pkg/utils/validators"

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

// GetAllEmailsIncoming returns all emails incoming.
func (uc *EmailUseCase) GetAllEmailsIncoming(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	return uc.repo.GetAllIncoming(login, offset, limit, ctx)
}

// GetAllEmailsSent returns all emails sent.
func (uc *EmailUseCase) GetAllEmailsSent(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	return uc.repo.GetAllSent(login, offset, limit, ctx)
}

// GetAllDraftEmails returns all draft emails.
func (uc *EmailUseCase) GetAllDraftEmails(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	return uc.repo.GetAllDraft(login, offset, limit, ctx)
}

// GetAllSpamEmails returns all draft emails.
func (uc *EmailUseCase) GetAllSpamEmails(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	return uc.repo.GetAllSpam(login, offset, limit, ctx)
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
	if sender == recipient {
		return uc.repo.AddProfileEmailMyself(email_id, sender, ctx)
	} else if validators.IsValidEmailFormat(sender) == true && recipient == "" {
		return uc.repo.AddProfileEmailMyself(email_id, sender, ctx)
	} else if validators.IsValidEmailFormat(sender) == true && validators.IsValidEmailFormat(recipient) == false {
		return uc.repo.AddProfileEmailMyself(email_id, sender, ctx)
	} else if validators.IsValidEmailFormat(sender) == false && validators.IsValidEmailFormat(recipient) == true {
		return uc.repo.AddProfileEmailMyself(email_id, recipient, ctx)
	}

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

// AddAttachment adds an attachment to the specified email.
func (uc *EmailUseCase) AddAttachment(fileID string, fileType string, emailID uint64, ctx context.Context) (uint64, error) {
	if validators.IsEmpty(fileID) || validators.IsEmpty(fileType) {
		return 0, fmt.Errorf("file id or file type is empty")
	}

	if emailID <= 0 {
		return 0, fmt.Errorf("invalid email id")
	}

	file_id, err := uc.repo.AddFile(fileID, fileType, ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to add file")
	}

	err = uc.repo.AddAttachment(emailID, file_id, ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to add attachment")
	}

	return file_id, nil
}

// GetFileByID returns the file with the specified ID.
func (uc *EmailUseCase) GetFileByID(fileID uint64, ctx context.Context) (*domain.File, error) {
	if fileID <= 0 {
		return nil, fmt.Errorf("invalid file id")
	}

	file, err := uc.repo.GetFileByID(fileID, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get file")
	}

	return file, nil
}

// GetFilesByEmailID returns all files attached to the specified email.
func (uc *EmailUseCase) GetFilesByEmailID(emailID uint64, ctx context.Context) ([]*domain.File, error) {
	if emailID <= 0 {
		return nil, fmt.Errorf("invalid email id")
	}

	files, err := uc.repo.GetFilesByEmailID(emailID, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get files")
	}

	return files, nil
}

// DeleteFileByID deletes the file with the specified ID.
func (uc *EmailUseCase) DeleteFileByID(fileID uint64, ctx context.Context) (bool, error) {
	if fileID <= 0 {
		return false, fmt.Errorf("invalid file id")
	}

	err := uc.repo.DeleteFileByID(fileID, ctx)
	if err != nil {
		return false, fmt.Errorf("failed to delete file")
	}

	return true, nil
}

// UpdateFileByID updates the information of the specified file.
func (uc *EmailUseCase) UpdateFileByID(fileID uint64, newFileID string, newFileType string, ctx context.Context) (bool, error) {
	if validators.IsEmpty(newFileID) || validators.IsEmpty(newFileType) {
		return false, fmt.Errorf("file id or file type is empty")
	}

	if fileID <= 0 {
		return false, fmt.Errorf("invalid file id")
	}

	file, err := uc.repo.GetFileByID(fileID, ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get file")
	}

	if newFileID != "" && file.FileId != newFileID {
		file.FileId = newFileID
	}
	if newFileType != "" && file.FileType != newFileType {
		file.FileType = newFileType
	}

	err = uc.repo.UpdateFileByID(fileID, file.FileId, file.FileType, ctx)
	if err != nil {
		return false, fmt.Errorf("failed to update file")
	}

	return true, nil
}

// AddFile add an file to database.
func (uc *EmailUseCase) AddFile(fileID string, fileType string, ctx context.Context) (uint64, error) {
	if validators.IsEmpty(fileID) || validators.IsEmpty(fileType) {
		return 0, fmt.Errorf("file id or file type is empty")
	}

	file_id, err := uc.repo.AddFile(fileID, fileType, ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to add file")
	}

	return file_id, nil
}
