#######################################
# Build the app binary
#######################################

FROM golang:1.20-alpine3.17 AS builder

WORKDIR /app

RUN update-ca-certificates

# copy insides of ./indexer/ to /app
COPY ./indexer/ .
# name of the built binary is "app"
RUN go mod tidy && \
    go build -o app

#######################################
# Create the production image
#######################################

FROM alpine:3.17

WORKDIR /app

# copy binary
COPY --from=builder /app/app .
# copy certs from the builder image
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY ./docker/local/indexer/entrypoint /entrypoint
COPY ./docker/local/indexer/start-indexer /start-indexer
COPY ./docker/local/indexer/start-api /start-api

RUN chmod +x /entrypoint && \
    chmod +x /start-indexer && \
    chmod +x /start-api

# curl is used in the entrypoint script
RUN apk add --no-cache curl

ENTRYPOINT ["/entrypoint"]
