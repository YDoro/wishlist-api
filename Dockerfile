FROM golang:1.24.2-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
 
RUN go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go build -o /bin/customer-api ./cmd/customer/main.go

EXPOSE 8080