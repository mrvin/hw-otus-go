## Домашнее задание №13 «API к Календарю»
Необходимо реализовать HTTP и GRPC API для сервиса календаря.

Методы API в принципе идентичны методам хранилища и [описаны в ТЗ](./CALENDAR.MD).

Для GRPC API необходимо:
* создать отдельную директорию для Protobuf спецификаций;
* создать Protobuf файлы с описанием всех методов API, объектов запросов и ответов (
т.к. объект Event будет использоваться во многих ответах разумно выделить его в отдельный message);
* создать отдельный пакет для кода GRPC сервера;
* добавить в Makefile команду `generate`; `make generate` - вызывает `go generate`, которая в свою очередь
генерирует код GRPC сервера на основе Protobuf спецификаций;
* написать код, связывающий GRPC сервер с методами доменной области (бизнес логикой);
* логировать каждый запрос по аналогии с HTTP API.

Для HTTP API необходимо:
* расширить "hello-world" сервер из [ДЗ №12](./12_README.md) до полноценного API;
* создать отдельный пакет для кода HTTP сервера;
* реализовать хэндлеры, при необходимости выделив структуры запросов и ответов;
* сохранить логирование запросов, реализованное в [ДЗ №12](./12_README.md).

Общие требования:
* должны быть реализованы все методы;
* календарь не должен зависеть от кода серверов;
* сервера должны запускаться на портах, указанных в конфиге сервиса.

**Можно использовать https://grpc-ecosystem.github.io/grpc-gateway/.**

### Критерии оценки
- Makefile заполнен и пайплайн зеленый - 1 балл
- Реализовано GRPC API и `make generate` - 3 балла
- Реализовано HTTP API - 2 балла
- Написаны юнит-тесты на API - до 2 баллов
- Понятность и чистота кода - до 2 баллов

#### Зачёт от 7 баллов

### Заметки
Клиент для http API
```bash
$ curl -ik --cert ../cert/clientCert.pem --key ../cert/clientKey.pem -X POST 'https://127.0.0.1:8088/signup' -H "Content-Type: application/json" -d '{"user_name":"Bob","password":"qwerty","email":"bob@mail.com"}'
$ curl -ik --cert ../cert/clientCert.pem --key ../cert/clientKey.pem -X GET 'https://127.0.0.1:8088/login' -H "Content-Type: application/json" -d '{"user_name":"Bob","password":"qwerty"}'
$ curl -ik --cert ../cert/clientCert.pem --key ../cert/clientKey.pem -X GET 'https://127.0.0.1:8088/user' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDAwNjA5MDYsImlhdCI6MTcwMDA2MDAwNiwidXNlcm5hbWUiOiJCb2IifQ.B2_m_bQFs4Q1zMRHISskFN6fZNGUscSOnj8-5Ce0Qug"

$ curl -ik --cert ../cert/clientCert.pem --key ../cert/clientKey.pem -X PUT 'https://127.0.0.1:8088/user' -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDAwNjA5MDYsImlhdCI6MTcwMDA2MDAwNiwidXNlcm5hbWUiOiJCb2IifQ.B2_m_bQFs4Q1zMRHISskFN6fZNGUscSOnj8-5Ce0Qug" -d '{"user_name":"Bob","password":"123456","email":"bobiiiii@mail.com"}'

$ curl -ik --cert ../cert/clientCert.pem --key ../cert/clientKey.pem -X GET 'https://127.0.0.1:8088/user' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDAwNjA5MDYsImlhdCI6MTcwMDA2MDAwNiwidXNlcm5hbWUiOiJCb2IifQ.B2_m_bQFs4Q1zMRHISskFN6fZNGUscSOnj8-5Ce0Qug"

$ curl -ik --cert ../cert/clientCert.pem --key ../cert/clientKey.pem -X POST 'https://127.0.0.1:8088/event' -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDAwNjA5MDYsImlhdCI6MTcwMDA2MDAwNiwidXNlcm5hbWUiOiJCb2IifQ.B2_m_bQFs4Q1zMRHISskFN6fZNGUscSOnj8-5Ce0Qug" -d '{
	"title":"Aliss Birthday",
	"description":"Birthday April 12, 1996. House party",
	"start_time":"2022-05-25T10:41:31Z",
	"stop_time":"2022-05-25T14:41:31Z",
	"user_id":1}'
$ curl -ik --cert ../cert/clientCert.pem --key ../cert/clientKey.pem -X GET 'https://127.0.0.1:8088/event?id=1' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDAwNjA5MDYsImlhdCI6MTcwMDA2MDAwNiwidXNlcm5hbWUiOiJCb2IifQ.B2_m_bQFs4Q1zMRHISskFN6fZNGUscSOnj8-5Ce0Qug"

$ curl -ik --cert ../cert/clientCert.pem --key ../cert/clientKey.pem -X PUT 'https://127.0.0.1:8088/event' -H "Content-Type: application/json"  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDAwNjA5MDYsImlhdCI6MTcwMDA2MDAwNiwidXNlcm5hbWUiOiJCb2IifQ.B2_m_bQFs4Q1zMRHISskFN6fZNGUscSOnj8-5Ce0Qug" -d '{
	"id":1,
	"title":"Bob Birthday",
	"description":"Birthday April 17, 1996. House party",
	"start_time":"2022-05-25T10:41:31Z",
	"stop_time":"2022-05-25T14:41:31Z",
	"user_id":1}'
	
$ curl -ik --cert ../cert/clientCert.pem --key ../cert/clientKey.pem -X GET 'https://127.0.0.1:8088/event?id=1' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDAwNjA5MDYsImlhdCI6MTcwMDA2MDAwNiwidXNlcm5hbWUiOiJCb2IifQ.B2_m_bQFs4Q1zMRHISskFN6fZNGUscSOnj8-5Ce0Qug"

$ curl -i -X DELETE 'http://127.0.0.1:8088/event?id=1' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTk4MDcyNTAsImlhdCI6MTY5OTgwNjM1MCwidXNlcm5hbWUiOiJCb2IifQ.0g0WPM57GVCmi83SHi7mok8zgTlQrPSP8xAtzUVoha4"

$ curl -i -X GET 'http://127.0.0.1:8088/event?id=1' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTk4MDcyNTAsImlhdCI6MTY5OTgwNjM1MCwidXNlcm5hbWUiOiJCb2IifQ.0g0WPM57GVCmi83SHi7mok8zgTlQrPSP8xAtzUVoha4"

$ curl -i -X DELETE 'http://127.0.0.1:8088/user' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTk4MDcyNTAsImlhdCI6MTY5OTgwNjM1MCwidXNlcm5hbWUiOiJCb2IifQ.0g0WPM57GVCmi83SHi7mok8zgTlQrPSP8xAtzUVoha4"

$ curl -i -X GET 'http://127.0.0.1:8088/users' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTk4MDcyNTAsImlhdCI6MTY5OTgwNjM1MCwidXNlcm5hbWUiOiJCb2IifQ.0g0WPM57GVCmi83SHi7mok8zgTlQrPSP8xAtzUVoha4"
```

### Ссылки:
- [Список кодов состояния HTTP](https://ru.wikipedia.org/wiki/Список_кодов_состояния_HTTP)
