-- +migrate Up
-- Создание таблицы пользователей (profile)
CREATE TABLE IF NOT EXISTS profile (
    id SERIAL PRIMARY KEY,
    login TEXT UNIQUE CHECK (LENGTH(login) <= 50),
    password TEXT CHECK (LENGTH(password) <= 200),
    firstname TEXT CHECK (LENGTH(firstname) <= 50),
    surname TEXT CHECK (LENGTH(surname) <= 50),
    patronymic TEXT CHECK (LENGTH(patronymic) <= 50),
    gender TEXT CHECK (LENGTH(gender) <= 10),
    birthday DATE,
    registration_date DATE,
    avatar_id TEXT CHECK (LENGTH(avatar_id) <= 200),
    phone_number TEXT CHECK (LENGTH(phone_number) <= 20),
    description TEXT CHECK (LENGTH(description) <= 300)
);

-- Создание таблицы сессий (session)
CREATE TABLE IF NOT EXISTS session (
    id TEXT PRIMARY KEY CHECK (char_length(id) <= 50),
    profile_id INTEGER REFERENCES profile(id) ON DELETE CASCADE,
    creation_date TIMESTAMP,
    device TEXT CHECK (LENGTH(device) <= 100),
    life_time INTEGER,
    csrf_token TEXT CHECK (char_length(csrf_token) <= 50)
);

-- Вставка начальных данных в таблицу users
INSERT INTO profile
    (id, login, password, firstname, surname, patronymic, gender, birthday, registration_date, avatar_id, phone_number, description)
VALUES
    (1, 'sergey@mailhub.ru', '$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe', 'Sergey', 'Fedasov', 'Aleksandrovich', 'Male', '2003-08-20', NOW(), '', '+77777777777', 'Description'),
    (2, 'ivan@mailhub.ru', '$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe', 'Ivan', 'Karpov', 'Aleksandrovich', 'Male', '2003-10-17', NOW(), '', '+79697045539', 'Description'),
    (3, 'max@mailhub.ru', '$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe', 'Maxim', 'Frelich', 'Aleksandrovich', 'Male', '2003-08-20', NOW(), '', '+79099099090', 'Description')
ON CONFLICT (login) DO NOTHING;
