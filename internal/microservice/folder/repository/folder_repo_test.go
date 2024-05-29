package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"mail/internal/pkg/logger"
	"mail/internal/pkg/utils/constants"

	domain "mail/internal/microservice/models/domain_models"
)

func GetCTX() context.Context {
	ctx := context.WithValue(context.Background(), interface{}(string(constants.LoggerKey)), logger.InitializationBdLog(nil))
	ctx2 := context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), []string{"testID"})

	return ctx2
}

func TestNewFolderRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := FolderRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	FolderRepo := NewFolderRepository(repo.DB)

	assert.Equal(t, repo, *FolderRepo)
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := FolderRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()
	folder := &domain.Folder{
		ProfileId: 1,
		Name:      "Test Folder",
	}

	t.Run("CreateSuccessfully", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(`
			INSERT INTO folder \(profile_id, name\)
			VALUES \(\$1, \$2\)
			RETURNING id
		`).
			WithArgs(folder.ProfileId, folder.Name).
			WillReturnRows(rows)

		id, folderRes, err := repo.Create(folder, ctx)
		assert.NoError(t, err)
		assert.Equal(t, uint32(1), id)
		assert.Equal(t, folder, folderRes)
	})

	t.Run("CreateFail", func(t *testing.T) {
		mock.ExpectQuery(`
			INSERT INTO folder \(profile_id, name\)
			VALUES \(\$1, \$2\)
			RETURNING id
		`).
			WithArgs(folder.ProfileId, folder.Name).
			WillReturnError(fmt.Errorf("failed to insert folder"))

		id, folderRes, err := repo.Create(folder, ctx)
		assert.Error(t, err)
		assert.Equal(t, uint32(0), id)
		assert.Equal(t, &domain.Folder{}, folderRes)
	})
}

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := FolderRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()
	folder := &domain.Folder{
		ProfileId: 1,
		Name:      "Test Folder",
	}

	t.Run("GetAllSuccessfully", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "profile_id"}).
			AddRow(1, "Name 1", 1).
			AddRow(2, "Name 2", 2).
			AddRow(3, "Name 3", 3)

		expectedFolders := []*domain.Folder{
			{ID: 1, Name: "Name 1", ProfileId: 1},
			{ID: 2, Name: "Name 2", ProfileId: 2},
			{ID: 3, Name: "Name 3", ProfileId: 3},
		}

		mock.ExpectQuery(`
			SELECT folder.id, folder.name, folder.profile_id FROM folder
			WHERE profile_id = \$1
		`).WillReturnRows(rows)

		folders, err := repo.GetAll(folder.ProfileId, 0, 0, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedFolders, folders)
	})

	t.Run("WithOffsetAndLimit", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "profile_id"}).
			AddRow(2, "Name 2", 2).
			AddRow(3, "Name 3", 3)

		expectedFolders := []*domain.Folder{
			{ID: 2, Name: "Name 2", ProfileId: 2},
			{ID: 3, Name: "Name 3", ProfileId: 3},
		}

		mock.ExpectQuery(`
			SELECT folder.id, folder.name, folder.profile_id FROM folder
			WHERE profile_id = \$1
			OFFSET \$2 LIMIT \$3
		`).WillReturnRows(rows)

		folders, err := repo.GetAll(folder.ProfileId, 1, 2, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedFolders, folders)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery(`
			SELECT folder.id, folder.name, folder.profile_id FROM folder
			WHERE profile_id = \$1
			OFFSET \$2 LIMIT \$3
		`).WillReturnError(sql.ErrNoRows)

		folders, err := repo.GetAll(folder.ProfileId, 1, 2, ctx)
		assert.Error(t, err)
		assert.Nil(t, folders)
	})
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := FolderRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()
	folder := &domain.Folder{
		ID:        1,
		ProfileId: 1,
		Name:      "Test Folder",
	}

	t.Run("DeleteSuccessfully", func(t *testing.T) {
		mock.ExpectExec(`
			DELETE FROM folder
			WHERE folder.id = \$1 AND profile_id = \$2
		`).WithArgs(folder.ID, folder.ProfileId).WillReturnResult(sqlmock.NewResult(0, 1))

		folderStatus, err := repo.Delete(folder.ID, folder.ProfileId, ctx)
		assert.NoError(t, err)
		assert.True(t, folderStatus)
	})

	t.Run("DeleteFailedNoRowsAffected", func(t *testing.T) {
		mock.ExpectExec(`
			DELETE FROM folder
			WHERE folder.id = \$1 AND profile_id = \$2
		`).WithArgs(folder.ID, folder.ProfileId).WillReturnResult(sqlmock.NewResult(0, 0))

		folderStatus, err := repo.Delete(folder.ID, folder.ProfileId, ctx)
		assert.Error(t, err)
		assert.False(t, folderStatus)
	})

	t.Run("DeleteFailedDBError", func(t *testing.T) {
		mock.ExpectExec(`
			DELETE FROM folder
			WHERE folder.id = \$1 AND profile_id = \$2
		`).WithArgs(folder.ID, folder.ProfileId).WillReturnError(fmt.Errorf("database error"))

		folderStatus, err := repo.Delete(folder.ID, folder.ProfileId, ctx)
		assert.Error(t, err)
		assert.False(t, folderStatus)
	})
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := FolderRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()
	folder := &domain.Folder{
		ID:        1,
		ProfileId: 1,
		Name:      "Test Folder",
	}

	t.Run("UpdateSuccessfully", func(t *testing.T) {
		mock.ExpectExec(`
			UPDATE folder
			SET
				name = \$1
			WHERE
				folder.id = \$2 AND folder.profile_id = \$3
		`).WithArgs(folder.Name, folder.ID, folder.ProfileId).WillReturnResult(sqlmock.NewResult(0, 1))

		folderStatus, err := repo.Update(folder, ctx)
		assert.NoError(t, err)
		assert.True(t, folderStatus)
	})

	t.Run("UpdateFailedNoRowsAffected", func(t *testing.T) {
		mock.ExpectExec(`
			UPDATE folder
			SET
				name = \$1
			WHERE
				folder.id = \$2 AND folder.profile_id = \$3
		`).WithArgs(folder.Name, folder.ID, folder.ProfileId).WillReturnResult(sqlmock.NewResult(0, 0))

		folderStatus, err := repo.Update(folder, ctx)
		assert.Error(t, err)
		assert.False(t, folderStatus)
	})

	t.Run("UpdateFailedDBError", func(t *testing.T) {
		mock.ExpectExec(`
			UPDATE folder
			SET
				name = \$1
			WHERE
				folder.id = \$2 AND folder.profile_id = \$3
		`).WithArgs(folder.Name, folder.ID, folder.ProfileId).WillReturnError(fmt.Errorf("database error"))

		folderStatus, err := repo.Update(folder, ctx)
		assert.Error(t, err)
		assert.False(t, folderStatus)
	})
}

func TestAddEmailFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := FolderRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()
	folder_id := uint32(1)
	email_id := uint32(1)

	t.Run("AddEmailFolderSuccessfully", func(t *testing.T) {
		mock.ExpectExec(`
			INSERT INTO folder_email \(folder_id, email_id\)
			VALUES \(\$1, \$2\)
		`).WithArgs(folder_id, email_id).WillReturnResult(sqlmock.NewResult(0, 1))

		status, err := repo.AddEmailFolder(folder_id, email_id, ctx)
		assert.NoError(t, err)
		assert.True(t, status)
	})

	t.Run("AddEmailFolderFailedDBError", func(t *testing.T) {
		mock.ExpectExec(`
			INSERT INTO folder_email \(folder_id, email_id\)
			VALUES \(\$1, \$2\)
		`).WithArgs(folder_id, email_id).WillReturnError(fmt.Errorf("database error"))

		status, err := repo.AddEmailFolder(folder_id, email_id, ctx)
		assert.Error(t, err)
		assert.False(t, status)
	})
}

func TestDeleteEmailFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := FolderRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()
	folder_id := uint32(1)
	email_id := uint32(1)

	t.Run("DeleteSuccessfully", func(t *testing.T) {
		mock.ExpectExec(`
			DELETE FROM folder_email
			WHERE folder_id = \$1 AND email_id = \$2
		`).WithArgs(folder_id, email_id).WillReturnResult(sqlmock.NewResult(0, 1))

		status, err := repo.DeleteEmailFolder(folder_id, email_id, ctx)
		assert.NoError(t, err)
		assert.True(t, status)
	})

	t.Run("DeleteFailedNoRowsAffected", func(t *testing.T) {
		mock.ExpectExec(`
			DELETE FROM folder_email
			WHERE folder_id = \$1 AND email_id = \$2
		`).WithArgs(folder_id, email_id).WillReturnResult(sqlmock.NewResult(0, 0))

		status, err := repo.DeleteEmailFolder(folder_id, email_id, ctx)
		assert.Error(t, err)
		assert.False(t, status)
	})

	t.Run("DeleteFailedDBError", func(t *testing.T) {
		mock.ExpectExec(`
			DELETE FROM folder_email
			WHERE folder_id = \$1 AND email_id = \$2
		`).WithArgs(folder_id, email_id).WillReturnError(fmt.Errorf("database error"))

		status, err := repo.Delete(folder_id, email_id, ctx)
		assert.Error(t, err)
		assert.False(t, status)
	})
}

func TestCheckFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := FolderRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()
	folder := &domain.Folder{
		ID:        1,
		ProfileId: 1,
		Name:      "Test Folder",
	}

	t.Run("CheckFolderSuccessfully", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "profile_id"}).
			AddRow(1, "Name 1", 1).
			AddRow(2, "Name 2", 2).
			AddRow(3, "Name 3", 3)

		mock.ExpectQuery(`
			SELECT id, profile_id, name 
			FROM folder 
			WHERE folder.id = \$1 AND folder.profile_id = \$2
		`).WillReturnRows(rows)

		status, err := repo.CheckFolder(folder.ID, folder.ProfileId, ctx)
		assert.NoError(t, err)
		assert.True(t, status)
	})

	t.Run("CheckFolderFail", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "profile_id"})

		mock.ExpectQuery(`
				SELECT folder.id, folder.name, folder.profile_id FROM folder
				WHERE profile_id = \$1
				OFFSET \$2 LIMIT \$3
			`).WillReturnRows(rows)

		status, err := repo.CheckFolder(folder.ID, folder.ProfileId, ctx)
		assert.Error(t, err)
		assert.False(t, status)
	})
}

func TestCheckEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := FolderRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()
	email_id := uint32(1)
	profile_id := uint32(1)

	t.Run("CheckEmailSuccessfully", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"profile_id", "email_id"}).
			AddRow(1, 1).
			AddRow(2, 2).
			AddRow(3, 3)

		mock.ExpectQuery(`
			SELECT profile_id, email_id 
			FROM profile_email 
			WHERE email_id = \$1 AND profile_id = \$2
		`).WillReturnRows(rows)

		status, err := repo.CheckEmail(email_id, profile_id, ctx)
		assert.NoError(t, err)
		assert.True(t, status)
	})

	t.Run("CheckEmailFail", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"profile_id", "email_id"})

		mock.ExpectQuery(`
			SELECT profile_id, email_id 
			FROM profile_email 
			WHERE email_id = \$1 AND profile_id = \$2
		`).WillReturnRows(rows)

		status, err := repo.CheckFolder(email_id, profile_id, ctx)
		assert.Error(t, err)
		assert.False(t, status)
	})
}

func TestGetAllEmails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := FolderRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()
	folder := &domain.Folder{
		ID:        1,
		ProfileId: 1,
		Name:      "Test Folder",
	}

	t.Run("GetAllEmailsSuccessfully", func(t *testing.T) {
		expectedEmails := []*domain.Email{
			{ID: 1, Topic: "Topic 1", Text: "Text 1"},
			{ID: 2, Topic: "Topic 2", Text: "Text 2"},
			{ID: 3, Topic: "Topic 3", Text: "Text 3"},
		}

		rows := sqlmock.NewRows([]string{"id", "topic", "text"}).
			AddRow(1, "Topic 1", "Text 1").
			AddRow(2, "Topic 2", "Text 2").
			AddRow(3, "Topic 3", "Text 3")

		mock.ExpectQuery(`
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \$1 
			LEFT JOIN folder_email ON e.id = folder_email.email_id
			WHERE folder_email.folder_id = \$2
			ORDER BY e.date_of_dispatch DESC
		`).WillReturnRows(rows)

		emails, err := repo.GetAllEmails(folder.ID, folder.ProfileId, 0, 0, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmails, emails)
	})

	t.Run("WithOffsetAndLimit", func(t *testing.T) {
		expectedEmails := []*domain.Email{
			{ID: 2, Topic: "Topic 2", Text: "Text 2"},
			{ID: 3, Topic: "Topic 3", Text: "Text 3"},
		}

		rows := sqlmock.NewRows([]string{"id", "topic", "text"}).
			AddRow(2, "Topic 2", "Text 2").
			AddRow(3, "Topic 3", "Text 3")

		mock.ExpectQuery(`
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \$1 
			LEFT JOIN folder_email ON e.id = folder_email.email_id
			WHERE folder_email.folder_id = \$2
			ORDER BY e.date_of_dispatch DESC
			OFFSET \$3 LIMIT \$4
		`).WillReturnRows(rows)

		emails, err := repo.GetAllEmails(folder.ID, folder.ProfileId, 1, 2, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmails, emails)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery(`
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \$1 
			LEFT JOIN folder_email ON e.id = folder_email.email_id
			WHERE folder_email.folder_id = \$2
			ORDER BY e.date_of_dispatch DESC
		`).WillReturnError(sql.ErrNoRows)

		emails, err := repo.GetAllEmails(folder.ID, folder.ProfileId, 0, 0, ctx)
		assert.Error(t, err)
		assert.Equal(t, fmt.Errorf("DB no have emails in folder"), err)
		assert.Nil(t, emails)
	})
}

func TestGetAvatarFileIDByLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := &FolderRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()
	login := "test"

	t.Run("GetAvatarFileIDByLogin_Success", func(t *testing.T) {
		fileId := "123"
		rows := sqlmock.NewRows([]string{"file_id"}).AddRow(fileId)
		mock.ExpectQuery(`
            SELECT f.file_id
            FROM file f
            LEFT JOIN profile p ON p.avatar_id = f.id
            WHERE p.login = \$1
        `).WithArgs(login).WillReturnRows(rows)

		_, err := repo.GetAvatarFileIDByLogin(login, ctx)

		assert.NoError(t, err)
	})

	t.Run("GetAvatarFileIDByLogin_NoAvatar", func(t *testing.T) {
		mock.ExpectQuery(`
            SELECT f.file_id
            FROM file f
            LEFT JOIN profile p ON p.avatar_id = f.id
            WHERE p.login = \$1
        `).WithArgs(login).WillReturnError(sql.ErrNoRows)

		_, err := repo.GetAvatarFileIDByLogin(login, ctx)

		assert.Error(t, err)
	})

	t.Run("GetAvatarFileIDByLogin_DBError", func(t *testing.T) {
		mock.ExpectQuery(`
            SELECT f.file_id
            FROM file f
            LEFT JOIN profile p ON p.avatar_id = f.id
            WHERE p.login = \$1
        `).WithArgs(login).WillReturnError(errors.New("DB error"))

		_, err := repo.GetAvatarFileIDByLogin(login, ctx)

		assert.Error(t, err)
	})
}

func TestGetAllFolderName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := &FolderRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()
	emailID := uint32(123)

	t.Run("GetAllFolderName_Success", func(t *testing.T) {
		expectedFolders := []*domain.Folder{
			{ID: 1, Name: "Inbox"},
			{ID: 2, Name: "Sent"},
			{ID: 3, Name: "Drafts"},
		}

		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Inbox").
			AddRow(2, "Sent").
			AddRow(3, "Drafts")

		mock.ExpectQuery(`
            SELECT f.id, f.name
            FROM folder f
            INNER JOIN folder_email fe ON f.id = fe.folder_id
            WHERE fe.email_id = \$1
        `).WithArgs(emailID).WillReturnRows(rows)

		folders, err := repo.GetAllFolderName(emailID, ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedFolders, folders)
	})

	t.Run("GetAllFolderName_NoFolders", func(t *testing.T) {
		mock.ExpectQuery(`
            SELECT f.id, f.name
            FROM folder f
            INNER JOIN folder_email fe ON f.id = fe.folder_id
            WHERE fe.email_id = \$1
        `).WithArgs(emailID).WillReturnError(sql.ErrNoRows)

		folders, err := repo.GetAllFolderName(emailID, ctx)

		assert.Error(t, err)
		assert.Nil(t, folders)
	})

	t.Run("GetAllFolderName_DBError", func(t *testing.T) {
		mock.ExpectQuery(`
            SELECT f.id, f.name
            FROM folder f
            INNER JOIN folder_email fe ON f.id = fe.folder_id
            WHERE fe.email_id = \$1
        `).WithArgs(emailID).WillReturnError(errors.New("DB error"))

		folders, err := repo.GetAllFolderName(emailID, ctx)

		assert.Error(t, err)
		assert.Nil(t, folders)
	})
}
