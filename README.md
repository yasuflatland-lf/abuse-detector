# Abuse Detector
This application is for preventing phishing sites are created on Studio. 

## Requirements
- Go 1.14.9 >=
- Docker 2.4.0.0 >=
- Docker Compose 1.27.4 >=

## Usage
### Request verification
This API verifies if the site does not include malicious links, such as fishing.
```
http://localhost:3000/verify?url=https://www.google.com/
```
The response would look like below if the site is not malicious.
```
{
    "strategyName": "",
    "link": [],
    "malicious": false,
    "statusCode": 200,
    "error": null
}
```
If it's malicious, the response looks like below.
```
{
    "strategyName": "TransparencyReportVerifyStrategy",
    "link": ["http://sucursalvirtualpersonas-sa.com"],
    "malicious": true,
    "statusCode": 200,
    "error": null
}
```
## How to build
```shell script
go build *.go
```    
## How to run for Development
1. Create `.env` based off from `.env.default`. For API keys required, please refer documents below in this README. 
1. In `.env` file, Remove `production` string from `COMMON_APP_ENV` as follows.
    ```
    COMMON_APP_ENV=
    ```
1. Start Chrome Headless Server
    ```
    docker run -d -p 9222:9222 --rm --name headless-shell --shm-size 2G chromedp/headless-shell
    ```
1. Run server as below. `realize` command allows Hot reloading.
    ```shell script
    realize start
    ```

## How to run for production
1. Create `.env` based off from `.env.default`
1. Set API Keys accordingly.
1. Run command below.
    ```
    docker-compose up
    ```
   
## How to run all tests
```
go test -v -race -run=. -bench=. ./...
```   

## How to build Docker image
This is how to build and confirm the image is built correctly.
```
docker build -t studio-abuse-detector .
docker run -p 3000:3000 -d --name studio-abuse-detector studio-abuse-detector:latest
curl localhost:3000
```

## Opearation Related
### How to remove all images including running
```~~~~
docker rm -f `docker ps -qa`
```
### How to access an image
```
docker-compose exec app /bin/sh
```

## Appendix
- [cdp, Chrome Dev Tools Protocl](https://github.com/mafredri/cdp)
- [Headless Chrome server base for Dockerfile, Zenika/alpine-chrome](https://github.com/Zenika/alpine-chrome)

### How to get API key for urlscan.io
1. Go to `https://urlscan.io/` and create an account.
1. Go to [Settings & API](https://urlscan.io/user/profile/) and create an API Key
1. Copy the `Key` and set it to `URLSCAN_API_KEY` in the `.env` file

### How to get API key for Google Safe Browsing API
1. Access to [Google API Console](https://console.developers.google.com/) and create a project
1. Create API key in the project.
1. Look for `Google Safe Browsing API` in `Liberary` tab and add it for the API Key created.
1. Copy the `Key` and set it to `GOOGLE_SAFE_BROWSING_API_KEY` in the `.env` file

## Caveat
- Chrome Headless server in use may need load balancing for a more massive load of access.
- Test links are real phishing sites for now. They become offline or removed in the short term, so tests highly likely to fail.