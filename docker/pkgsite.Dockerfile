FROM golang AS builder

# Clone the pkgsite repository. Using commit from 11/18/2024
RUN git clone "https://go.googlesource.com/pkgsite" && \
cd pkgsite && \
git checkout cfc082e0779d6ed60edb737ed05e0a93bf9e7006 && \
CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -a -installsuffix cgo -o /pkgsite cmd/frontend/main.go

FROM redis:7.4.1-alpine

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates
RUN apk add bash

RUN mkdir -p /data
RUN mkdir -p /conf

COPY ./conf/redis.conf /conf/redis.conf
COPY ./entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

COPY --from=builder /pkgsite /
COPY --from=builder /go/pkgsite/static /static
COPY --from=builder /go/pkgsite/third_party /third_party

EXPOSE 8888
