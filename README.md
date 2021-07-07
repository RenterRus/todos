# Todos
API Заметок <br>
Использует связку HTTP + gRPC
<br>
Целевое взаимодействие происходит через HTTP
<br>
Полученная команда перенаправляется на GRPC сервер, который в свою очередь дергает методы SQLite (Не БД в памяти)

## Запуск
Программа считывает заданные флаги с дефолтными значениями:<br>
1. `--http` ("127.0.0.1:9999") HTTP Server addr<br>
2. `--grpc` (":9998") GRPC Server addr<br>
3. `--db` ("todos") DB name<br>
4. `--table` ("todo") DB table name<br>

## HTTP API
1. `/GetTodos` - Получает открытые задачи <br>
2. `/Create?message=<text>` - Создает задачу без таймера на закрытие <br>
3. `/Update?id=<id>&message=<text>` - Обновляет описание задачи <br>
4. `/Close?id=<id>` - Закрывает задачу <br>
5. `/SetExpTimeout?id=<id>&date=<7.7.2021>&time=<23:57:11>` - Устанавливает таймер на закрытие