-- Импортировать необходимые библиотеки
local wrk = require "wrk"

wrk.method = "POST"
wrk.body = '{"topic": "WRK","text": "WRK_Test","readStatus": false,"mark": false,"replyToEmailId": 0,"draftStatus": false,"spamStatus": false,"senderEmail": "senderWrk@mailhub.su","recipientEmail": "wrk@mailhub.su"}'
wrk.headers["Content-Type"] = "application/json"
wrk.headers["X-Csrf-Token"] = "3b718de1065efe9c8e578ea98452bea8"
wrk.headers["Cookie"] = "session_id=dd001c65474c2617ef0be801e1439336"

-- Указать общее количество запросов
wrk.requests = 100000

-- Настроить параметры потоков
wrk.thread = function()
    -- Указать, что хотим отправить 1000 запросов на каждый поток
    wrk.connections = 400 -- Количество одновременных подключений на поток
    wrk.duration = "1000s" -- Длительность тестирования на поток
    wrk.requests = 20000 -- Количество запросов на поток
end
