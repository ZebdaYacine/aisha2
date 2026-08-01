FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api
FROM alpine:3.22
RUN addgroup -S aisha && adduser -S -G aisha aisha
WORKDIR /app
COPY --from=build /out/api ./api
COPY openapi ./openapi
USER aisha
EXPOSE 8080
ENTRYPOINT ["./api"]
