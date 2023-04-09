#######################################
# Build the indexer binary
#######################################

FROM golang:1.20-alpine3.17 AS builder

WORKDIR /app
# copy insides of ./indexer/ to /app
COPY ./indexer/ .
# name of the built binary is "indexer"
RUN go mod tidy && \
    go build -o indexer 

#######################################
# Build the indexer image
#######################################

FROM alpine:3.17

WORKDIR /app

COPY --from=builder /app/indexer .

COPY ./docker/indexer/entrypoint /entrypoint
COPY ./docker/indexer/start-indexer /start-indexer
COPY ./docker/indexer/start-server /start-server

RUN chmod +x /entrypoint && \
    chmod +x /start-indexer && \
    chmod +x /start-server

RUN apk add --no-cache curl

ENTRYPOINT ["/entrypoint"]
