#syntax=docker/dockerfile:1.7

# Stage 1: Compile Tailwind CSS (always fresh, never stale)
# Scans all .templ files for class names and generates the final CSS.
# This eliminates the staleness risk of committing a static app.css.
FROM node:22-slim AS css
WORKDIR /src
COPY . .
# node:22-slim no longer ships pnpm on PATH; install the pinned version
# (matches website/package.json packageManager) via npm.
# pnpm 11 fails installs whose dependency build scripts were not explicitly
# approved (ERR_PNPM_IGNORED_BUILDS). @parcel/watcher, a transitive dep of
# @tailwindcss/cli, ships prebuilt binaries, so its install script is safe
# to approve via a pnpm-workspace.yaml allowBuilds entry.
RUN npm install -g pnpm@11.20.0 \
    && pnpm init -y \
    && printf 'allowBuilds:\n  "@parcel/watcher": true\n' > pnpm-workspace.yaml \
    && pnpm add tailwindcss @tailwindcss/cli \
    && pnpm dlx @tailwindcss/cli \
        -i examples/demo/demo.css \
        -o examples/demo/static/app.css \
        --minify

# Stage 2: Build templ + Go binary
FROM golang:1.26 AS builder

ENV GOEXPERIMENT=jsonv2
ENV CGO_ENABLED=0

WORKDIR /src

# The root go.mod replaces all 6 sub-modules with local paths, so
# `go mod download` needs their go.mod files present to resolve the graph.
COPY go.mod go.sum ./
COPY utils/go.mod utils/
COPY icons/go.mod icons/
COPY errorpage/go.mod errorpage/
COPY charts/echarts/go.mod charts/echarts/
COPY htmx/go.mod htmx/
COPY datastar/go.mod datastar/
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
