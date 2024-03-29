package email

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	domain "mail/pkg/domain/models"
	"testing"
	"time"
)

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := EmailRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	t.Run("NoOffsetAndLimit", func(t *testing.T) {
		expectedEmails := []*domain.Email{
			{ID: 1, Topic: "Topic 1", Text: "Text 1"},
			{ID: 2, Topic: "Topic 2", Text: "Text 2"},
			{ID: 3, Topic: "Topic 3", Text: "Text 3"},
		}

		rows := sqlmock.NewRows([]string{"id", "topic", "text"}).
			AddRow(1, "Topic 1", "Text 1").
			AddRow(2, "Topic 2", "Text 2").
			AddRow(3, "Topic 3", "Text 3")

		mock.ExpectQuery(`SELECT \* FROM emails`).WillReturnRows(rows)

		emails, err := repo.GetAll(-1, -1)
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

		mock.ExpectQuery(`SELECT \* FROM emails`).WillReturnRows(rows)

		emails, err := repo.GetAll(-1, -1)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmails, emails)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM emails`).WillReturnError(sql.ErrNoRows)

		emails, err := repo.GetAll(-1, -1)
		assert.Error(t, err)
		assert.Nil(t, emails)
	})
}

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := EmailRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	t.Run("EmailExists", func(t *testing.T) {
		expectedEmail := &domain.Email{ID: 1, Topic: "Topic 1", Text: "Text 1"}
		rows := sqlmock.NewRows([]string{"id", "topic", "text"}).AddRow(expectedEmail.ID, expectedEmail.Topic, expectedEmail.Text)
		mock.ExpectQuery(`SELECT \* FROM emails WHERE id = \$1`).WithArgs(1).WillReturnRows(rows)

		email, err := repo.GetByID(expectedEmail.ID)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmail, email)
	})

	t.Run("EmailNotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM emails WHERE id = \$1`).WithArgs(1).WillReturnError(sql.ErrNoRows)

		email, err := repo.GetByID(1)
		assert.Nil(t, email)
		assert.Error(t, err)
		expectedErrorMessage := fmt.Sprintf("email with id %d not found", 1)
		assert.EqualError(t, err, expectedErrorMessage)
	})

}

func TestAddEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := EmailRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	t.Run("EmailAddedSuccessfully", func(t *testing.T) {
		email := &domain.Email{
			Topic:          "Test Topic",
			Text:           "Test Text",
			PhotoID:        "",
			ReadStatus:     false,
			Flag:           false,
			Deleted:        false,
			DateOfDispatch: time.Now(),
			DraftStatus:    false,
			SenderID:       1,
			RecipientID:    2,
		}
		//mock.ExpectQuery(`SELECT MAX(id) FROM emails`)
		mock.ExpectExec(`INSERT INTO emails`).
			WithArgs(email.Topic, email.Text, sqlmock.AnyArg(), email.PhotoID, email.SenderID, email.RecipientID, email.ReadStatus, email.Deleted, email.DraftStatus, email.Flag).
			WillReturnResult(sqlmock.NewResult(1, 1))

		emailRes, err := repo.Add(email)
		assert.NoError(t, err)
		assert.Equal(t, email, emailRes)
	})

	t.Run("EmailAddFailed", func(t *testing.T) {
		email := &domain.Email{
			Topic:          "Test Topic",
			Text:           "Test Text",
			PhotoID:        "",
			ReadStatus:     false,
			Flag:           false,
			Deleted:        false,
			DateOfDispatch: time.Now(),
			DraftStatus:    false,
			SenderID:       1,
			RecipientID:    2,
		}

		mock.ExpectQuery(`INSERT INTO emails`).
			WithArgs(email.Topic, email.Text, sqlmock.AnyArg(), email.PhotoID, email.SenderID, email.RecipientID, email.ReadStatus, email.Deleted, email.DraftStatus, email.Flag).
			WillReturnError(fmt.Errorf("failed to insert email"))

		emailRes, err := repo.Add(email)
		assert.Error(t, err)
		assert.Equal(t, &domain.Email{}, emailRes)
	})
}

func TestUpdateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := EmailRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	t.Run("EmailUpdatedSuccessfully", func(t *testing.T) {
		newEmail := &domain.Email{
			ID:             1,
			Topic:          "Test Topic",
			Text:           "Test Text",
			PhotoID:        "",
			ReadStatus:     false,
			Flag:           false,
			Deleted:        false,
			DateOfDispatch: time.Now(),
			DraftStatus:    false,
			SenderID:       1,
			RecipientID:    2,
		}

		mock.ExpectExec(`UPDATE emails`).
			WithArgs(newEmail.Topic, newEmail.Text, newEmail.PhotoID, newEmail.ReadStatus, newEmail.Deleted, newEmail.DraftStatus, sqlmock.AnyArg(), newEmail.Flag, newEmail.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		updated, err := repo.Update(newEmail)

		assert.NoError(t, err)
		assert.True(t, updated)
	})

	t.Run("EmailUpdateFailedNoRowsAffected", func(t *testing.T) {
		newEmail := &domain.Email{
			ID:             2,
			Topic:          "Test Topic",
			Text:           "Test Text",
			PhotoID:        "",
			ReadStatus:     false,
			Flag:           false,
			Deleted:        false,
			DateOfDispatch: time.Now(),
			DraftStatus:    false,
			SenderID:       1,
			RecipientID:    2,
		}

		mock.ExpectExec(`UPDATE emails`).
			WithArgs(newEmail.Topic, newEmail.Text, newEmail.PhotoID, newEmail.ReadStatus, newEmail.Deleted, newEmail.DraftStatus, sqlmock.AnyArg(), newEmail.Flag, newEmail.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		updated, err := repo.Update(newEmail)

		assert.Error(t, err)
		assert.False(t, updated)
	})

	t.Run("EmailUpdateFailedDBError", func(t *testing.T) {
		newEmail := &domain.Email{
			ID:             2,
			Topic:          "Test Topic",
			Text:           "Test Text",
			PhotoID:        "",
			ReadStatus:     false,
			Flag:           false,
			Deleted:        false,
			DateOfDispatch: time.Now(),
			DraftStatus:    false,
			SenderID:       1,
			RecipientID:    2,
		}

		mock.ExpectExec(`UPDATE emails`).
			WithArgs(newEmail.Topic, newEmail.Text, newEmail.PhotoID, newEmail.ReadStatus, newEmail.Deleted, newEmail.DraftStatus, sqlmock.AnyArg(), newEmail.Flag, newEmail.ID).
			WillReturnError(fmt.Errorf("database error"))

		updated, err := repo.Update(newEmail)

		assert.Error(t, err)
		assert.False(t, updated)
	})
}

func TestDeleteEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := EmailRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	t.Run("EmailDeletedSuccessfully", func(t *testing.T) {
		emailID := uint64(1)

		mock.ExpectExec(`DELETE FROM emails`).WithArgs(emailID).WillReturnResult(sqlmock.NewResult(0, 1))

		deleted, err := repo.Delete(emailID)

		assert.NoError(t, err)
		assert.True(t, deleted)
	})

	t.Run("EmailDeleteFailedNoRowsAffected", func(t *testing.T) {
		emailID := uint64(2)

		mock.ExpectExec(`DELETE FROM email`).WithArgs(emailID).WillReturnResult(sqlmock.NewResult(0, 0))

		deleted, err := repo.Delete(emailID)

		assert.Error(t, err)
		assert.False(t, deleted)
	})

	t.Run("EmailDeleteFailedDBError", func(t *testing.T) {
		emailID := uint64(3)

		mock.ExpectExec(`DELETE FROM emails`).WithArgs(emailID).WillReturnError(fmt.Errorf("database error"))

		deleted, err := repo.Delete(emailID)

		assert.Error(t, err)
		assert.False(t, deleted)
	})
}
