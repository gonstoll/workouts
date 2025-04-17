# Workouts

Workouts is a simple Go applicaton that allows you to track your workouts.
Built on a [FEM](https://frontendmasters.com/courses/complete-go) course.

## Run locally

This app uses PostgreSQL as a database through Docker. To connect to a
database, simply run:

```bash
docker compose up -d --build
```

This will start a PostgreSQL database on port 5432. To inspect changes to it
you can run

```bash
psql -h localhost -p 5432 -U postgres
```

Once your database is connected, start the application by running:

```bash
go run main.go
```
