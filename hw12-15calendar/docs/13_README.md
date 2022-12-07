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
$ curl -i -X POST 'http://127.0.0.1:8088/users' -H "Content-Type: application/json" -d '{"name":"Bob","email":"bob@mail.com"}'
$ curl -i -X GET 'http://127.0.0.1:8088/users?id=1'

$ curl -i -X PUT 'http://127.0.0.1:8088/users' -H "Content-Type: application/json" -d '{"id":1, "name":"Alis","email":"bob@mail.com"}'
$ curl -i -X GET 'http://127.0.0.1:8088/users?id=1'

$ curl -i -X POST 'http://127.0.0.1:8088/events' -H "Content-Type: application/json" -d '{
	"title":"Aliss Birthday",
	"description":"Birthday April 12, 1996. House party",
	"start_time":"2022-05-25T10:41:31Z",
	"stop_time":"2022-05-25T14:41:31Z",
	"user_id":1}'
$ curl -i -X GET 'http://127.0.0.1:8088/events?userid=1&id=1'

$ curl -i -X PUT 'http://127.0.0.1:8088/events' -H "Content-Type: application/json" -d '{
	"id":1,
	"title":"Bob Birthday",
	"description":"Birthday April 17, 1996. House party",
	"start_time":"2022-05-25T10:41:31Z",
	"stop_time":"2022-05-25T14:41:31Z",
	"user_id":1}'
$ curl -i -X GET 'http://127.0.0.1:8088/events?userid=1&id=1'

$ curl -i -X DELETE 'http://127.0.0.1:8088/events?id=1'
$ curl -i -X GET 'http://127.0.0.1:8088/events?userid=1&id=1'

$ curl -i -X DELETE 'http://127.0.0.1:8088/users?id=1'
$ curl -i -X GET 'http://127.0.0.1:8088/users?id=1'
```

### Ссылки:
- [Список кодов состояния HTTP](https://ru.wikipedia.org/wiki/Список_кодов_состояния_HTTP)
