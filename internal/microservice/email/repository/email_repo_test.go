package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

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
			SpamStatus:     false,
			SenderEmail:    "ivan@mailhub.su",
			RecipientEmail: "sergey@mailhub.su",
			ReplyToEmailID: 0,
		}

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(`
			INSERT INTO email \(topic, text, date_of_dispatch, sender_email, recipient_email, isRead, isDeleted, isDraft, isSpam, reply_to_email_id, is_important\)
			VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10, \$11\)
			RETURNING id
		`).
			WithArgs(email.Topic, email.Text, sqlmock.AnyArg(), email.SenderEmail, email.RecipientEmail, email.ReadStatus, email.Deleted, email.DraftStatus, email.SpamStatus, nil, email.Flag).
			WillReturnRows(rows)

		mock.ExpectExec(`
			INSERT INTO email_file \(email_id, file_id\)
			SELECT \$1, p.avatar_id
			FROM profile p
			WHERE p.login = \$2
		`).
			WithArgs(1, email.SenderEmail).
			WillReturnResult(sqlmock.NewResult(1, 1))

		id, emailRes, err := repo.Add(email, ctx)
		assert.NoError(t, err)
		assert.Equal(t, uint64(1), id)
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
			INSERT INTO email \(topic, text, date_of_dispatch, sender_email, recipient_email, isRead, isDeleted, isDraft, isSpam, reply_to_email_id, is_important\)
			VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10, \$11\)
			RETURNING id
		`).WithArgs(email.Topic, email.Text, sqlmock.AnyArg(), email.SenderEmail, email.RecipientEmail, email.ReadStatus, email.Deleted, email.DraftStatus, email.SpamStatus, nil, email.Flag).
			WillReturnError(fmt.Errorf("failed to insert email"))

		mock.ExpectExec(`
			INSERT INTO email_file \(email_id, file_id\)
			SELECT \$1, p.avatar_id
			FROM profile p
			WHERE p.login = \$2
		`).
			WithArgs(1, email.SenderEmail).
			WillReturnResult(sqlmock.NewResult(1, 1))

		id, emailRes, err := repo.Add(email, ctx)
		assert.Error(t, err)
		assert.Equal(t, uint64(0), id)
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

	t.Run("ddProfileEmailSuccessfully different login", func(t *testing.T) {
		emailID := int64(1)
		sender := "sender_test@mailhub.su"
		recipient := "recipient_test@mailhub.su"

		mock.ExpectExec(`
			INSERT INTO profile_email \(profile_id, email_id\)
			VALUES \(\(SELECT id FROM profile WHERE login=\$1\), \$3\), \(\(SELECT id FROM profile WHERE login=\$2\), \$3\)
		`).
			WithArgs(sender, recipient, emailID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AddProfileEmail(uint64(emailID), sender, recipient, ctx)
		assert.NoError(t, err)
	})

	t.Run("ddProfileEmailSuccessfully similar login", func(t *testing.T) {
		emailID := int64(1)
		sender := "test@mailhub.su"
		recipient := "test@mailhub.su"

		mock.ExpectExec(`
			INSERT INTO profile_email \(profile_id, email_id\)
			VALUES \(\(SELECT id FROM profile WHERE login=\$1\), \$2\)
		`).
			WithArgs(sender, emailID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AddProfileEmail(uint64(emailID), sender, recipient, ctx)
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

		err := repo.AddProfileEmail(uint64(emailID), sender, recipient, ctx)
		assert.Error(t, err)
	})
}

func TestAddProfileEmailMyself(t *testing.T) {
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

	t.Run("ddProfileEmailSuccessfully different login", func(t *testing.T) {
		emailID := int64(1)
		sender := "sender_test@mailhub.su"

		mock.ExpectExec(`
			INSERT INTO profile_email \(profile_id, email_id\)
			VALUES \(\(SELECT id FROM profile WHERE login=\$1\), \$2\)
		`).
			WithArgs(sender, emailID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AddProfileEmailMyself(uint64(emailID), sender, ctx)
		assert.NoError(t, err)
	})

	t.Run("ddProfileEmailFailed", func(t *testing.T) {
		emailID := int64(1)
		sender := "sender_test@mailhub.su"

		mock.ExpectExec(`
			INSERT INTO profile_email \(profile_id, email_id\)
			VALUES \(\(SELECT id FROM profile WHERE login=\$1\), \$3\)
		`).
			WithArgs(sender, emailID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.AddProfileEmailMyself(uint64(emailID), sender, ctx)
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
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.isSpam, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \(
				SELECT id FROM profile WHERE login = \$1
			\)
			WHERE e.recipient_email = \$1 AND e.isSpam = false AND e.isDraft = false
			ORDER BY e.date_of_dispatch DESC
		`).WillReturnRows(rows)

		emails, err := repo.GetAllIncoming(login, -1, -1, ctx)
		assert.NoError(t, err)
		assert.NotNil(t, emails)
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
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.isSpam, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \(
				SELECT id FROM profile WHERE login = \$1
			\)
			WHERE e.recipient_email = \$1 AND e.isSpam = false AND e.isDraft = false
			ORDER BY e.date_of_dispatch DESC
			OFFSET \$2 LIMIT \$3
		`).WillReturnRows(rows)

		emails, err := repo.GetAllIncoming(login, 1, 2, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmails, emails)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery(`
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.isSpam, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \(
				SELECT id FROM profile WHERE login = \$1
			\)
			WHERE e.recipient_email = \$1 AND e.isSpam = false AND e.isDraft = false
			ORDER BY e.date_of_dispatch DESC
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
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.isSpam, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \(
				SELECT id FROM profile WHERE login = \$1
			\)
			WHERE e.sender_email = \$1 AND e.isSpam = false AND e.isDraft = false
			ORDER BY e.date_of_dispatch DESC
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
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.isSpam, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \(
				SELECT id FROM profile WHERE login = \$1
			\)
			WHERE e.sender_email = \$1 AND e.isSpam = false AND e.isDraft = false
			ORDER BY e.date_of_dispatch DESC
			OFFSET \$2 LIMIT \$3
		`).WillReturnRows(rows)

		emails, err := repo.GetAllSent(login, 1, 2, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmails, emails)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery(`
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \(
				SELECT id FROM profile WHERE login = \$1
			\)
			WHERE e.sender_email = \$1 AND e.isSpam = false AND e.isDraft = false
			ORDER BY e.date_of_dispatch DESC
		`).WillReturnError(sql.ErrNoRows)

		emails, err := repo.GetAllSent(login, -1, -1, ctx)
		assert.Error(t, err)
		assert.Nil(t, emails)
	})
}

func TestGetAllDraft(t *testing.T) {
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
			{ID: 1, Topic: "Topic 1", Text: "Text 1", SenderEmail: "test@mailhub.su", DraftStatus: true},
			{ID: 2, Topic: "Topic 2", Text: "Text 2", SenderEmail: "test@mailhub.su", DraftStatus: true},
			{ID: 3, Topic: "Topic 3", Text: "Text 3", SenderEmail: "test@mailhub.su", DraftStatus: true},
		}

		rows := sqlmock.NewRows([]string{"id", "topic", "text", "sender_email", "isdraft"}).
			AddRow(1, "Topic 1", "Text 1", "test@mailhub.su", true).
			AddRow(2, "Topic 2", "Text 2", "test@mailhub.su", true).
			AddRow(3, "Topic 3", "Text 3", "test@mailhub.su", true)

		mock.ExpectQuery(`
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.isSpam, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \(
				SELECT id FROM profile WHERE login = \$1
			\)
			WHERE e.sender_email = \$1 AND e.isDraft = true
			ORDER BY e.date_of_dispatch DESC
		`).WillReturnRows(rows)

		emails, err := repo.GetAllDraft(login, -1, -1, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmails, emails)
	})

	t.Run("WithOffsetAndLimit", func(t *testing.T) {
		expectedEmails := []*domain.Email{
			{ID: 2, Topic: "Topic 2", Text: "Text 2", SenderEmail: "test@mailhub.su", DraftStatus: true},
			{ID: 3, Topic: "Topic 3", Text: "Text 3", SenderEmail: "test@mailhub.su", DraftStatus: true},
		}

		rows := sqlmock.NewRows([]string{"id", "topic", "text", "sender_email", "isdraft"}).
			AddRow(2, "Topic 2", "Text 2", "test@mailhub.su", true).
			AddRow(3, "Topic 3", "Text 3", "test@mailhub.su", true)

		mock.ExpectQuery(`
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.isSpam, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \(
				SELECT id FROM profile WHERE login = \$1
			\)
			WHERE e.sender_email = \$1 AND e.isDraft = true
			ORDER BY e.date_of_dispatch DESC
			OFFSET \$2 LIMIT \$3
		`).WillReturnRows(rows)

		emails, err := repo.GetAllDraft(login, 1, 2, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmails, emails)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery(`
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \(
				SELECT id FROM profile WHERE login = \$1
			\)
			WHERE e.sender_email = \$1 AND e.isDraft = true
			ORDER BY e.date_of_dispatch DESC
		`).WillReturnError(sql.ErrNoRows)

		emails, err := repo.GetAllDraft(login, -1, -1, ctx)
		assert.Error(t, err)
		assert.Nil(t, emails)
	})
}

func TestGetAllSpam(t *testing.T) {
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
			{ID: 1, Topic: "Topic 1", Text: "Text 1", SenderEmail: "test@mailhub.su", SpamStatus: true},
			{ID: 2, Topic: "Topic 2", Text: "Text 2", SenderEmail: "test@mailhub.su", SpamStatus: true},
			{ID: 3, Topic: "Topic 3", Text: "Text 3", SenderEmail: "test@mailhub.su", SpamStatus: true},
		}

		rows := sqlmock.NewRows([]string{"id", "topic", "text", "sender_email", "isspam"}).
			AddRow(1, "Topic 1", "Text 1", "test@mailhub.su", true).
			AddRow(2, "Topic 2", "Text 2", "test@mailhub.su", true).
			AddRow(3, "Topic 3", "Text 3", "test@mailhub.su", true)

		mock.ExpectQuery(`
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.isSpam, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \(
				SELECT id FROM profile WHERE login = \$1
			\)
			WHERE e.recipient_email = \$1 AND e.isSpam = true
			ORDER BY e.date_of_dispatch DESC
		`).WillReturnRows(rows)

		emails, err := repo.GetAllSpam(login, -1, -1, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmails, emails)
	})

	t.Run("WithOffsetAndLimit", func(t *testing.T) {
		expectedEmails := []*domain.Email{
			{ID: 2, Topic: "Topic 2", Text: "Text 2", SenderEmail: "test@mailhub.su", DraftStatus: true},
			{ID: 3, Topic: "Topic 3", Text: "Text 3", SenderEmail: "test@mailhub.su", DraftStatus: true},
		}

		rows := sqlmock.NewRows([]string{"id", "topic", "text", "sender_email", "isdraft"}).
			AddRow(2, "Topic 2", "Text 2", "test@mailhub.su", true).
			AddRow(3, "Topic 3", "Text 3", "test@mailhub.su", true)

		mock.ExpectQuery(`
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.isSpam, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \(
				SELECT id FROM profile WHERE login = \$1
			\)
			WHERE e.recipient_email = \$1 AND e.isSpam = true
			ORDER BY e.date_of_dispatch DESC
			OFFSET \$2 LIMIT \$3
		`).WillReturnRows(rows)

		emails, err := repo.GetAllSpam(login, 1, 2, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmails, emails)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery(`
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \(
				SELECT id FROM profile WHERE login = \$1
			\)
			WHERE e.sender_email = \$1 AND e.isSpam = true
			ORDER BY e.date_of_dispatch DESC
		`).WillReturnError(sql.ErrNoRows)

		emails, err := repo.GetAllSpam(login, -1, -1, ctx)
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
		mock.ExpectQuery(`
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.isSpam, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \(
				SELECT id FROM profile WHERE login = \$2
			\)
			WHERE e.id = \$1
		`).WithArgs(1, login).WillReturnRows(rows)

		email, err := repo.GetByID(expectedEmail.ID, login, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedEmail, email)
	})

	t.Run("EmailNotFound", func(t *testing.T) {
		mock.ExpectQuery(`
			SELECT DISTINCT e.id, e.topic, e.text, e.date_of_dispatch, e.sender_email, e.recipient_email, e.isRead, e.isDeleted, e.isDraft, e.isSpam, e.reply_to_email_id, e.is_important
			FROM email e
			JOIN profile_email pe ON e.id = pe.email_id
			JOIN profile p ON pe.profile_id = \(
				SELECT id FROM profile WHERE login = \$2
			\)
			WHERE e.id = \$1
		`).WithArgs(1, login).WillReturnError(sql.ErrNoRows)

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
			SpamStatus:     false,
			SenderEmail:    "ivan@mailhub.su",
			RecipientEmail: "sergey@mailhub.su",
		}

		mock.ExpectExec(`
			UPDATE email
			SET
				topic = \$1, 
				text = \$2, 
				isread = \$3, 
				isdeleted = \$4, 
				isdraft = \$5, 
				isspam = \$6,
				reply_to_email_id = \$7, 
				is_important = \$8
			WHERE
				id = \$9 AND sender_email = \$10
		`).
			WithArgs(newEmail.Topic, newEmail.Text, newEmail.ReadStatus, newEmail.Deleted, newEmail.DraftStatus, newEmail.SpamStatus, sqlmock.AnyArg(), newEmail.Flag, newEmail.ID, newEmail.SenderEmail).
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
				isread = \$3, 
				isdeleted = \$4, 
				isdraft = \$5, 
				isspam = \$6,
				reply_to_email_id = \$7, 
				is_important = \$8
			WHERE
				id = \$9 AND sender_email = \$10
		`).
			WithArgs(newEmail.Topic, newEmail.Text, newEmail.ReadStatus, newEmail.Deleted, newEmail.DraftStatus, newEmail.SpamStatus, sqlmock.AnyArg(), newEmail.Flag, newEmail.ID, newEmail.SenderEmail).
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
				isread = \$3, 
				isdeleted = \$4, 
				isdraft = \$5, 
				isspam = \$6,
				reply_to_email_id = \$7, 
				is_important = \$8
			WHERE
				id = \$9 AND sender_email = \$10
		`).
			WithArgs(newEmail.Topic, newEmail.Text, newEmail.ReadStatus, newEmail.Deleted, newEmail.DraftStatus, newEmail.SpamStatus, sqlmock.AnyArg(), newEmail.Flag, newEmail.ID, newEmail.SenderEmail).
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

		mock.ExpectExec(`
			DELETE FROM profile_email
			WHERE profile_id = \(
				SELECT profile_id
				FROM profile_email pe
				JOIN profile p ON pe.profile_id = p.id
				WHERE email_id = \$1 AND p.login = \$2
			\)
			AND email_id = \$1
		`).WithArgs(emailID, login).WillReturnResult(sqlmock.NewResult(0, 1))

		deleted, err := repo.Delete(emailID, login, ctx)

		assert.NoError(t, err)
		assert.True(t, deleted)
	})

	t.Run("EmailDeleteFailedNoRowsAffected", func(t *testing.T) {
		emailID := uint64(2)

		mock.ExpectExec(`
			DELETE FROM profile_email
			WHERE profile_id = \(
				SELECT profile_id
				FROM profile_email pe
				JOIN profile p ON pe.profile_id = p.id
				WHERE email_id = \$1 AND p.login = \$2
			\)
			AND email_id = \$1
		`).WithArgs(emailID, login).WillReturnResult(sqlmock.NewResult(0, 0))

		deleted, err := repo.Delete(emailID, login, ctx)

		assert.Error(t, err)
		assert.False(t, deleted)
	})

	t.Run("EmailDeleteFailedDBError", func(t *testing.T) {
		emailID := uint64(3)

		mock.ExpectExec(`
			DELETE FROM profile_email
			WHERE profile_id = \(
				SELECT profile_id
				FROM profile_email pe
				JOIN profile p ON pe.profile_id = p.id
				WHERE email_id = \$1 AND p.login = \$2
			\)
			AND email_id = \$1
		`).WithArgs(emailID, login).WillReturnError(fmt.Errorf("database error"))

		deleted, err := repo.Delete(emailID, login, ctx)

		assert.Error(t, err)
		assert.False(t, deleted)
	})
}

func TestAddFile(t *testing.T) {
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

	t.Run("FileAddedSuccessfully", func(t *testing.T) {
		fileID := "file123"
		fileType := "text/plain"
		fileName := "PDF"
		fileSize := "10101010"
		expectedID := uint64(1)

		mock.ExpectQuery("INSERT INTO file").
			WithArgs(fileID, fileType, fileName, fileSize).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

		id, err := repo.AddFile(fileID, fileType, fileName, fileSize, ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
	})

	t.Run("FileAddFailed", func(t *testing.T) {
		fileID := "file123"
		fileType := "text/plain"
		fileName := "PDF"
		fileSize := "10101010"

		mock.ExpectQuery("INSERT INTO file").
			WithArgs(fileID, fileType, fileName, fileSize).
			WillReturnError(fmt.Errorf("database error"))

		id, err := repo.AddFile(fileID, fileType, fileName, fileSize, ctx)

		assert.Error(t, err)
		assert.Zero(t, id)
	})
}

func TestAddAttachment(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := EmailRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()

	t.Run("AttachmentAddedSuccessfully", func(t *testing.T) {
		emailID := uint64(1)
		fileID := uint64(2)

		mock.ExpectExec("INSERT INTO email_file").
			WithArgs(emailID, fileID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.AddAttachment(emailID, fileID, ctx)

		assert.NoError(t, err)
	})

	t.Run("AttachmentAddFailed", func(t *testing.T) {
		emailID := uint64(1)
		fileID := uint64(2)

		mock.ExpectExec("INSERT INTO email_file").
			WithArgs(emailID, fileID).
			WillReturnError(fmt.Errorf("database error"))

		err := repo.AddAttachment(emailID, fileID, ctx)

		assert.Error(t, err)
	})
}

func TestGetFileByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := EmailRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()

	t.Run("FileFound", func(t *testing.T) {
		Id := uint64(1)
		fileID := "file123"
		fileType := "text/plain"
		fileName := "PDF"
		fileSize := "10101010"
		expectedFile := &domain.File{ID: Id, FileId: fileID, FileType: fileType, FileName: fileName, FileSize: fileSize}

		mock.ExpectQuery("SELECT file_id, file_type, file_name, file_size FROM file").
			WithArgs(uint64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"file_id", "file_type", "file_name", "file_size"}).AddRow(fileID, fileType, fileName, fileSize))

		file, err := repo.GetFileByID(uint64(1), ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedFile, file)
	})

	t.Run("FileNotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT file_id, file_type, file_name, file_size FROM file").
			WithArgs(uint64(2)).
			WillReturnError(sql.ErrNoRows)

		file, err := repo.GetFileByID(uint64(2), ctx)

		assert.Error(t, err)
		assert.Nil(t, file)
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery("SELECT file_id, file_type FROM file").
			WithArgs(uint64(3)).
			WillReturnError(fmt.Errorf("database error"))

		file, err := repo.GetFileByID(uint64(3), ctx)

		assert.Error(t, err)
		assert.Nil(t, file)
	})
}

func TestGetFilesByEmailID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := EmailRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()

	t.Run("FilesFound", func(t *testing.T) {
		emailID := uint64(1)
		expectedFiles := []*domain.File{
			{ID: 1, FileId: "file123", FileType: "text/plain", FileName: "PDF", FileSize: "10101010"},
			{ID: 2, FileId: "file456", FileType: "image/jpeg", FileName: "PDF", FileSize: "10101010"},
		}

		rows := sqlmock.NewRows([]string{"id", "file_id", "file_type", "file_name", "file_size"}).
			AddRow(1, "file123", "text/plain", "PDF", "10101010").
			AddRow(2, "file456", "image/jpeg", "PDF", "10101010")

		mock.ExpectQuery("SELECT f.id, f.file_id, f.file_type, f.file_name, f.file_size FROM file").
			WithArgs(emailID).
			WillReturnRows(rows)

		files, err := repo.GetFilesByEmailID(emailID, ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedFiles, files)
	})

	t.Run("NoFilesFound", func(t *testing.T) {
		emailID := uint64(2)

		mock.ExpectQuery("SELECT f.id, f.file_id, f.file_type FROM file").
			WithArgs(emailID).
			WillReturnError(sql.ErrNoRows)

		files, err := repo.GetFilesByEmailID(emailID, ctx)

		assert.Error(t, err)
		assert.Empty(t, files)
	})

	t.Run("DBError", func(t *testing.T) {
		emailID := uint64(3)

		mock.ExpectQuery("SELECT f.id, f.file_id, f.file_type FROM file").
			WithArgs(emailID).
			WillReturnError(fmt.Errorf("database error"))

		files, err := repo.GetFilesByEmailID(emailID, ctx)

		assert.Error(t, err)
		assert.Nil(t, files)
	})
}

func TestDeleteFileByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := EmailRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()

	t.Run("FileDeletedSuccessfully", func(t *testing.T) {
		fileID := uint64(1)

		mock.ExpectExec("DELETE FROM file").
			WithArgs(fileID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteFileByID(fileID, ctx)

		assert.NoError(t, err)
	})

	t.Run("FileDeleteFailedNoRowsAffected", func(t *testing.T) {
		fileID := uint64(2)

		mock.ExpectExec("DELETE FROM file").
			WithArgs(fileID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.DeleteFileByID(fileID, ctx)

		assert.NoError(t, err)
	})

	t.Run("FileDeleteFailedDBError", func(t *testing.T) {
		fileID := uint64(3)

		mock.ExpectExec("DELETE FROM file").
			WithArgs(fileID).
			WillReturnError(fmt.Errorf("database error"))

		err := repo.DeleteFileByID(fileID, ctx)

		assert.Error(t, err)
	})
}

func TestUpdateFileByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := EmailRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()

	t.Run("FileUpdatedSuccessfully", func(t *testing.T) {
		fileID := uint64(1)
		newFileID := "newFileID"
		newFileType := "newFileType"
		newFileName := "PDF"
		newFileSize := "10101010"

		mock.ExpectExec("UPDATE file").
			WithArgs(newFileID, newFileType, newFileName, newFileSize, fileID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateFileByID(fileID, newFileID, newFileType, newFileName, newFileSize, ctx)

		assert.NoError(t, err)
	})

	t.Run("FileUpdateFailedNoRowsAffected", func(t *testing.T) {
		fileID := uint64(2)
		newFileID := "newFileID"
		newFileType := "newFileType"
		newFileName := "PDF"
		newFileSize := "10101010"

		mock.ExpectExec("UPDATE file").
			WithArgs(newFileID, newFileType, newFileName, newFileSize, fileID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateFileByID(fileID, newFileID, newFileType, newFileName, newFileSize, ctx)

		assert.NoError(t, err)
	})

	t.Run("FileUpdateFailedDBError", func(t *testing.T) {
		fileID := uint64(3)
		newFileID := "newFileID"
		newFileType := "newFileType"
		newFileName := "PDF"
		newFileSize := "10101010"

		mock.ExpectExec("UPDATE file").
			WithArgs(newFileID, newFileType, newFileName, newFileSize, fileID).
			WillReturnError(fmt.Errorf("database error"))

		err := repo.UpdateFileByID(fileID, newFileID, newFileType, newFileName, newFileSize, ctx)

		assert.Error(t, err)
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

	repo := EmailRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()

	t.Run("GetAvatarFileIDByLogin_AvatarExists", func(t *testing.T) {
		login := "test_login"
		expectedAvatarID := "avatar_id"

		mock.ExpectQuery(`
            SELECT f\.file_id
            FROM file f
            LEFT JOIN profile p ON p\.avatar_id = f\.id
            WHERE p\.login = \$1
        `).
			WithArgs(login).
			WillReturnRows(sqlmock.NewRows([]string{"file_id"}).AddRow(expectedAvatarID))

		avatarID, err := repo.GetAvatarFileIDByLogin(login, ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedAvatarID, avatarID)
	})

	t.Run("GetAvatarFileIDByLogin_NoAvatar", func(t *testing.T) {
		login := "test_login"

		rows := sqlmock.NewRows([]string{"file_id"})
		mock.ExpectQuery(`
            SELECT f\.file_id
            FROM file f
            LEFT JOIN profile p ON p\.avatar_id = f\.id
            WHERE p\.login = \$1
        `).
			WithArgs(login).
			WillReturnRows(rows)

		avatarID, err := repo.GetAvatarFileIDByLogin(login, ctx)
		assert.Error(t, err)
		assert.Empty(t, avatarID)
	})

	t.Run("GetAvatarFileIDByLogin_DBError", func(t *testing.T) {
		login := "test_login"

		mock.ExpectQuery(`
            SELECT f\.file_id
            FROM file f
            LEFT JOIN profile p ON p\.avatar_id = f\.id
            WHERE p\.login = \$1
        `).
			WithArgs(login).
			WillReturnError(errors.New("DB error"))

		avatarID, err := repo.GetAvatarFileIDByLogin(login, ctx)
		assert.Error(t, err)
		assert.Empty(t, avatarID)
	})
}
