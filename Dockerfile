# Base image olarak Go'yu kullanıyoruz
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install gotest.tools/gotestsum@latest

RUN wget -O /usr/local/bin/wait-for-it https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && \
    chmod +x /usr/local/bin/wait-for-it
    
RUN go build -o server ./cmd/server

# ✅ Test aşaması için "tester" stage'i düzeltilmiş hali
FROM builder AS tester
WORKDIR /app
ENTRYPOINT [ "sh", "-c" ]
CMD [ "echo 'Starting Tests...'; gotestsum --format standard-verbose ./internal/... -cover -v" ]


# Runtime aşaması için yeni bir image oluştur
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
CMD ["./server"]
