# trial-fiber

[Fiber 开发框架](https://docs.gofiber.io/)的模板项目，可用于快速选型、微服务单元等。

需要确保系统中已安装 [Golang](https://go.dev/) 。

## [Optional] Setup Mirror for Mainland China

You can skip this one if not approaching the internet from within mainland China.

```bash
go env -w GOPROXY=https://goproxy.io,direct # 官方
go env -w GOPROXY=https://goproxy.cn,direct # 七牛云

go env GOPROXY # 确认信息
```

## Usage

```bash
# Create .env file
cp .env_template .env # specify database connection info

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