# trial-fiber

[Fiber 开发框架](https://docs.gofiber.io/)的模板项目，可用于快速选型、微服务单元等。

需要确保系统中已安装 [Golang](https://go.dev/) 。

## Usage

```bash
# Create .env file
cp .env_template .env

# Install dependencies
go get

# Update dependencies
go get ./...
go mod tidy

# Run
go run .

# Compile
go build
```

## [Optional] Setup Mirror for Mainland China

You can skip this one if not approaching the internet from within mainland China.

```bash
go env -w GOPROXY=https://goproxy.io,direct # 官方
go env -w GOPROXY=https://goproxy.cn,direct # 七牛云

go env GOPROXY # 确认信息
```
