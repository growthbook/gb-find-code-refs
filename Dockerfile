FROM alpine:3.19.0

RUN apk update
RUN apk add --no-cache git
RUN apk add --no-cache openssh

COPY gb-find-code-refs /usr/local/bin/gb-find-code-refs

ENTRYPOINT ["gb-find-code-refs"]
