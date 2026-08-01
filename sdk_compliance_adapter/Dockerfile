FROM golang:1.26.5-alpine AS builder

# The official golang images set GOTOOLCHAIN=local, so a base older than the
# `go` directive in go.mod fails the build outright rather than fetching a
# toolchain. `auto` lets a future go.mod bump resolve itself against the
# checksum database instead of hard-failing here.
ENV GOTOOLCHAIN=auto

WORKDIR /app

# Copy the SDK source
COPY go.mod go.sum ./
COPY *.go ./

# Copy adapter
COPY sdk_compliance_adapter/main.go /app/adapter/

# Download dependencies
RUN go mod download

# Build adapter
RUN cd /app/adapter && go build -o /app/adapter-server main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/adapter-server /app/adapter-server

EXPOSE 8080

CMD ["/app/adapter-server"]
