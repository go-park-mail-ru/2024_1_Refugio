-- Импортировать необходимые библиотеки
local wrk = require "wrk"

wrk.method = "POST"
wrk.body = '{"firstKey": "somedata", "secondKey": "somedata"}'
wrk.body = '{
    "topic": "WRK",
    "text": "WRK_Test",
    "readStatus": false,
    "mark": false,
    "replyToEmailId": 0,
    "draftStatus": false,
    "spamStatus": false,
    "senderEmail": "wrk@mailhub.su",
    "recipientEmail": "ivan@mailhub.su"}'
wrk.headers["Content-Type"] = "application/json"
wrk.headers["X-Csrf-Token"] = "e5bf01a0faf9177e583c49b0db26e10f"
wrk.headers["session_id"] = "4ec9f9423c727e9e701ad2f2357503b7"

wrk.requests = 10