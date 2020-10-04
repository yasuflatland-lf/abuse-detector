FROM golang:1.14.9-alpine3.12 as build

WORKDIR /go/app

COPY . .
COPY .env .
COPY .realize.yaml .

RUN apk add --no-cache git \
 && go build -o app

FROM golang:1.14.9-alpine3.12

WORKDIR /app

COPY --from=build /go/app/app .

RUN apk add --no-cache alpine-sdk git \
  && go get -u github.com/oxequa/realize \
  && addgroup go \
  && adduser -D -G go go \
  && chown -R go:go /app/app \
  && chmod +x /app/app

# CMD ["go", "run", "main.go"]
CMD [ "realize", "start" ]