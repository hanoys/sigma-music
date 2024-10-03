FROM golang:1.22-alpine AS builder

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /usr/src/app
COPY . .

FROM alpine:latest

COPY --from=builder /usr/src/app .
COPY --from=builder /go/bin/migrate .
CMD ./migrate -path . -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up
