# Workouts

Workouts is a simple Go applicaton that allows you to track your workouts.
Built on a [FEM](https://frontendmasters.com/courses/complete-go) course.

## Prerequisites

- [Docker](https://docs.docker.com/desktop)
- [Go](https://go.dev/doc/install)
- [Goose](https://github.com/pressly/goose)
- [PostgreSQL](https://www.postgresql.org/download) _(optional, for debugging)_

## Run locally

This app uses PostgreSQL as a database through Docker. To connect to a
database, simply run:

```bash
docker compose up -d --build
```

This will start a PostgreSQL database on port 5432. To inspect changes to it
you can run:

```bash
psql -h localhost -p 5432 -U postgres
```

Once your databse is connected, run migrations with:

```bash
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
```

Start the application by running:

```bash
go run main.go
```

## Testing

This application spins up a PostgreSQL test database on port 5433. Same as with
the main database, to inspect changes to it you can run:

```bash
psql -h localhost -p 5433 -U postgres
```
