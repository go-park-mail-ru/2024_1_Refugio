package email

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	repo = CreateFakeEmails()
)

func TestGetAll(t *testing.T) {
	emails, err := repo.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(emails); i++ {
		if *emails[i] != *FakeEmails[i] {
			t.Error("emails[i] != FakeEmails[i]")
			assert.Equal(t, FakeEmails[i], emails[i])
		}
	}
}

func TestGetByID(t *testing.T) {
	for i := 1; i*2 < len(FakeEmails); i *= 2 {
		email, err := repo.GetByID(uint64(i))
		if err != nil {
			fmt.Println(err)
			return
		}
		if *email != *FakeEmails[i-1] {
			t.Error("email != FakeEmails[i]")
			assert.Equal(t, FakeEmails[i], email)
			return
		}
	}
}

func TestAdd(t *testing.T) {
	repo = NewEmailMemoryRepository()
	arrEmails := []*Email{
		{
			ID:             1,
			Topic:          "Приглашение на мероприятие",
			Text:           "Уважаемый Гость, приглашаем вас на наше предстоящее мероприятие. Будем рады видеть вас в числе участников!",
			PhotoID:        "https://bizprinting.ru/image/catalog/products/full_1913.jpg",
			ReadStatus:     false,
			Mark:           "",
			Deleted:        false,
			DateOfDispatch: time.Date(2023, time.December, 1, 12, 0, 0, 0, time.UTC),
			ReplyToEmailID: 0,
			DraftStatus:    false,
		},
		{
			ID:             2,
			Topic:          "Акция: Скидка 50% на новую коллекцию",
			Text:           "Уважаемый клиент, представляем вам нашу новую коллекцию. Только сегодня у вас есть уникальная возможность получить скидку 50% на все товары!",
			PhotoID:        "https://tc-edelweiss.ru/wp-content/uploads/2020/11/zenden-50-discont-new-collection.jpg",
			ReadStatus:     true,
			Mark:           "",
			Deleted:        false,
			DateOfDispatch: time.Date(2024, time.January, 2, 10, 30, 0, 0, time.UTC),
			ReplyToEmailID: 0,
			DraftStatus:    false,
		},
		{
			ID:             3,
			Topic:          "Подтверждение заказа",
			Text:           "Ваш заказ №123456 успешно подтвержден. Ожидайте доставку в ближайшие дни. Спасибо за покупку!",
			PhotoID:        "",
			ReadStatus:     true,
			Mark:           "",
			Deleted:        false,
			DateOfDispatch: time.Date(2024, time.January, 5, 15, 0, 0, 0, time.UTC),
			ReplyToEmailID: 0,
			DraftStatus:    false,
		},
	}
	for _, email := range arrEmails {
		_, err := repo.Add(email)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	mails, err := repo.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(arrEmails); i++ {
		if *arrEmails[i] != *mails[i] {
			t.Error("arrEmails[i] != mails[i]")
			assert.Equal(t, FakeEmails[i], mails[i])
			return
		}
	}
}

func TestUpdate(t *testing.T) {
	arrEmails := []*Email{
		{
			ID:             1,
			Topic:          "",
			Text:           "Уважаемый Гость, приглашаем вас на наше предстоящее мероприятие. Будем рады видеть вас в числе участников!",
			PhotoID:        "https://bizprinting.ru/image/catalog/products/full_1913.jpg",
			ReadStatus:     false,
			Mark:           "",
			Deleted:        false,
			DateOfDispatch: time.Date(2023, time.December, 1, 12, 0, 0, 0, time.UTC),
			ReplyToEmailID: 0,
			DraftStatus:    false,
		},
		{
			ID:             2,
			Topic:          "",
			Text:           "Уважаемый клиент, представляем вам нашу новую коллекцию. Только сегодня у вас есть уникальная возможность получить скидку 50% на все товары!",
			PhotoID:        "https://tc-edelweiss.ru/wp-content/uploads/2020/11/zenden-50-discont-new-collection.jpg",
			ReadStatus:     true,
			Mark:           "",
			Deleted:        false,
			DateOfDispatch: time.Date(2024, time.January, 2, 10, 30, 0, 0, time.UTC),
			ReplyToEmailID: 0,
			DraftStatus:    false,
		},
		{
			ID:             3,
			Topic:          "",
			Text:           "Ваш заказ №123456 успешно подтвержден. Ожидайте доставку в ближайшие дни. Спасибо за покупку!",
			PhotoID:        "",
			ReadStatus:     true,
			Mark:           "",
			Deleted:        false,
			DateOfDispatch: time.Date(2024, time.January, 5, 15, 0, 0, 0, time.UTC),
			ReplyToEmailID: 0,
			DraftStatus:    false,
		},
	}

	for _, email := range arrEmails {
		b, err := repo.Update(email)
		if !b {
			t.Error("This ID not found")
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	mails, err := repo.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(arrEmails); i++ {
		if *arrEmails[i] != *mails[i] {
			t.Error("arrEmails[i] != mails[i]")
			assert.Equal(t, FakeEmails[i], mails[i])
			return
		}
	}
}

func TestDelete(t *testing.T) {
	repo = CreateFakeEmails()
	expectedEmails := []*Email{
		{
			ID:             1,
			Topic:          "Приглашение на мероприятие",
			Text:           "Уважаемый Гость, приглашаем вас на наше предстоящее мероприятие. Будем рады видеть вас в числе участников!",
			PhotoID:        "https://bizprinting.ru/image/catalog/products/full_1913.jpg",
			ReadStatus:     false,
			Mark:           "",
			Deleted:        false,
			DateOfDispatch: time.Date(2023, time.December, 1, 12, 0, 0, 0, time.UTC),
			ReplyToEmailID: 0,
			DraftStatus:    false,
		},
		{
			ID:             2,
			Topic:          "Акция: Скидка 50% на новую коллекцию",
			Text:           "Уважаемый клиент, представляем вам нашу новую коллекцию. Только сегодня у вас есть уникальная возможность получить скидку 50% на все товары!",
			PhotoID:        "https://tc-edelweiss.ru/wp-content/uploads/2020/11/zenden-50-discont-new-collection.jpg",
			ReadStatus:     true,
			Mark:           "",
			Deleted:        false,
			DateOfDispatch: time.Date(2024, time.January, 2, 10, 30, 0, 0, time.UTC),
			ReplyToEmailID: 0,
			DraftStatus:    false,
		},
		{
			ID:             3,
			Topic:          "Подтверждение заказа",
			Text:           "Ваш заказ №123456 успешно подтвержден. Ожидайте доставку в ближайшие дни. Спасибо за покупку!",
			PhotoID:        "",
			ReadStatus:     true,
			Mark:           "",
			Deleted:        false,
			DateOfDispatch: time.Date(2024, time.January, 5, 15, 0, 0, 0, time.UTC),
			ReplyToEmailID: 0,
			DraftStatus:    false,
		},
		{
			ID:             4,
			Topic:          "Поздравляем с днем рождения!",
			Text:           "Дорогой друг, от всей души поздравляем тебя с днем рождения! Пусть этот день будет полон радости и удивительных сюрпризов!",
			PhotoID:        "https://img.redzhina.ru/img/ff/39/ff398e29d9ab85278b84089047722f77.jpg",
			ReadStatus:     false,
			Mark:           "",
			Deleted:        false,
			DateOfDispatch: time.Date(2024, time.January, 27, 8, 0, 0, 0, time.UTC),
			ReplyToEmailID: 0,
			DraftStatus:    false,
		},
		{
			ID:             5,
			Topic:          "Новости о продукте",
			Text:           "Мы рады представить вам последние новости о нашем продукте. Новые функции, улучшения и многое другое ждут вас!",
			PhotoID:        "",
			ReadStatus:     true,
			Mark:           "",
			Deleted:        false,
			DateOfDispatch: time.Date(2024, time.February, 5, 14, 45, 0, 0, time.UTC),
			ReplyToEmailID: 0,
			DraftStatus:    false,
		},
	}
	deleteId := []int{6, 7, 8, 9, 10, 11, 12, 13, 14, 16, 17}

	for _, id := range deleteId {
		b, err := repo.Delete(uint64(id))
		if !b {
			t.Error("This ID not found")
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	mails, err := repo.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(expectedEmails); i++ {
		if *expectedEmails[i] != *mails[i] {
			t.Error("arrEmails[i] != mails[i]")
			assert.Equal(t, FakeEmails[i], mails[i])
			return
		}
	}
}
