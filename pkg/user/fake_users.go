package user

var FakeUsers = map[uint32]*User{
	1: {
		ID:       1,
		Name:     "Sergey",
		Surname:  "Fedasov",
		Login:    "sergey@mail.ru",
		Password: "$2a$14$frlAr5Y9N4NWOsUbWH9K.eq1tsOUnoEJnV8dlEnlgEb5s7F6R1CpS",
		AvatarId: "image",
	},
	2: {
		ID:       2,
		Name:     "Ivan",
		Surname:  "Karpov",
		Login:    "ivan@mail.ru",
		Password: "$2a$14$frlAr5Y9N4NWOsUbWH9K.eq1tsOUnoEJnV8dlEnlgEb5s7F6R1CpS",
		AvatarId: "image",
	},
	3: {ID: 3,
		Name:     "Max",
		Surname:  "Frelich",
		Login:    "max@mail.ru",
		Password: "$2a$14$frlAr5Y9N4NWOsUbWH9K.eq1tsOUnoEJnV8dlEnlgEb5s7F6R1CpS",
		AvatarId: "image",
	},
}
