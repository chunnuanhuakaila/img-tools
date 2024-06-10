FROM golang:1.19-alpine
WORKDIR $GOPATH/src/search
COPY . .
COPY etc /app/etc
RUN go env -w GOPROXY=https://goproxy.cn,direct && go build -o /app/img && go clean --modcache
FROM alpine:latest
#RUN apk update && apk add curl
WORKDIR /app
COPY --from=0 /app .
ENTRYPOINT ["./img"]