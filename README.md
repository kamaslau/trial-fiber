# trial-fiber

Template [Fiber framework](https://docs.gofiber.io/) project for fast prototyping, trial, or micro-service unit usage.

Make sure that you have [Golang](https://go.dev/) installed already.

## [Optional] Setup Mirror for Mainland China

You can skip this one if not approaching the internet from within mainland China.

```bash
go env -w GOPROXY=https://goproxy.io,direct # Official
go env -w GOPROXY=https://goproxy.cn,direct # QiniuCloud/七牛云

go env GOPROXY # Double check
```

## Usage

```bash
# Create .env file
cp .env.sample .env # specify database connection info

# Install dependencies
go get

# [Optional] Update dependencies
go get -u
go mod tidy

# Run
## With live-reloading (air)
go install github.com/air-verse/air@latest # This line only needs to be run once
air

## Without live-reloading
go run .

# Compile
go build
```

## Deploy with docker

```bash
docker build . -t trial-fiber:latest

docker stop trial-fiber && \
docker rm trial-fiber && \
docker run --name trial-fiber -p 3000:3000 -d --restart always --net=host trial-fiber:latest
```

## References/Credits

- [Go Fiber: Start Building RESTful APIs on Golang (Feat. GORM)](https://dev.to/percoguru/getting-started-with-apis-in-golang-feat-fiber-and-gorm-2n34)
