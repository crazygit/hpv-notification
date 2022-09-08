FROM golang:1.19-alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

ARG REGION=US
ARG VERSION=unknown
ARG GIT_COMMIT_HASH=unknown
ARG BUILD_TIME=unknown


WORKDIR /build
RUN apk add --no-cache --virtual .build-deps \
    upx \
    ca-certificates \
    gcc \
    g++

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -a -o hpv-notification main.go && upx hpv

FROM alpine:3
WORKDIR /app
COPY ./config.yaml /app/config.yaml
COPY --from=builder /build/hpv-notification .

CMD ["/app/hpv-notification", "fetch"]
