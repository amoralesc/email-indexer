#######################################
# Build stage
#######################################

FROM node:18.15-alpine3.17 AS builder

WORKDIR /app

COPY ./emails-app/package*.json .
RUN npm install

COPY ./emails-app/ .
RUN npm run build

#######################################
# Production stage
#######################################

FROM nginx:stable-alpine3.17

RUN mkdir /app

COPY --from=builder /app/dist /app
COPY ./docker/web/nginx.conf /etc/nginx/nginx.conf
