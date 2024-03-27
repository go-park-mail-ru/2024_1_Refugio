-- +migrate Up
-- Создание таблицы пользователей (users)
CREATE TABLE IF NOT EXISTS users (
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

-- Создание таблицы сессий (sessions)
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY CHECK (char_length(id) <= 50),
    user_id INTEGER REFERENCES users(id),
    creation_date TIMESTAMP,
    device TEXT CHECK (LENGTH(device) <= 100),
    life_time INTEGER,
    csrf_token TEXT CHECK (char_length(csrf_token) <= 50)
);

-- Вставка начальных данных в таблицу users
INSERT INTO users
    (id, login, password, firstname, surname, patronymic, gender, birthday, registration_date, avatar_id, phone_number, description)
VALUES
    (1, 'sergey@mailhub.ru', '$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe', 'Sergey', 'Fedasov', 'Aleksandrovich', 'Male', '2003-08-20', NOW(), '', '+77777777777', 'Description'),
    (2, 'ivan@mailhub.ru', '$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe', 'Ivan', 'Karpov', 'Aleksandrovich', 'Male', '2003-10-17', NOW(), '', '+79697045539', 'Description'),
    (3, 'max@mailhub.ru', '$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe', 'Maxim', 'Frelich', 'Aleksandrovich', 'Male', '2003-08-20', NOW(), '', '+79099099090', 'Description')
ON CONFLICT (login) DO NOTHING;

-- Создание таблицы писем (emails)
CREATE TABLE IF NOT EXISTS emails (
    id SERIAL PRIMARY KEY,
    topic TEXT,
    text TEXT,
    date_of_dispatch DATE,      /*TIMESTAMPTZ NOT NULL DEFAULT '2022-08-10 10:10:00',*/
    photoid TEXT CHECK (LENGTH(photoid) <= 200),
    sender_id INTEGER REFERENCES users(id),
    recipient_id INTEGER REFERENCES users(id),
    read_status BOOLEAN NOT NULL,
    deleted_status BOOLEAN NOT NULL,
    draft_status BOOLEAN NOT NULL,
    reply_to_email_id INTEGER REFERENCES emails(id) ON DELETE NO ACTION DEFAULT NULL,
    flag BOOLEAN NOT NULL
);

-- Вставка начальных данных в таблицу emails
INSERT INTO emails
    (id, topic, text, date_of_dispatch, photoid, sender_id, recipient_id, read_status, deleted_status, draft_status, reply_to_email_id, flag)
VALUES
    (1, 'Topic1 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 1, 2, False, False, False, Null, False),
    (2, 'Topic2 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 1, 3, False, False, False, Null, False),
    (3, 'Topic3 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 2, 3, False, False, False, Null, False),
    (4, 'Topic4 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 3, 1, False, False, False, Null, False),
    (5, 'Topic5 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 3, 2, False, False, False, Null, False),
    (6, 'Topic6 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 1, 2, False, False, False, Null, False)
;
