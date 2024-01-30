# Используем официальный образ Golang
FROM golang:latest

# Устанавливаем пакеты, включая PostgreSQL client
RUN apt-get update && apt-get install -y postgresql-client

# Копируем все файлы вашего приложения в контейнер
COPY . /go/src/app

# Устанавливаем рабочую директорию
WORKDIR /go/src/app

# Собираем приложение
RUN go build -o main ./cmd/app

# Экспортируем порт, который слушает приложение
EXPOSE 8080

# Команда, которая будет запущена при старте контейнера
CMD ["./main"]
