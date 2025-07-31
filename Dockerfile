FROM golang:alpine AS base
WORKDIR /app
# [Optional] Setup Mirror for Mainland China
# RUN go env -w GOPROXY=https://goproxy.cn,direct # 七牛云
# RUN go env GOPROXY

FROM base AS builder
COPY . .
RUN go clean --modcache
RUN go mod tidy
RUN go mod verify
RUN go build -o ./tmp/main ./src

FROM base AS runner
COPY --from=builder /app/tmp/main ./
COPY .env ./

EXPOSE 3000

CMD ["./main"]