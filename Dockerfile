# Copyright 2025-2026 The MathWorks, Inc.

# Stage 1: Build the proxy executable
FROM golang:1.26.1 AS builder
WORKDIR /app
COPY . /app
RUN go version
RUN ls
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o proxy /app/main.go

# Stage 2: Build the proxy image
FROM scratch
LABEL maintainer="The MathWorks"
COPY --from=builder /app/proxy /proxy

# Add license files
COPY --from=builder /app/licenses/ /licenses

ENTRYPOINT ["./proxy"]
