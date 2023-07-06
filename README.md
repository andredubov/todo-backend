## REST API to create TODO lists [Backend Application] ![GO][go-badge]

[go-badge]: https://img.shields.io/github/go-mod/go-version/andredubov/todo-backend

## Build & Run (Locally)
### Prerequisites
- go 1.20
- docker & docker-compose
- [swag](https://github.com/swaggo/swag) (<i>optional</i>, used to re-generate swagger documentation)

Create .env file in root directory and add following values:
```dotenv
APP_ENV=local

HTTP_HOST=localhost
HTTP_PORT=8080

DB_HOST=postgres
DB_PORT=5432
DB_NAME=postgres
DB_USER=postgres
DB_PASSWORD=qwerty
DB_SSL_MODE=disable

PASSWORD_SALT=salt
JWT_SIGNING_KEY=key
```

Use `make run` to build&run project.