# Gator: a blog aggregator in Go

## Docker

```sh
docker run -d \
  --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=gator \
  -p 5432:5432 \
  postgres:18
```

## Goose

```sh
goose -dir sql/schema postgres "postgres://postgres:password@127.0.0.1:5432/gator?sslmode=disable" up
```
