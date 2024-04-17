# trial-fiber

[Fiber 开发框架](https://docs.gofiber.io/)的模板项目，可用于快速选型、微服务单元等。

需要确保系统中已安装 [Golang](https://go.dev/) 。

```bash
# 创建环境配置文件
cp .env_template .env

# 初始化
go mod init trial-fiber
go get github.com/gofiber/fiber/v3

# 开发
go fmt main.go
go run main.go

# 编译
go build

# 部署
go run trial-fiber
```

## 中国大陆镜像

```bash
go env -w GOPROXY=https://goproxy.io,direct // 官方镜像
```
