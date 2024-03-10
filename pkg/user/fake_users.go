package user

var FakeUsers = map[uint32]*User{
	1: {
		ID:       1,
		Name:     "Sergey",
		Surname:  "Fedasov",
		Login:    "sergey@mail.ru",
		Password: "$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe",
		AvatarId: "image",
	},
	2: {
		ID:       2,
		Name:     "Ivan",
		Surname:  "Karpov",
		Login:    "ivan@mail.ru",
		Password: "$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe",
		AvatarId: "image",
	},
	3: {ID: 3,
		Name:     "Max",
		Surname:  "Frelich",
		Login:    "max@mail.ru",
		Password: "$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe",
		AvatarId: "image",
	},
}
