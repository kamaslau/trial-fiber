# trial-fiber

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
![Repository size](https://img.shields.io/github/repo-size/kamaslau/trial-fiber?color=56BEB8)

Template [Fiber framework](https://docs.gofiber.io/) project for fast prototyping, trial, or micro-service unit usage.

Make sure that you have [Golang](https://go.dev/) installed already.

- [GORM](https://gorm.io/docs/) as ORM, to [operate with Postgre database](https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL)
- [Redis](https://redis.io/docs/latest/) as Cache, [etcd](https://etcd.io/docs/latest/) can also be a good choice
- [NATS](https://docs.nats.io/) as Message Queue
- [InfluxDB](https://docs.influxdata.com/) for Time Series Data (logs, etc.)
- ~~[GQLGen](https://gqlgen.com/) as GraphQL~~

## [🇨🇳 Optional] Setup Mirror for Mainland China

You can safely skip this part if not approaching the internet from within mainland China.

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

Create a `trial-fiber` database (you can use another name, just config it in the `.env` file) in PostgreSQL/MySQL/MariaDB of yours, then follow these steps below.

```bash
# Create .env file
cp .env.sample .env # at least specify a database connection info

# Install dependencies
go get ./...

# [Optional] Update dependencies
go get -u ./...
go mod tidy
```

### Run
Without live-reloading
```bash
go run ./src/...
```
With live-reloading (via air-verse/air)
```bash
# Install Air
go install github.com/air-verse/air@latest

# Run with
air # Linux/macOS
air -c .air.windows.toml # Windows

```
For more detailed instructions, e.g. solutions on working cross OS, checkout [https://manual.kamaslau.com/golang/tools/air.html](https://manual.kamaslau.com/golang/tools/air.html).


### Compile

Remember to put a .env file in to the same directory with the executable file compiled
```bash
go build -o ./dist/main ./src

```

## Deploy (via docker)

```bash
docker build . -t trial-fiber:latest

docker network create trial-backend && \
docker stop trial-fiber && \
docker rm trial-fiber && \
docker run --name trial-fiber --network trial-backend -p 3000:3000 -d --restart always --net=host trial-fiber:latest
```

Make sure that other service components (database, cache, MQ, etc.) are within the same network, e.g.,

```bash
docker network connect trial-backend trial-fiber
docker network connect trial-backend database-xxx
docker network connect trial-backend mq-xxx
```

## References/Credits

- [Go Fiber: Start Building RESTful APIs on Golang (Feat. GORM)](https://dev.to/percoguru/getting-started-with-apis-in-golang-feat-fiber-and-gorm-2n34)
- [Building a GraphQL Server with Go Backend Tutorial | Getting Started](https://www.howtographql.com/graphql-go/0-introduction/)
