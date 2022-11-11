#!/bin/bash
set -e

export $(cat .env.default | grep -v ^# | xargs);
APP_DIR=$(pwd) go run main.go

