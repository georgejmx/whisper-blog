FROM node:20-alpine AS client-builder

WORKDIR /tmp
COPY client/ /tmp/client
COPY routes/templates/ /tmp/routes/templates
WORKDIR /tmp/client
ENV NODE_ENV=development
RUN npm install
RUN npm run lint
RUN npm run bundle

FROM golang:1.18-alpine AS server-builder
RUN apk add build-base

WORKDIR /app
COPY go.* /app
COPY config/ /app/config
COPY types/ /app/types
COPY security/ /app/security
COPY routes/ /app/routes
COPY controller/ /app/controller
COPY utils/ /app/utils
COPY words/ /app/words
COPY main.go /app/main.go
COPY --from=client-builder /tmp/client/public /app/client/public
RUN go build -o server .

FROM alpine:latest
RUN apk add build-base && apk add sqlite-dev

WORKDIR /root
COPY --from=server-builder /app/server /root
RUN mkdir -p data
CMD [ "/root/server" ]