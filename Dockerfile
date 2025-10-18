FROM golang:1.24.1-alpine

RUN apk add --no-cache postgresql-client curl && \
    go install github.com/mitranim/gow@latest && \
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
    go install github.com/go-task/task/v3/cmd/task@latest && \
    go install github.com/aarondl/sqlboiler/v4@latest && \
    go install github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-psql@latest && \
    go install github.com/stephenafamo/boilingfactory@latest && \
    go install gotest.tools/gotestsum@latest && \
    go install github.com/stephenafamo/boilingseed@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest

WORKDIR /app
COPY go.mod go.sum /app
RUN go mod download -x
