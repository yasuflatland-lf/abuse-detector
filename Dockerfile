FROM golang:1.14.9-alpine3.12 as build

WORKDIR /go/app

COPY . .
COPY .env .

RUN apk add --no-cache git \
 && go build -o app

FROM alpine:3.12.0

WORKDIR /app

COPY --from=build /go/app/app .

RUN apk add --update --no-cache go git \
  && export GOPATH=/root/go \
  && export PATH=${GOPATH}/bin:/usr/local/go/bin:$PATH \
  && export GOBIN=$GOROOT/bin \
  && mkdir -p ${GOPATH}/src ${GOPATH}/bin \
  && addgroup go \
  && adduser -D -G go go \
  && chown -R go:go /app/app \
  && chmod +x /app/app

CMD ["go", "run", "main.go"]