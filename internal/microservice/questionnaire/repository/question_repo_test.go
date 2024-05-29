package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"mail/internal/microservice/models/domain_models"
	"mail/internal/pkg/logger"
	"mail/internal/pkg/utils/constants"
)

func GetCTX() context.Context {
	ctx := context.WithValue(context.Background(), interface{}(string(constants.LoggerKey)), logger.InitializationBdLog(nil))
	ctx2 := context.WithValue(ctx, interface{}(string(constants.RequestIDKey)), []string{"testID"})

	return ctx2
}

func TestQuestionAnswerRepository_GetAllQuestions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := QuestionAnswerRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()

	t.Run("Success", func(t *testing.T) {
		expectedQuestions := []*domain_models.Question{
			{ID: 1, Text: "Question 1", MinResult: "Min Text 1", MaxResult: "Max Text 1", DopQuestion: "Dop Question 1"},
			{ID: 2, Text: "Question 2", MinResult: "Min Text 2", MaxResult: "Max Text 2", DopQuestion: "Dop Question 2"},
		}

		rows := sqlmock.NewRows([]string{"id", "text", "min_text", "max_text", "dop_question"}).
			AddRow(1, "Question 1", "Min Text 1", "Max Text 1", "Dop Question 1").
			AddRow(2, "Question 2", "Min Text 2", "Max Text 2", "Dop Question 2")
		query := "SELECT question.id, question.text, question.min_text, question.max_text, question.dop_question FROM question"
		mock.ExpectQuery(query).WillReturnRows(rows)

		questions, err := repo.GetAllQuestions(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedQuestions, questions)
	})

	t.Run("DBError", func(t *testing.T) {
		query := "SELECT question.id, question.text, question.min_text, question.max_text, question.dop_question FROM question"
		mock.ExpectQuery(query).WillReturnError(sql.ErrConnDone)

		questions, err := repo.GetAllQuestions(ctx)

		assert.Error(t, err)
		assert.Nil(t, questions)
	})

	t.Run("NoRows", func(t *testing.T) {
		query := "SELECT question.id, question.text, question.min_text, question.max_text, question.dop_question FROM question"
		mock.ExpectQuery(query).WillReturnError(sql.ErrNoRows)

		questions, err := repo.GetAllQuestions(ctx)

		assert.Error(t, err)
		assert.Nil(t, questions)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestQuestionAnswerRepository_GetAllAnswers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := QuestionAnswerRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()

	t.Run("Success", func(t *testing.T) {
		expectedAnswers := []*domain_models.Answer{
			{ID: 1, QuestionID: 1, Login: "user1", Mark: 4, Text: "Answer 1"},
			{ID: 2, QuestionID: 2, Login: "user2", Mark: 3, Text: "Answer 2"},
		}

		rows := sqlmock.NewRows([]string{"id", "question_id", "login", "mark", "text"}).
			AddRow(1, 1, "user1", 4, "Answer 1").
			AddRow(2, 2, "user2", 3, "Answer 2")
		query := "SELECT answer.id, answer.question_id, answer.login, answer.mark, question.text FROM answer LEFT JOIN question ON question.id = answer.question_id"
		mock.ExpectQuery(query).WillReturnRows(rows)

		answers, err := repo.GetAllAnswers(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedAnswers, answers)
	})

	t.Run("DBError", func(t *testing.T) {
		query := "SELECT answer.id, answer.question_id, answer.login, answer.mark, question.text FROM answer LEFT JOIN question ON question.id = answer.question_id"
		mock.ExpectQuery(query).WillReturnError(sql.ErrConnDone)

		answers, err := repo.GetAllAnswers(ctx)

		assert.Error(t, err)
		assert.Nil(t, answers)
	})

	t.Run("NoRows", func(t *testing.T) {
		query := "SELECT answer.id, answer.question_id, answer.login, answer.mark, question.text FROM answer LEFT JOIN question ON question.id = answer.question_id"
		mock.ExpectQuery(query).WillReturnError(sql.ErrNoRows)

		answers, err := repo.GetAllAnswers(ctx)

		assert.Error(t, err)
		assert.Nil(t, answers)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestQuestionAnswerRepository_AddQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := QuestionAnswerRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()

	t.Run("Success", func(t *testing.T) {
		expectedQuestion := &domain_models.Question{
			ID:          1,
			Text:        "Test Question",
			MinResult:   "Min Text",
			MaxResult:   "Max Text",
			DopQuestion: "Additional Question",
		}

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		insertQuery := `INSERT INTO question \(text, min_text, max_text, dop_question\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`
		mock.ExpectQuery(insertQuery).WillReturnRows(rows)

		result, err := repo.AddQuestion(expectedQuestion, ctx)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("DBError", func(t *testing.T) {
		insertQuery := `INSERT INTO question \(text, min_text, max_text, dop_question\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`
		mock.ExpectQuery(insertQuery).WillReturnError(sql.ErrConnDone)

		result, err := repo.AddQuestion(&domain_models.Question{}, ctx)

		assert.Error(t, err)
		assert.False(t, result)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestQuestionAnswerRepository_AddAnswer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := QuestionAnswerRepository{
		DB: sqlx.NewDb(mockDB, "sqlmock"),
	}

	ctx := GetCTX()

	t.Run("Success", func(t *testing.T) {
		expectedAnswer := &domain_models.Answer{
			QuestionID: 1,
			Login:      "testuser",
			Mark:       4,
			Text:       "Test answer",
		}

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		insertQuery := `INSERT INTO answer \(question_id, login, mark, text\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`
		mock.ExpectQuery(insertQuery).WillReturnRows(rows)

		result, err := repo.AddAnswer(expectedAnswer, ctx)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("DBError", func(t *testing.T) {
		insertQuery := `INSERT INTO answer \(question_id, login, mark, text\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`
		mock.ExpectQuery(insertQuery).WillReturnError(sql.ErrConnDone)

		result, err := repo.AddAnswer(&domain_models.Answer{}, ctx)

		assert.Error(t, err)
		assert.False(t, result)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
