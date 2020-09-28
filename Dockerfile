FROM golang:1.14.9-alpine3.12 as build

WORKDIR /go/app

COPY . .
COPY .env .

RUN apk add --no-cache git \
 && go build -o app

FROM alpine:3.12.0

WORKDIR /app

COPY --from=build /go/app/app .

RUN addgroup go \
  && adduser -D -G go go \
  && chown -R go:go /app/app

CMD ["./app"]