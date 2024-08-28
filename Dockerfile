# Указываем базовый образ с Go, используем последний стабильный образ
FROM golang:1.21-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app
RUN apk add --no-cache gcc musl-dev

# Копируем go.mod и go.sum, чтобы установить зависимости
COPY go.mod go.sum ./
ENV CGO_ENABLED=1
# Загружаем зависимости
RUN go mod download

# Копируем все файлы проекта в рабочую директорию
COPY . .

# Компилируем приложение
RUN go build -o main .

# Указываем команду для запуска приложения
CMD ["./main"]
