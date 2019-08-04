FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git
COPY . /build
WORKDIR /build
ENV GO111MODULES=on
RUN go get
RUN go build

FROM alpine:latest
RUN apk --no-cache --update add ca-certificates
COPY sql /app/sql
COPY .env /app/.env
COPY --from=builder  /build/codersranktask /app/codersranktask
WORKDIR /app

CMD ["./codersranktask"]
