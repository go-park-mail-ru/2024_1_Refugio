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
	emails, err := uc.repo.GetAllIncoming(login, offset, limit, ctx)
	if err != nil {
		return nil, err
	}

	for _, email := range emails {
		if validators.IsValidEmailFormat(email.SenderEmail) {
			email.PhotoID, err = uc.repo.GetAvatarFileIDByLogin(email.SenderEmail, ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	return emails, nil
}

// GetAllEmailsSent returns all emails sent.
func (uc *EmailUseCase) GetAllEmailsSent(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	emails, err := uc.repo.GetAllSent(login, offset, limit, ctx)
	if err != nil {
		return nil, err
	}

	for _, email := range emails {
		if validators.IsValidEmailFormat(email.RecipientEmail) {
			email.PhotoID, err = uc.repo.GetAvatarFileIDByLogin(email.RecipientEmail, ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	return emails, nil
}

// GetAllDraftEmails returns all draft emails.
func (uc *EmailUseCase) GetAllDraftEmails(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	emails, err := uc.repo.GetAllDraft(login, offset, limit, ctx)
	if err != nil {
		return nil, err
	}

	for _, email := range emails {
		if validators.IsValidEmailFormat(email.RecipientEmail) {
			email.PhotoID, err = uc.repo.GetAvatarFileIDByLogin(email.RecipientEmail, ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	return emails, nil
}

// GetAllSpamEmails returns all draft emails.
func (uc *EmailUseCase) GetAllSpamEmails(login string, offset, limit int64, ctx context.Context) ([]*domain.Email, error) {
	emails, err := uc.repo.GetAllSpam(login, offset, limit, ctx)
	if err != nil {
		return nil, err
	}

	for _, email := range emails {
		if validators.IsValidEmailFormat(email.SenderEmail) {
			email.PhotoID, err = uc.repo.GetAvatarFileIDByLogin(email.SenderEmail, ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	return emails, nil
}

// GetEmailByID returns the email by its ID.
func (uc *EmailUseCase) GetEmailByID(id uint64, login string, ctx context.Context) (*domain.Email, error) {
	email, err := uc.repo.GetByID(id, login, ctx)
	if err != nil {
		return nil, err
	}

	if validators.IsValidEmailFormat(email.SenderEmail) && validators.IsValidEmailFormat(email.RecipientEmail) {
		if email.SenderEmail == login {
			email.PhotoID, err = uc.repo.GetAvatarFileIDByLogin(email.RecipientEmail, ctx)
			if err != nil {
				return nil, err
			}
		} else if email.RecipientEmail == login {
			email.PhotoID, err = uc.repo.GetAvatarFileIDByLogin(email.SenderEmail, ctx)
			if err != nil {
				return nil, err
			}
		}
	}

	return email, nil
}

// CreateEmail creates a new email.
func (uc *EmailUseCase) CreateEmail(newEmail *domain.Email, ctx context.Context) (uint64, *domain.Email, error) {
	return uc.repo.Add(newEmail, ctx)
}

// CreateProfileEmail creates a new profile_email
func (uc *EmailUseCase) CreateProfileEmail(emailId uint64, sender, recipient string, ctx context.Context) error {
	if sender == recipient {
		return uc.repo.AddProfileEmailMyself(emailId, sender, ctx)
	} else if validators.IsValidEmailFormat(sender) && recipient == "" {
		return uc.repo.AddProfileEmailMyself(emailId, sender, ctx)
	} else if validators.IsValidEmailFormat(sender) && !validators.IsValidEmailFormat(recipient) {
		return uc.repo.AddProfileEmailMyself(emailId, sender, ctx)
	} else if !validators.IsValidEmailFormat(sender) && validators.IsValidEmailFormat(recipient) {
		return uc.repo.AddProfileEmailMyself(emailId, recipient, ctx)
	}

	return uc.repo.AddProfileEmail(emailId, sender, recipient, ctx)
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
func (uc *EmailUseCase) AddAttachment(fileID, fileType, fileName, fileSize string, emailID uint64, ctx context.Context) (uint64, error) {
	if validators.IsEmpty(fileID) || validators.IsEmpty(fileType) || validators.IsEmpty(fileName) || validators.IsEmpty(fileSize) {
		return 0, fmt.Errorf("file id or file type or file name or file size is empty")
	}

	if emailID <= 0 {
		return 0, fmt.Errorf("invalid email id")
	}

	fileId, err := uc.repo.AddFile(fileID, fileType, fileName, fileSize, ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to add file")
	}

	err = uc.repo.AddAttachment(emailID, fileId, ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to add attachment")
	}

	return fileId, nil
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
func (uc *EmailUseCase) UpdateFileByID(fileID uint64, newFileID, newFileType, newFileName, newFileSize string, ctx context.Context) (bool, error) {
	if validators.IsEmpty(newFileID) || validators.IsEmpty(newFileType) || validators.IsEmpty(newFileName) || validators.IsEmpty(newFileSize) {
		return false, fmt.Errorf("file id or file type or file name or file size is empty")
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
	if newFileName != "" && file.FileName != newFileName {
		file.FileName = newFileName
	}
	if newFileSize != "" && file.FileSize != newFileSize {
		file.FileSize = newFileSize
	}

	err = uc.repo.UpdateFileByID(fileID, file.FileId, file.FileType, file.FileName, file.FileSize, ctx)
	if err != nil {
		return false, fmt.Errorf("failed to update file")
	}

	return true, nil
}

// AddFile add an file to database.
func (uc *EmailUseCase) AddFile(fileID, fileType, fileName, fileSize string, ctx context.Context) (uint64, error) {
	if validators.IsEmpty(fileID) || validators.IsEmpty(fileType) || validators.IsEmpty(fileName) || validators.IsEmpty(fileSize) {
		return 0, fmt.Errorf("file id or file type or file name or file size is empty")
	}

	fileId, err := uc.repo.AddFile(fileID, fileType, fileName, fileSize, ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to add file")
	}

	return fileId, nil
}

// AddFileToEmail add a file to an email.
func (uc *EmailUseCase) AddFileToEmail(emailID uint64, fileID uint64, ctx context.Context) error {
	if emailID <= 0 || fileID <= 0 {
		return fmt.Errorf("invalid file id")
	}

	err := uc.repo.AddAttachment(emailID, fileID, ctx)
	if err != nil {
		return fmt.Errorf("failed to add attachment")
	}

	return nil
}
