-- +migrate Up
-- Создание таблицы вложений (file)
CREATE TABLE IF NOT EXISTS file (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    file_id TEXT CHECK (LENGTH(file_id) <= 200),
    file_type TEXT CHECK (LENGTH(file_id) <= 200) NOT NULL DEFAULT '',
    file_name TEXT CHECK (LENGTH(file_id) <= 200) NOT NULL DEFAULT '',
    file_size TEXT CHECK (LENGTH(file_id) <= 200) NOT NULL DEFAULT ''
);

-- Создание таблицы пользователей (profile)
CREATE TABLE IF NOT EXISTS profile (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    login TEXT NOT NULL UNIQUE CHECK (LENGTH(login) <= 50),
    password_hash TEXT NOT NULL CHECK (LENGTH(password_hash) <= 200),
    firstname TEXT NOT NULL CHECK (LENGTH(firstname) <= 50),
    surname TEXT NOT NULL CHECK (LENGTH(surname) <= 50),
    patronymic TEXT CHECK (LENGTH(patronymic) <= 50),
    gender TEXT NOT NULL CHECK (gender = 'Male' OR gender = 'Female' OR gender = 'Other'),
    birthday DATE,
    registration_date DATE NOT NULL DEFAULT CURRENT_DATE,
    avatar_id INTEGER REFERENCES file(id) ON DELETE NO ACTION DEFAULT NULL,
    phone_number TEXT CHECK (LENGTH(phone_number) <= 20),
    description TEXT CHECK (LENGTH(description) <= 300),
    vkid INTEGER DEFAULT 0
);

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
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    topic TEXT,
    text TEXT,
    date_of_dispatch TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    photoid TEXT CHECK (LENGTH(photoid) <= 200),
    sender_email TEXT NOT NULL CHECK (LENGTH(sender_email) <= 50),
    recipient_email TEXT NOT NULL CHECK (LENGTH(recipient_email) <= 50),
    isRead BOOLEAN NOT NULL,
    isDeleted BOOLEAN NOT NULL,
    isDraft BOOLEAN NOT NULL,
    isSpam BOOLEAN NOT NULL,
    reply_to_email_id INTEGER REFERENCES email(id) ON DELETE NO ACTION DEFAULT NULL,
    is_important BOOLEAN NOT NULL
);

-- Создание таблицы вложений (file)
CREATE TABLE IF NOT EXISTS email_file (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    email_id INT REFERENCES email(id) ON DELETE CASCADE,
    file_id INT REFERENCES file(id) ON DELETE CASCADE
);

-- Создание таблицы письма пользователя (profile_email)
CREATE TABLE IF NOT EXISTS profile_email (
    profile_id INT,
    email_id INT,
    PRIMARY KEY ( profile_id, email_id ),
    CONSTRAINT fk_profile FOREIGN KEY (profile_id) REFERENCES profile(id) ON DELETE CASCADE,
    CONSTRAINT fk_email FOREIGN KEY (email_id)  REFERENCES email(id) ON DELETE CASCADE
);

-- Создание таблицы папок (folder)
CREATE TABLE IF NOT EXISTS folder (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    profile_id INTEGER REFERENCES profile(id) ON DELETE CASCADE,
    name TEXT NOT NULL CHECK (LENGTH(name) <= 100)
);

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
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    profile_id INTEGER REFERENCES profile(id) ON DELETE CASCADE,
    notifications_enabled BOOLEAN,
    language TEXT CHECK (LENGTH(language) <= 50)
);

-- Вставка начальных данных в таблицу file
INSERT INTO file
(file_id, file_type)
VALUES
    ('', 'PHOTO'), ('', 'PHOTO'), ('', 'PHOTO'), ('', 'PHOTO');

-- Вставка начальных данных в таблицу users
INSERT INTO profile
(login, password_hash, firstname, surname, patronymic, gender, birthday, registration_date, phone_number, description, avatar_id, vkid)
VALUES
    ('sergey@mailhub.su', '$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe', 'Sergey', 'Fedasov', 'Aleksandrovich', 'Male', '2003-08-20', NOW(), '+77777777777', 'Description', 1, 0),
    ('ivan@mailhub.su', '$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe', 'Ivan', 'Karpov', 'Aleksandrovich', 'Male', '2003-10-17', NOW(), '+79697045539', 'Description', 2, 0),
    ('max@mailhub.su', '$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe', 'Maxim', 'Frelich', 'Aleksandrovich', 'Male', '2003-08-20', NOW(), '+79099099090', 'Description', 3, 0),
    ('alex@mailhub.su', '$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe', 'Alexey', 'Khochevnikov', 'Aleksandrovich', 'Male', '2003-10-20', NOW(), '+79090007030', 'Description', 4, 0),
    ('fedasov@mailhub.su', '$2a$10$4PcooWbEMRjvdk2cMFumO.ajWaAclawIljtlfu2.2f5/fV8LkgEZe', 'Сергей', 'Федасов', 'Андреевич', 'Male', '2003-10-20', NOW(), '+79090007030', 'Description', 4, 344167564)
ON CONFLICT (login) DO NOTHING;

-- Вставка начальных данных в таблицу email
INSERT INTO email
(topic, text, date_of_dispatch, photoid, sender_email, recipient_email, isRead, isDeleted, isDraft, isSpam, reply_to_email_id, is_important)
VALUES
    ('Topic1 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 'sergey@mailhub.su', 'ivan@mailhub.su', False, False, False, False, Null, False),
    ('Topic2 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 'sergey@mailhub.su', 'max@mailhub.su', False, False, False, False, Null, False),
    ('Topic3 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 'ivan@mailhub.su', 'max@mailhub.su', False, False, False, False, Null, False),
    ('Topic4 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 'max@mailhub.su', 'sergey@mailhub.su', False, False, False, False, Null, False),
    ('Topic5 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 'max@mailhub.su', 'ivan@mailhub.su', False, False, False, False, Null, False),
    ('Topic6 Enough pretended estimating.', 'Laughing say assurance indulgence mean unlocked stairs denote above prudent get use latter margaret. Unreserved another abode blushes old steepest lady disposing enjoyment immediate prevailed charm. Looked ladies civil sigh. Because cold offended quiet bred the. Hastened outlived supported.', '2022-08-10 10:10:00', '', 'sergey@mailhub.su', 'ivan@mailhub.su', False, False, False, False, Null, False);

-- Вставка начальных данных в таблицу profile_email
INSERT INTO profile_email
(profile_id, email_id)
VALUES
    (1, 1), (2, 1), (1, 2), (3, 2), (2, 3), (3, 3), (3, 4), (1, 4), (3, 5), (2, 5), (1, 6), (2, 6);

-- Вставка начальных данных в таблицу email_file
INSERT INTO email_file
(email_id, file_id)
VALUES
    (1, 1), (2, 1), (3, 2), (4, 3), (5, 3), (6, 1);