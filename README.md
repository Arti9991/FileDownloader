# 12-07-2025

curl -v -X POST http://localhost:8080/task

curl -v -X POST -H "Content-Type: application/json" -d '[
{"url":"www.ya.ru"},
{"url":"www.dlya.ru"},
{"url":"www.Nya.ru"}]' http://localhost:8082/task/add