FROM golang:1.24.1-alpine

RUN apk add --no-cache postgresql-client && \
    go install github.com/mitranim/gow@latest && \
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
    go install github.com/go-task/task/v3/cmd/task@latest && \
    go install github.com/aarondl/sqlboiler/v4@latest && \
    go install github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-psql@latest

# RUN apk add --no-cache postgresql-client

# RUN chmod +x wait-for-postgres.sh

# CMD ["./wait-for-postgres.sh", "postgres:5432", "--", "air", "-c", ".air.toml"]
