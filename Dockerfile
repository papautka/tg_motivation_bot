# Этап сборки (builder stage)
# Используем официальный образ Golang с минимальной ОС Alpine
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы зависимостей Go в контейнер
COPY go.mod go.sum ./

# Загружаем зависимости проекта
RUN go mod download

# Копируем остальной код проекта в контейнер
COPY . .

# Компилируем Go-приложение
# CGO_ENABLED=0 — отключаем C-биндинги для статической сборки
# GOOS=linux — указываем целевую платформу
# Сборка бинарника с именем `bot` из директории `cmd`
RUN CGO_ENABLED=0 GOOS=linux go build -x -o bot ./cmd


# Финальный образ (production stage)
# Используем минимальный образ Alpine без Golang
FROM alpine:latest

# Копируем скомпилированный бинарник из builder-образа
COPY --from=builder /app/bot /bot

# Указываем команду по умолчанию — запуск бота
CMD ["/bot"]
