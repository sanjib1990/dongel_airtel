FROM golang:1.17.6-alpine3.15 as builder

ENV CGO_ENABLED=1

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN apk update  \
    && apk add --virtual  \
    build-dependencies  \
    build-base  \
    gcc  \
    && make build

FROM alpine:3.15.0
COPY --from=builder /app/dongel ./
CMD ["./dongel"]
