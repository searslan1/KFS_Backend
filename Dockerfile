# Base image olarak Go'yu kullanıyoruz
FROM golang:1.22-alpine AS builder

# Çalışma dizinini ayarla
WORKDIR /app

# Bağımlılıkları daha verimli yönetmek için mod ve sum dosyalarını ekle
COPY go.mod go.sum ./
RUN go mod download

# Tüm projeyi kopyala
COPY . .

# Test araçlarını yükle
RUN go install gotest.tools/gotestsum@latest

# Binary build et
RUN go build -o server ./cmd/server

# Unit Test aşaması için ayrı bir target
FROM builder AS tester
CMD ["gotestsum", "--format", "standard-verbose", "./internal/...", "-cover", "-v"]

# Runtime aşaması için yeni bir image oluştur
FROM alpine:latest

WORKDIR /app

# Builder aşamasından binary'yi kopyala
COPY --from=builder /app/server .

# Environment variables'ları kullanmak için gerekli
ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}
ENV DB_NAME=${DB_NAME}
ENV DB_USER=${DB_USER}
ENV DB_PASSWORD=${DB_PASSWORD}

# Uygulamayı çalıştır
CMD ["./server"]
