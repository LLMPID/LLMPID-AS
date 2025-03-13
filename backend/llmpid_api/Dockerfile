# ---- Build Stage ----
FROM golang:1.24-alpine  AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN mkdir bin/

RUN go build -o /src/bin/llmpid_api ./cmd/main.go

# ---- Run Stage ----
FROM alpine:latest

WORKDIR /api

COPY --from=build /src/bin/llmpid_api .
COPY --from=build /src/config/config.yaml ./config/config.yaml

EXPOSE 8080

CMD ["./llmpid_api"]