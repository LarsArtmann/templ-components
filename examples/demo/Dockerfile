#syntax=docker/dockerfile:1.7

# Stage 1: Compile Tailwind CSS (always fresh, never stale)
# Scans all .templ files for class names and generates the final CSS.
# This eliminates the staleness risk of committing a static app.css.
FROM node:22-slim AS css
WORKDIR /src
COPY . .
RUN npm init -y \
    && npm install tailwindcss @tailwindcss/cli \
    && npx @tailwindcss/cli \
        -i examples/demo/demo.css \
        -o examples/demo/static/app.css \
        --minify

# Stage 2: Build templ + Go binary
FROM golang:1.26 AS builder

ENV GOEXPERIMENT=jsonv2
ENV CGO_ENABLED=0

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Overwrite committed CSS with freshly compiled CSS from Stage 1
COPY --from=css /src/examples/demo/static/app.css examples/demo/static/app.css

RUN go install github.com/a-h/templ/cmd/templ@v0.3.1020 \
    && templ generate ./examples/demo/... \
    && go build -ldflags="-s -w" -o /demo ./examples/demo/

# Stage 3: Distroless static (no shell, no package manager, ~2MB)
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /demo /demo

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/demo"]
