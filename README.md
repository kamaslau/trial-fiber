# trial-fiber

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
![Repository size](https://img.shields.io/github/repo-size/kamaslau/trial-fiber?color=56BEB8)

Template [Fiber framework](https://docs.gofiber.io/) project for fast prototyping, trial, or micro-service unit usage.

Make sure that you have [Golang](https://go.dev/) installed already.

- ~~[Removed] [GQLGen](https://gqlgen.com/) as GraphQL~~
- [GORM](https://gorm.io/docs/) as ORM, to [operate with Postgres database](https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL)
- [Redis](https://redis.io/docs/latest/) as Cache
- [NATS](https://docs.nats.io/) as Message Queue
- [InfluxDB](https://docs.influxdata.com/) for Time Series Data (logs, etc.)

## [🇨🇳 Optional] Setup Mirror for Mainland China

You can skip this one if not approaching the internet from within mainland China.

```bash
go env -w GOPROXY=https://goproxy.io,direct # Official
go env -w GOPROXY=https://goproxy.cn,direct # QiniuCloud/七牛云

go env GOPROXY # Double check
```

## Service Endpoints

### RESTful

- Root: http://127.0.0.1:3000/api/v1/
- Demo with _Posts_ model: http://127.0.0.1:3000/api/v1/post

### GraphQL (TODO)

- Root: [http://127.0.0.1:3000/graphql/](http://127.0.0.1:3000/graphql/)
- Demo is not available for now

```bash
curl -X POST \
http://127.0.0.1:3000/graphql/ \
-H 'Content-Type: application/json' \
-d '{
  "query": "query { posts { id name } }"
}'
```

## Usage

Create a `trial-fiber` database (you can use another name, just config it in the `.env` file) in Postgres/MySQL/MariaDB of yours, then follow these steps below.

```bash
# Create .env file
cp .env.sample .env # at least specify a database connection info

# Install dependencies
go get ./src

# [Optional] Update dependencies
go get -u ./src
go mod tidy

# Run
## With live-reloading (via air-verse/air)
### Install and config air
go install github.com/air-verse/air@latest
air init
### Change the .air.toml file generated
#### macOS/Linux
[build]
bin = "./tmp/main"
cmd = "go build -o ./tmp/main ./src"
#### Windows
[build]
bin = "tmp\\main.exe"
cmd = "go build -o ./tmp/main.exe ./src"
### Use
air

## Without live-reloading
go run src/main.go

# Compile
## Remember to put a .env file in to the same directory with the executable file compiled
go build -o ./dist/main ./src
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
- [Building a GraphQL Server with Go Backend Tutorial | Getting Started](https://www.howtographql.com/graphql-go/0-introduction/)
