-- +migrate Up
-- Создание таблицы пользователей (profile)
CREATE TABLE IF NOT EXISTS profile (
    id INTEGER PRIMARY KEY,
    login TEXT NOT NULL UNIQUE CHECK (LENGTH(login) <= 50),
    password_hash TEXT NOT NULL CHECK (LENGTH(password_hash) <= 200),
    firstname TEXT NOT NULL CHECK (LENGTH(firstname) <= 50),
    surname TEXT NOT NULL CHECK (LENGTH(surname) <= 50),
    patronymic TEXT CHECK (LENGTH(patronymic) <= 50),
    gender TEXT NOT NULL CHECK (gender = 'Male' OR gender = 'Female' OR gender = 'Other'),
    birthday DATE,
    registration_date DATE NOT NULL DEFAULT CURRENT_DATE,
    file_id TEXT CHECK (LENGTH(file_id) <= 200),
    phone_number TEXT CHECK (LENGTH(phone_number) <= 20),
    description TEXT CHECK (LENGTH(description) <= 300)
);

-- Создание Sequence последовательности для profile.id
CREATE SEQUENCE profileId
START 1
INCREMENT 1
OWNED BY profile.id;

-- Создание таблицы сессий (session)
CREATE TABLE IF NOT EXISTS session (
    id TEXT PRIMARY KEY CHECK (LENGTH(id) <= 50),
    profile_id INTEGER REFERENCES profile(id) ON DELETE CASCADE,
    creation_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    device TEXT CHECK (LENGTH(device) <= 100),
    life_time INTEGER NOT NULL,
    csrf_token TEXT NOT NULL CHECK (LENGTH(csrf_token) <= 50)
);

-- Создание таблицы писем (email)
CREATE TABLE IF NOT EXISTS email (
    id INTEGER PRIMARY KEY,
    topic TEXT,
    text TEXT,
    date_of_dispatch TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    photoid TEXT CHECK (LENGTH(photoid) <= 200),
    sender_email TEXT NOT NULL CHECK (LENGTH(sender_email) <= 50),
    recipient_email TEXT NOT NULL CHECK (LENGTH(recipient_email) <= 50),
    isRead BOOLEAN NOT NULL,
    isDeleted BOOLEAN NOT NULL,
    isDraft BOOLEAN NOT NULL,
    reply_to_email_id INTEGER REFERENCES email(id) ON DELETE NO ACTION DEFAULT NULL,
    is_important BOOLEAN NOT NULL
);

-- Создание Sequence последовательности для email.id
CREATE SEQUENCE emailId
START 1
INCREMENT 1
OWNED BY email.id;

-- Создание таблицы письма пользователя (profile_email)
CREATE TABLE IF NOT EXISTS profile_email (
    profile_id INT,
    email_id INT,
    PRIMARY KEY ( profile_id, email_id ),
    CONSTRAINT fk_profile FOREIGN KEY (profile_id) REFERENCES profile(id),
    CONSTRAINT fk_email FOREIGN KEY (email_id)  REFERENCES email(id)
);

-- Создание таблицы вложений (file)
CREATE TABLE IF NOT EXISTS file (
    id INTEGER PRIMARY KEY,
    email_id INT REFERENCES email(id) ON DELETE CASCADE,
    file_id TEXT CHECK (LENGTH(file_id) <= 200)
);

-- Создание Sequence последовательности для file.id
CREATE SEQUENCE fileId
START 1
INCREMENT 1
OWNED BY file.id;

-- Создание таблицы папок (folder)
CREATE TABLE IF NOT EXISTS folder (
    id INTEGER PRIMARY KEY,
    profile_id INTEGER REFERENCES profile(id) ON DELETE CASCADE,
    name TEXT NOT NULL CHECK (LENGTH(name) <= 100)
);

-- Создание Sequence последовательности для folder.id
CREATE SEQUENCE folderId
START 1
INCREMENT 1
OWNED BY folder.id;

-- Создание таблицы связи папок с письмами (folder_email)
CREATE TABLE IF NOT EXISTS folder_email (
    folder_id INTEGER,
    email_id INTEGER,
    PRIMARY KEY ( folder_id, email_id ),
    CONSTRAINT fk_folder FOREIGN KEY (folder_id) REFERENCES folder(id) ON DELETE CASCADE,
    CONSTRAINT fk_email FOREIGN KEY  (email_id) REFERENCES email(id) ON DELETE CASCADE
);

-- Создание таблицы настроек (settings)
CREATE TABLE IF NOT EXISTS settings (
    id INTEGER PRIMARY KEY,
    profile_id INTEGER REFERENCES profile(id) ON DELETE CASCADE,
    notifications_enabled BOOLEAN,
    language TEXT CHECK (LENGTH(language) <= 50)
);

-- Создание Sequence последовательности для settings.id
CREATE SEQUENCE settingsId
START 1
INCREMENT 1
OWNED BY settings.id;

-- Вставка начальных данных в таблицу users
INSERT INTO profile
    (id, login, password_hash, firstname, surname, patronymic, gender, birthday, registration_date, file_id, phone_number, description)
VALUES
    (nextval('profileId'), 'sergey@mailhub.su', '$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe', 'Sergey', 'Fedasov', 'Aleksandrovich', 'Male', '2003-08-20', NOW(), '', '+77777777777', 'Description'),
    (nextval('profileId'), 'ivan@mailhub.su', '$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe', 'Ivan', 'Karpov', 'Aleksandrovich', 'Male', '2003-10-17', NOW(), '', '+79697045539', 'Description'),
    (nextval('profileId'), 'max@mailhub.su', '$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe', 'Maxim', 'Frelich', 'Aleksandrovich', 'Male', '2003-08-20', NOW(), '', '+79099099090', 'Description')
ON CONFLICT (login) DO NOTHING;

-- Вставка начальных данных в таблицу email
INSERT INTO email
    (id, topic, text, date_of_dispatch, photoid, sender_email, recipient_email, isRead, isDeleted, isDraft, reply_to_email_id, is_important)
VALUES
    (nextval('emailId'), 'Topic1 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 'sergey@mailhub.su', 'ivan@mailhub.su', False, False, False, Null, False),
    (nextval('emailId'), 'Topic2 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 'sergey@mailhub.su', 'max@mailhub.su', False, False, False, Null, False),
    (nextval('emailId'), 'Topic3 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 'ivan@mailhub.su', 'max@mailhub.su', False, False, False, Null, False),
    (nextval('emailId'), 'Topic4 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 'max@mailhub.su', 'sergey@mailhub.su', False, False, False, Null, False),
    (nextval('emailId'), 'Topic5 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 'max@mailhub.su', 'ivan@mailhub.su', False, False, False, Null, False),
    (nextval('emailId'), 'Topic6 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 'sergey@mailhub.su', 'ivan@mailhub.su', False, False, False, Null, False);

-- Вставка начальных данных в таблицу profile_email
INSERT INTO profile_email
    (profile_id, email_id)
VALUES
    (1, 1), (2, 1), (1, 2), (3, 2), (2, 3), (3, 3), (3, 4), (1, 4), (3, 5), (2, 5), (1, 6), (2, 6);

