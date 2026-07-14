#syntax=docker/dockerfile:1.7

# Stage 1: Build templ + Go binary
FROM golang:1.26 AS builder

ENV GOEXPERIMENT=jsonv2
ENV CGO_ENABLED=0

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@v0.3.1020 \
    && templ generate ./examples/demo/... \
    && go build -ldflags="-s -w" -o /demo ./examples/demo/

# Stage 2: Distroless static (no shell, no package manager, ~2MB)
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /demo /demo

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/demo"]
