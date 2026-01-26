# Stage 1: Base image and downloading dependencies
FROM golang:1.22 AS base
MAINTAINER Van Le <vanlt@phx-smartschool.com>

RUN mkdir -p /go/src/backend
WORKDIR /go/src/backend

COPY go.mod .
COPY go.sum .

RUN go mod download

# Stage 2: Build the go app
FROM base AS builder

# Sử dụng ARG để cấu hình GOPROXY và SERVICE_PATH
ARG GOPROXY=https://goproxy.cn
ARG SERVICE_PATH
ARG BRANCH_NAME
ARG SERVICE_TYPE
COPY . /go/src/backend
WORKDIR /go/src/backend/${SERVICE_PATH}

# Build the binary
RUN if [ "$SERVICE_TYPE" = "api" ]; then \
    go install github.com/swaggo/swag/cmd/swag@v1.16.3 && \
    swag init --parseDependency -g swagger_${BRANCH_NAME}.go; \
    else \
        echo "Skipping Swagger generation"; \
    fi && \
    GOPROXY=$GOPROXY mkdir -p bin && go build -o ./bin/service ./.


# Stage 3: Create the final image
FROM debian:stable-slim
MAINTAINER Van Le <vanlt@phx-smartschool.com>

ARG SERVICE_PATH
ARG BRANCH_NAME

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        netbase && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get autoremove -y && apt-get autoclean -y

# Sử dụng ARG để chỉ định SERVICE_PATH

COPY --from=builder /go/src/backend/${SERVICE_PATH}/bin /app
COPY --from=builder /go/src/backend/configs /app/configs

WORKDIR /app