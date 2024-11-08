FROM golang:1.22-alpine AS builder

RUN apk update && apk upgrade

WORKDIR /usr/src/app
COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY .. .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/echo ./cmd/benchmark_echo/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/gin ./cmd/benchmark_gin/main.go

#FROM scratch
#
#COPY --from=builder /usr/src/app/bin/ .
#COPY --from=builder /usr/src/app/.env.local .
#
#CMD ["./app"]
