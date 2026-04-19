FROM golang:1.21-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /api
WORKDIR /api

RUN touch $HOME/.netrc

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o ./router ./main.go

FROM alpine:latest

RUN apk update && apk add ca-certificates bash && rm -rf /var/cache/apk/*

RUN mkdir -p /api
WORKDIR /api
COPY --from=builder /api/router .

EXPOSE 8080

ENTRYPOINT ["./router"]
