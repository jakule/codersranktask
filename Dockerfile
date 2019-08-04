FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git
COPY . /build
WORKDIR /build
ENV GO111MODULES=on
RUN go get
RUN go build

FROM alpine:latest
RUN apk --no-cache --update add ca-certificates
COPY --from=builder  /build/codersranktask /codersranktask

CMD ["/codersranktask"]
