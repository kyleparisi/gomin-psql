#!/bin/bash
set -e

export $(cat .env.default | grep -v ^# | xargs);

echo Starting db:
psql -h $DB_HOST --quiet -d postgres -c "DROP DATABASE $DB_DATABASE;" || :
psql -h $DB_HOST --quiet -d postgres -c "CREATE DATABASE $DB_DATABASE;"

echo Running migrations:
./schema/manage.py makemigrations
./schema/manage.py migrate
python generator.py > src/app/types.go

echo Running format:
go fmt

echo Running tests:
APP_DIR=$(pwd) go test -v -p 1 ./...
