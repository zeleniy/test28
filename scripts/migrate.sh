#!/bin/sh

migrate -database "$DB_URL" -path "/app/db/migrations" "$@"
