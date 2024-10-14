FROM golang:1.22-alpine AS builder

RUN apk update && apk upgrade

WORKDIR /usr/src/app
COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY .. .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app ./cmd/app/main.go
CMD ["./bin/app"]

#FROM scratch
#
#COPY --from=builder /usr/src/app/bin/ .
#COPY --from=builder /usr/src/app/.env.local .
#
#CMD ["./app"]
