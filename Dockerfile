# Build stage
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64

RUN go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go build -o /bin/customer-api ./cmd/main.go

# Development stage
FROM golang:1.24.2-alpine AS dev

WORKDIR /app

COPY --from=builder /go/bin/migrate /bin/
COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/air-verse/air@latest
RUN go install go.uber.org/mock/mockgen@latest

EXPOSE 8080

CMD ["air"]

# Production stage
FROM alpine:latest AS prod

WORKDIR /app

RUN echo "**** Production stage ****"
COPY --from=builder /bin/customer-api /bin/
COPY --from=builder /go/bin/migrate /bin/
COPY ./internal/infra/db/postgres/migrations ./migrations/

EXPOSE 8080

CMD ["/bin/customer-api"]