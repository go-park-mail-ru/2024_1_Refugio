-- +migrate Up
-- Создание таблицы вопросов (question)
CREATE TABLE IF NOT EXISTS question (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    text TEXT CHECK (LENGTH(text) <= 200),
    min_text TEXT CHECK (LENGTH(text) <= 200),
    max_text TEXT CHECK (LENGTH(text) <= 200)
);

-- Создание таблицы ответов (answer)
CREATE TABLE IF NOT EXISTS answer (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    question_id INTEGER REFERENCES question(id) ON DELETE CASCADE,
    login TEXT NOT NULL CHECK (LENGTH(login) <= 50),
    mark INTEGER NOT NULL DEFAULT 0,
    text TEXT CHECK (LENGTH(text) <= 200) DEFAULT ''
);

-- Вставка начальных данных в таблицу question
INSERT INTO question
(text, min_text, max_text)
VALUES
    ('Насколько вы удовлетворены нашим продуктом?', 'Сильно неудовлетворен', 'Сильно удовлетворен'),
    ('Вам нравится интерфес сайта?', 'Очень ненравится', 'Очень нравится'),
    ('Много ли вы получаете спама?', 'Очень мало', 'Очень много');