FROM golang:1.24.2-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o /app/customer-api ./cmd/customer/main.go

EXPOSE 8080

# Defina o comando para executar as APIs
CMD ["./customer-api"]