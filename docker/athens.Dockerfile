FROM gomods/athens:v0.15.1

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*
RUN apk update && apk --no-cache add ca-certificates
RUN apk add bash

RUN mkdir -p /var/lib/athens

EXPOSE 3000

ENTRYPOINT [ "/sbin/tini", "--" ]

CMD ["athens-proxy"]
