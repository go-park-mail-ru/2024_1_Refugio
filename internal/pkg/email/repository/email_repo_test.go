package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	domain "mail/internal/models/domain_models"
	"mail/internal/pkg/logger"
	"os"
	"testing"
	"time"
)

func GetCTX() context.Context {
	requestID := "testID"

	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
	}
	defer f.Close()

	c := context.WithValue(context.Background(), "logger", logger.InitializationBdLog(f))
	ctx := context.WithValue(c, "requestID", requestID)
	return ctx
}

func TestNewEmailRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := EmailRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	EmailRepo := NewEmailRepository(repo.DB)

	assert.Equal(t, repo, *EmailRepo)
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

	ctx := GetCTX()

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
			SenderEmail:    "ivan@mailhub.su",
			RecipientEmail: "sergey@mailhub.su",
			ReplyToEmailID: 0,
		}

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(`
			INSERT INTO email \(topic, text, date_of_dispatch, photoid, sender_email, recipient_email, read_status, deleted_status, draft_status, reply_to_email_id, flag\) 
			VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10, \$11\) 
			RETURNING id
		`).
			WithArgs(email.Topic, email.Text, sqlmock.AnyArg(), email.PhotoID, email.SenderEmail, email.RecipientEmail, email.ReadStatus, email.Deleted, email.DraftStatus, nil, email.Flag).
			WillReturnRows(rows)

		id, emailRes, err := repo.Add(email, ctx)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
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
			SenderEmail:    "ivan@mailhub.su",
			RecipientEmail: "sergey@mailhub.su",
		}

		mock.ExpectQuery(`
			INSERT INTO email \(topic, text, date_of_dispatch, photoid, sender_email, recipient_email, read_status, deleted_status, draft_status, flag\) 
			VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10\) 
			RETURNING id
		`).WithArgs(email.Topic, email.Text, sqlmock.AnyArg(), email.PhotoID, email.SenderEmail, email.RecipientEmail, email.ReadStatus, email.Deleted, email.DraftStatus, email.Flag).
			WillReturnError(fmt.Errorf("failed to insert email"))

		id, emailRes, err := repo.Add(email, ctx)
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
		assert.Equal(t, &domain.Email{}, emailRes)
	})
}

func TestAddProfileEmail(t *testing.T) {
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

	ctx := GetCTX()

	t.Run("ddProfileEmailSuccessfully", func(t *testing.T) {
		emailID := int64(1)
		sender := "sender_test@mailhub.su"
		recipient := "recipient_test@mailhub.su"

		mock.ExpectExec(`
			INSERT INTO profile_email \(profile_id, email_id\)
			VALUES \(\(SELECT id FROM profile WHERE login=\$1\), \$3\), \(\(SELECT id FROM profile WHERE login=\$2\), \$3\)
		`).
			WithArgs(sender, recipient, emailID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AddProfileEmail(emailID, sender, recipient, ctx)
		assert.NoError(t, err)
	})

	t.Run("ddProfileEmailFailed", func(t *testing.T) {
		emailID := int64(1)
		sender := "sender_test@mailhub.su"
		recipient := "recipient_test@mailhub.su"

		mock.ExpectExec(`
			INSERT INTO profile_email \(profile_id, email_id\)
			VALUES \(\(SELECT id FROM profile WHERE login=\$1\), \$3\), \(\(SELECT id FROM profile WHERE login=\$2\), \$3\)
		`).
			WithArgs(recipient, emailID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.AddProfileEmail(emailID, sender, recipient, ctx)
		assert.Error(t, err)
	})
}

func TestFindEmail(t *testing.T) {
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

	login := "test@nailhub.su"
	ctx := GetCTX()

	t.Run("FindEmailSuccessfully", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"login"}).AddRow(login)
		mock.ExpectQuery(`SELECT \* FROM profile WHERE login = \$1`).
			WithArgs(login).
			WillReturnRows(rows)

		err := repo.FindEmail(login, ctx)
		assert.NoError(t, err)
	})

	t.Run("FindEmailFailed", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM profile WHERE login = \$1`).
			WithArgs(login).
			WillReturnError(fmt.Errorf("database error"))

		err := repo.FindEmail(login, ctx)
		assert.Error(t, err)
	})
}

func TestGetAllIncoming(t *testing.T) {
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

	login := "test@mailhub.su"
	ctx := GetCTX()

	t.Run("NoOffsetAndLimit", func(t *testing.T) {
		expectedEmails := []*domain.Email{
			{ID: 1, Topic: "Topic 1", Text: "Text 1", RecipientEmail: "test@mailhub.su"},
			{ID: 2, Topic: "Topic 2", Text: "Text 2", RecipientEmail: "test@mailhub.su"},
			{ID: 3, Topic: "Topic 3", Text: "Text 3", RecipientEmail: "test@mailhub.su"},
		}

		rows := sqlmock.NewRows([]string{"id", "topic", "text", "recipient_email"}).
			AddRow(1, "Topic 1", "Text 1", "test@mailhub.su").
			AddRow(2, "Topic 2", "Text 2", "test@mailhub.su").
			AddRow(3, "Topic 3", "Text 3", "test@mailhub.su")

		mock.ExpectQuery(`
			SELECT \* FROM email
			WHERE recipient_email = \$1
			ORDER BY date_of_dispatch ASC
		`).WillReturnRows(rows)

		emails, err := repo.GetAllIncoming(login, -1, -1, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmails, emails)
	})

	t.Run("WithOffsetAndLimit", func(t *testing.T) {
		expectedEmails := []*domain.Email{
			{ID: 2, Topic: "Topic 2", Text: "Text 2", RecipientEmail: "test@mailhub.su"},
			{ID: 3, Topic: "Topic 3", Text: "Text 3", RecipientEmail: "test@mailhub.su"},
		}

		rows := sqlmock.NewRows([]string{"id", "topic", "text", "recipient_email"}).
			AddRow(2, "Topic 2", "Text 2", "test@mailhub.su").
			AddRow(3, "Topic 3", "Text 3", "test@mailhub.su")

		mock.ExpectQuery(`
			SELECT \* FROM email
			WHERE recipient_email = \$1
			ORDER BY date_of_dispatch ASC  
			OFFSET \$2 LIMIT \$3
		`).WillReturnRows(rows)

		emails, err := repo.GetAllIncoming(login, 1, 2, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmails, emails)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery(`
			SELECT \* FROM email
			WHERE recipient_email = \$1
			ORDER BY date_of_dispatch ASC
		`).WillReturnError(sql.ErrNoRows)

		emails, err := repo.GetAllIncoming(login, -1, -1, ctx)
		assert.Error(t, err)
		assert.Nil(t, emails)
	})
}

func TestGetAllSent(t *testing.T) {
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

	login := "test@mailhub.su"
	ctx := GetCTX()

	t.Run("NoOffsetAndLimit", func(t *testing.T) {
		expectedEmails := []*domain.Email{
			{ID: 1, Topic: "Topic 1", Text: "Text 1", SenderEmail: "test@mailhub.su", ReadStatus: true},
			{ID: 2, Topic: "Topic 2", Text: "Text 2", SenderEmail: "test@mailhub.su", ReadStatus: true},
			{ID: 3, Topic: "Topic 3", Text: "Text 3", SenderEmail: "test@mailhub.su", ReadStatus: true},
		}

		rows := sqlmock.NewRows([]string{"id", "topic", "text", "sender_email"}).
			AddRow(1, "Topic 1", "Text 1", "test@mailhub.su").
			AddRow(2, "Topic 2", "Text 2", "test@mailhub.su").
			AddRow(3, "Topic 3", "Text 3", "test@mailhub.su")

		mock.ExpectQuery(`
			SELECT \* FROM email
			WHERE sender_email = \$1
			ORDER BY date_of_dispatch ASC
		`).WillReturnRows(rows)

		emails, err := repo.GetAllSent(login, -1, -1, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmails, emails)
	})

	t.Run("WithOffsetAndLimit", func(t *testing.T) {
		expectedEmails := []*domain.Email{
			{ID: 2, Topic: "Topic 2", Text: "Text 2", SenderEmail: "test@mailhub.su", ReadStatus: true},
			{ID: 3, Topic: "Topic 3", Text: "Text 3", SenderEmail: "test@mailhub.su", ReadStatus: true},
		}

		rows := sqlmock.NewRows([]string{"id", "topic", "text", "sender_email"}).
			AddRow(2, "Topic 2", "Text 2", "test@mailhub.su").
			AddRow(3, "Topic 3", "Text 3", "test@mailhub.su")

		mock.ExpectQuery(`
			SELECT \* FROM email
			WHERE sender_email = \$1
			ORDER BY date_of_dispatch ASC  
			OFFSET \$2 LIMIT \$3
		`).WillReturnRows(rows)

		emails, err := repo.GetAllSent(login, 1, 2, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmails, emails)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery(`
			SELECT \* FROM email
			WHERE sender_email = \$1
			ORDER BY date_of_dispatch ASC
		`).WillReturnError(sql.ErrNoRows)

		emails, err := repo.GetAllSent(login, -1, -1, ctx)
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

	login := "test@mailhub.su"
	ctx := GetCTX()

	t.Run("EmailExists", func(t *testing.T) {
		expectedEmail := &domain.Email{ID: 1, Topic: "Topic 1", Text: "Text 1", SenderEmail: login}
		rows := sqlmock.NewRows([]string{"id", "topic", "text", "sender_email"}).AddRow(expectedEmail.ID, expectedEmail.Topic, expectedEmail.Text, login)
		mock.ExpectQuery(`SELECT \* FROM email WHERE id = \$1 AND \(recipient_email = \$2 OR sender_email = \$2\)`).WithArgs(1, login).WillReturnRows(rows)

		email, err := repo.GetByID(expectedEmail.ID, login, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmail, email)
	})

	t.Run("EmailNotFound", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM email WHERE id = \$1 AND \(recipient_email = \$2 OR sender_email = \$2\)`).WithArgs(1, login).WillReturnError(sql.ErrNoRows)

		email, err := repo.GetByID(1, login, ctx)
		assert.Nil(t, email)
		assert.Error(t, err)
		expectedErrorMessage := fmt.Sprintf("email with id %d not found", 1)
		assert.EqualError(t, err, expectedErrorMessage)
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

	ctx := GetCTX()

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
			SenderEmail:    "ivan@mailhub.su",
			RecipientEmail: "sergey@mailhub.su",
		}

		mock.ExpectExec(`
			UPDATE email
			SET
				topic = \$1, 
				text = \$2, 
				photoid = \$3,
				read_status = \$4, 
				deleted_status = \$5, 
				draft_status = \$6, 
				reply_to_email_id = \$7, 
				flag = \$8
			WHERE
				id = \$9 AND sender_email = \$10
		`).
			WithArgs(newEmail.Topic, newEmail.Text, newEmail.PhotoID, newEmail.ReadStatus, newEmail.Deleted, newEmail.DraftStatus, sqlmock.AnyArg(), newEmail.Flag, newEmail.ID, newEmail.SenderEmail).
			WillReturnResult(sqlmock.NewResult(0, 1))

		updated, err := repo.Update(newEmail, ctx)

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
			SenderEmail:    "ivan@mailhub.su",
			RecipientEmail: "sergey@mailhub.su",
		}

		mock.ExpectExec(`
			UPDATE email
			SET
				topic = \$1, 
				text = \$2, 
				photoid = \$3,
				read_status = \$4, 
				deleted_status = \$5, 
				draft_status = \$6, 
				reply_to_email_id = \$7, 
				flag = \$8
			WHERE
				id = \$9 AND sender_email = \$10
		`).
			WithArgs(newEmail.Topic, newEmail.Text, newEmail.PhotoID, newEmail.ReadStatus, newEmail.Deleted, newEmail.DraftStatus, sqlmock.AnyArg(), newEmail.Flag, newEmail.ID, newEmail.SenderEmail).
			WillReturnResult(sqlmock.NewResult(0, 0))

		updated, err := repo.Update(newEmail, ctx)

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
			SenderEmail:    "ivan@mailhub.su",
			RecipientEmail: "sergey@mailhub.su",
		}

		mock.ExpectExec(`
			UPDATE email
			SET
				topic = \$1, 
				text = \$2, 
				photoid = \$3,
				read_status = \$4, 
				deleted_status = \$5, 
				draft_status = \$6, 
				reply_to_email_id = \$7, 
				flag = \$8
			WHERE
				id = \$9 AND sender_email = \$10
		`).
			WithArgs(newEmail.Topic, newEmail.Text, newEmail.PhotoID, newEmail.ReadStatus, newEmail.Deleted, newEmail.DraftStatus, sqlmock.AnyArg(), newEmail.Flag, newEmail.ID, newEmail.SenderEmail).
			WillReturnError(fmt.Errorf("database error"))

		updated, err := repo.Update(newEmail, ctx)

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

	login := "test@mailhub.su"
	ctx := GetCTX()

	t.Run("EmailDeletedSuccessfully", func(t *testing.T) {
		emailID := uint64(1)

		mock.ExpectExec(`DELETE FROM email WHERE id = \$1 AND \(recipient_email = \$2 OR sender_email = \$2\)`).WithArgs(emailID, login).WillReturnResult(sqlmock.NewResult(0, 1))

		deleted, err := repo.Delete(emailID, login, ctx)

		assert.NoError(t, err)
		assert.True(t, deleted)
	})

	t.Run("EmailDeleteFailedNoRowsAffected", func(t *testing.T) {
		emailID := uint64(2)

		mock.ExpectExec(`DELETE FROM email WHERE id = \$1 AND \(recipient_email = \$2 OR sender_email = \&2\)`).WithArgs(emailID, login).WillReturnResult(sqlmock.NewResult(0, 0))

		deleted, err := repo.Delete(emailID, login, ctx)

		assert.Error(t, err)
		assert.False(t, deleted)
	})

	t.Run("EmailDeleteFailedDBError", func(t *testing.T) {
		emailID := uint64(3)

		mock.ExpectExec(`DELETE FROM email WHERE id = \$1 AND \(recipient_email = \$2 OR sender_email = \&2\)`).WithArgs(emailID, login).WillReturnError(fmt.Errorf("database error"))

		deleted, err := repo.Delete(emailID, login, ctx)

		assert.Error(t, err)
		assert.False(t, deleted)
	})
}
