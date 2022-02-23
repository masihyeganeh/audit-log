FROM golang:1.17 as builder

WORKDIR /go/src/masihyeganeh/audit-log
COPY go.mod .
COPY go.sum .

ENV GO111MODULE=on
RUN go mod download

FROM builder as server_builder
WORKDIR /go/src/masihyeganeh/audit-log

COPY . .

RUN GIT_COMMIT=$(git rev-parse --short HEAD) \
 && BUILD_TIME=$(date +%Y-%m-%d-%H:%M:%S) \
 && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.GitCommit=$GIT_COMMIT -X main.BuildTime=$BUILD_TIME" -o server cmd/server/*.go

FROM debian:stretch-slim
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
WORKDIR /usr/local/

COPY --from=server_builder /go/src/masihyeganeh/audit-log/server .
# COPY --from=server_builder /go/src/masihyeganeh/audit-log/configs ./configs # TODO

ENTRYPOINT ["./server"]
