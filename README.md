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
    "strategyName": "TransparencyReportVerifyStrategy",
    "link": ["https://www.google.com/"],
    "malicious": false,
    "statusCode": 200,
    "error": {
        "Name": "google-chrome",
        "Err": {}
    }
}
```
If it's malicious, the response looks like below.
```
{
    "strategyName": "TransparencyReportVerifyStrategy",
    "link": ["https://zonabn1segura-pe.com/"],
    "malicious": true,
    "statusCode": 200,
    "error": null
}
```
## How to build
```shell script
go build *.go
```    
## How to run production
1. Create `.env` based off from `.env.default`. For API keys required, please refer documents below in this README. 
1. Set API Keys accordingly. For API keys required, please refer documents below in this README.
1. Run command below.
    ```
    docker-compose up
    ```
    For the initial start, run as below.
    ```
    docker-compose up --build
    ```
## How to run for development
1. Create `.env` based off from `.env.default`. For API keys required, please refer documents below in this README. 
1. Set API Keys accordingly. For API keys required, please refer documents below in this README.
1. Comment out `CMD ["./app"]` and remove comment of `CMD [ "realize", "start" ]` instead to enable realize for hot reloading.
1. Run command below.
    ```
    docker-compose up
    ```
    For the initial start, run as below.
    ```
    docker-compose up --build
    ```

## How to run for debugging with IDE, such as Goland
1. Create `.env` based off from `.env.default`. For API keys required, please refer documents below in this README.
1. Configure `COMMON_APP_ENV=`, no strings. (Default should be `production`) 
1. Set API Keys accordingly. For API keys required, please refer documents below in this README.
1. Comment out `CMD ["./app"]` and remove comment of `CMD [ "realize", "start" ]` instead to enable realize for hot reloading.
1. Spin up servers as below
   
   Spin up chrome headless server
   ```
   docker run -d -p 9222:9222 --rm --name headless-shell --shm-size 2G chromedp/headless-shell
   ```
1. Then, right click `main.go` and debug run on Goland IDE. 
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