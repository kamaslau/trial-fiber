FROM golang:alpine AS base
WORKDIR /app
RUN go env -w GOPROXY=https://goproxy.cn,direct # 七牛云
RUN go env GOPROXY # 确认信息

FROM base AS builder
COPY . .
RUN go get
RUN go mod verify
RUN go build -o trial-fiber

FROM base AS runner
COPY --from=builder /app/trial-fiber ./
COPY .env ./

EXPOSE 3000

CMD ["./trial-fiber"]