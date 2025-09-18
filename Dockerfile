FROM golang:1.24.1-alpine

RUN go install github.com/mitranim/gow@latest
# RUN apk add --no-cache postgresql-client

# RUN chmod +x wait-for-postgres.sh

# CMD ["./wait-for-postgres.sh", "postgres:5432", "--", "air", "-c", ".air.toml"]
