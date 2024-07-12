FROM golang:1.22-alpine AS build-env

RUN apk update && apk add --no-cache gcc musl-dev make

# Создание рабочей директории
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download -x

COPY . .

# Установка переменной окружения для пути к базе данных
ENV SQLITE_PATH=/app/data/sqlite/storage.db

ENV CGO_ENABLED=1

RUN sh .github/docker-build.sh

FROM alpine:latest

WORKDIR /
COPY --from=build-env /app/build/app /app

# Создание директории для базы данных в контейнере
RUN mkdir -p /data/sqlite

CMD ["/app"]