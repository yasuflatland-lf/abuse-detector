# Abuse Detector
This application is for previnting phishing sites are created on Studio.

## Requirements
- Go 1.14.9 >=
- Docker
- Docker Compose

## How to build
```shell script
go build *.go
```    
## How to run
```shell script
go run *.go
```
## How to build Docker image
This is how to build and confirm the image is built correctly.
```
docker build -t studio-abuse-detector .
docker run -p 3000:3000 -d --name studio-abuse-detector studio-abuse-detector:latest
curl localhost:3000
```

## How to remove all images
```
docker rm -f `docker ps -qa`
```