FROM golang:1.17-alpine as builder
WORKDIR /app
copy . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o jellyfin-proxy .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/jellyfin-proxy /usr/local/bin/jellyfin-proxy

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/jellyfin-proxy"]
