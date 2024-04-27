-- +migrate Up
-- Создание таблицы вопросов (question)
CREATE TABLE IF NOT EXISTS question (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    text TEXT CHECK (LENGTH(text) <= 200),
    min_text TEXT CHECK (LENGTH(min_text) <= 200),
    max_text TEXT CHECK (LENGTH(max_text) <= 200),
    dop_question TEXT CHECK (LENGTH(dop_question) <= 200)
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
(text, min_text, max_text, dop_question)
VALUES
    ('Насколько вы удовлетворены общим пользовательским опытом при использовании Mailhub?', 'Отлично', 'Очень плохо', 'Что конкретно вам не понравилось в нашем приложении?'),
    ('Как бы вы оценили удобство и интуитивность интерфейса нашего веб-приложения?', 'Очень удобно и интуитивно', 'Совершенно неудобно', 'Предложите какие-либо улучшения или изменения, которые могли бы сделать интерфейс более удобным.'),
    ('Пожалуйста, оцените качество поддержки клиентов, предоставляемой нами.', 'Очень высокое качество поддержки', 'Очень низкое качество поддержки', 'В чем вы видите наши сильные и слабые стороны в области клиентской поддержки?'),
    ('Насколько эффективно работает функционал отправки и получения сообщений?', 'Очень эффективно', 'Совершенно неэффективно', 'Есть ли какие-то особенности функционала отправки и получения сообщений, которые вы хотели бы изменить или доработать?'),
    ('Удовлетворены ли вы скоростью работы Mailhub и быстродействием его функций?', 'Полностью удовлетворен(а)', 'Совершенно не удовлетворен(а)', 'На каких этапах использования приложения вы столкнулись с задержками или низкой производительностью?');

-- Вставка начальных данных в таблицу answer
INSERT INTO answer
(question_id, login, mark)
VALUES
    (1, 'ivan@mailhub.su', 4),
    (2, 'ivan@mailhub.su', 3),
    (3, 'ivan@mailhub.su', 5),
    (4, 'ivan@mailhub.su', 4),
    (5, 'ivan@mailhub.su', 3),
    (1, 'max@mailhub.su', 5),
    (2, 'max@mailhub.su', 4),
    (3, 'max@mailhub.su', 3),
    (4, 'max@mailhub.su', 4),
    (5, 'max@mailhub.su', 5),
    (1, 'serega@mailhub.su', 2),
    (2, 'serega@mailhub.su', 1),
    (3, 'serega@mailhub.su', 5),
    (4, 'serega@mailhub.su', 3),
    (5, 'serega@mailhub.su', 2);