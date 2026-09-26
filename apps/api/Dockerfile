# syntax=docker/dockerfile:1.7

FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="-s -w" -o /out/api ./cmd

# Keep only the certificates, user database, and health probe needed by the
# scratch runtime. This avoids shipping Alpine's shell and package manager.
FROM alpine:3.22 AS runtime-assets
RUN apk add --no-cache ca-certificates \
    && addgroup -S aisha \
    && adduser -S -G aisha aisha

FROM scratch
WORKDIR /app
COPY --from=runtime-assets /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=runtime-assets /etc/passwd /etc/passwd
COPY --from=runtime-assets /etc/group /etc/group
COPY --from=runtime-assets /lib/ld-musl-x86_64.so.1 /lib/ld-musl-x86_64.so.1
COPY --from=runtime-assets /bin/busybox /bin/busybox
COPY --from=build /out/api /app/api
COPY --from=build /src/openapi /app/openapi
USER aisha
EXPOSE 8088
ENTRYPOINT ["./api"]
