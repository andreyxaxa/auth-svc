# Auth service

## Обзор

Часть сервиса аутентификации. 

JWT-токены, чистая архитектура.


- Swagger - http://localhost:8080/swagger
- Конфиг - [config/config.go](https://github.com/andreyxaxa/auth-svc/blob/main/config/config.go); Читается из `.env` файла.
- Логгер - [pkg/logger/logger.go](https://github.com/andreyxaxa/auth-svc/blob/main/pkg/logger/logger.go); Интерфейс позволяет подменить логгер.
- Dependency injection - [internal/controller/http/v1/controller.go](https://github.com/andreyxaxa/auth-svc/blob/main/internal/controller/http/v1/controller.go), [internal/usecase/session/session.go](https://github.com/andreyxaxa/auth-svc/blob/main/internal/usecase/session/session.go)
- Graceful shutdown - [internal/app/app.go](https://github.com/andreyxaxa/auth-svc/blob/main/internal/app/app.go).
- Удобная и гибкая конфигурация HTTP сервера - [pkg/httpserver/options.go](https://github.com/andreyxaxa/auth-svc/blob/main/pkg/httpserver/options.go).
  Позволяет конфигурировать сервер в конструкторе таким образом:
  ```go
  httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port))
  ```
- В слое хэндлеров применяется версионирование - [internal/controller/http/v1](https://github.com/andreyxaxa/auth-svc/tree/main/internal/controller/http/v1).
  Для версии v2 нужно будет просто добавить папку `http/v2` с таким же содержимым, в файле [internal/controller/http/router.go](https://github.com/andreyxaxa/auth-svc/blob/main/internal/controller/http/router.go) добавить строку:
  ```go
  {
      v1.NewSessionRoutes(apiV1Group, s, l)
  }

  {
      v2.NewSessionRoutes(apiV1Group, s, l)
  }
  ```

## Запуск

Клонируем репозиторий, выполняем:
```
make compose-up
```

## Прочие `make` команды
Зависимости:
```
make deps
```
docker compose down:
```
make compose-down
```
