# Chirpy

A small Twitter-like HTTP API in Go, backed by Postgres. Users post short messages ("chirps"), authenticate with JWTs, and can be upgraded to a "Chirpy Red" membership through a payment provider webhook.

## Stack

Go stdlib `net/http`, Postgres, [sqlc](https://sqlc.dev) for type-safe queries, [goose](https://github.com/pressly/goose) for migrations, [devenv](https://devenv.sh) + direnv for the dev environment.

## Setup

With direnv installed, `cd` into the repo and everything loads automatically: Go, Postgres, sqlc, goose, and the variables from `.env`.

Create a `.env` file:

```sh
DB_URL="postgres://USER:PASSWORD@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
JWT_SECRET="your-secret"
POLKA_KEY="your-polka-api-key"
```

Start Postgres, run the migrations, then start the server:

```sh
devenv up                                    # in a separate terminal
goose -dir sql/schema postgres "$DB_URL" up
go run ./cmd/api                             # listens on :8080
```

## Endpoints

| Method | Path | Auth | Description |
| --- | --- | --- | --- |
| GET | `/api/healthz` | none | Health check |
| POST | `/api/users` | none | Create a user |
| PUT | `/api/users` | Bearer | Update email and password |
| POST | `/api/login` | none | Returns a user, an access token, and a refresh token |
| POST | `/api/refresh` | Bearer (refresh) | Returns a new access token |
| POST | `/api/revoke` | Bearer (refresh) | Revokes a refresh token |
| POST | `/api/chirps` | Bearer | Create a chirp, max 140 chars |
| GET | `/api/chirps` | none | List chirps |
| GET | `/api/chirps/{chirpID}` | none | Get one chirp |
| DELETE | `/api/chirps/{chirpID}` | Bearer | Delete your own chirp |
| POST | `/api/polka/webhooks` | ApiKey | Upgrade a user to Chirpy Red |
| GET | `/admin/metrics` | none | Fileserver hit count |
| POST | `/admin/reset` | none | Deletes all users, `dev` platform only |

`GET /api/chirps` takes two optional query parameters:

- `author_id`: a user UUID, filters to that author's chirps (filtered in SQL)
- `sort`: `asc` (default) or `desc`, by `created_at`

Access tokens are sent as `Authorization: Bearer <token>` and last 1 hour, refresh tokens last 60 days. The Polka webhook expects `Authorization: ApiKey <key>` and returns 401 if it does not match `POLKA_KEY`.

## Development

```sh
go test ./...                                # unit tests (auth package)
sqlc generate                                # regenerate internal/database after editing sql/queries
goose -dir sql/schema postgres "$DB_URL" up  # apply new migrations
```

Migrations live in `sql/schema`, queries in `sql/queries`. Never edit `internal/database` by hand, it is generated.

## Layout

```
cmd/api/           handlers, routing, config
internal/auth/     password hashing, JWTs, header parsing
internal/database/ sqlc-generated code
sql/schema/        goose migrations
sql/queries/       sqlc queries
```
