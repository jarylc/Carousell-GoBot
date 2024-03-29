FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -ldflags="-w -s"

FROM alpine
ENV UID=1000 \
    GID=1000
RUN apk add --no-cache su-exec tzdata
COPY --from=builder /app/carousell-gobot /app/carousell-gobot
COPY entrypoint.sh /app/entrypoint.sh
ENTRYPOINT [ "/app/entrypoint.sh" ]
CMD ["/app/carousell-gobot -c /data/config.yaml -s /data/state.json"]
