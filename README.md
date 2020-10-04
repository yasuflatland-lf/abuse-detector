# Abuse Detector
This application is for preventing phishing sites are created on Studio.

## Requirements
- Go 1.14.9 >=
- Docker
- Docker Compose

## How to build
```shell script
go build *.go
```    
## How to run for Development
```shell script
go run main.go
```
How to run production
```
docker-compose up
```
Request verification
```
http://localhost:3000/verify?url=https://www.google.com/
```

## How to build Docker image
This is how to build and confirm the image is built correctly.
```
docker build -t studio-abuse-detector .
docker run -p 3000:3000 -d --name studio-abuse-detector studio-abuse-detector:latest
curl localhost:3000
```

## How to run Chrome Headless server at local
```
docker run -d -p 9222:9222 --rm --name headless-shell --shm-size 2G chromedp/headless-shell
```
## How to remove all images
```~~~~
docker rm -f `docker ps -qa`
```

## Appendix
- [cdp, Chrome Dev Tools Protocl](https://github.com/mafredri/cdp)
- [Headless Chrome server base for Dockerfile, Zenika/alpine-chrome](https://github.com/Zenika/alpine-chrome)