FROM golang:alpine AS base
WORKDIR /app
# [Optional] Setup Mirror for Mainland China
RUN go env -w GOPROXY=https://goproxy.cn,direct # 七牛云
RUN go env GOPROXY # 确认信息

FROM base AS builder
COPY . .
RUN go clean --modcache
RUN go get
RUN go mod verify
RUN go build -o main

FROM base AS runner
COPY --from=builder /app/main ./
COPY .env ./

EXPOSE 3000

CMD ["./main"]