FROM golang:1.18-alpine

# installing gcc compiler as needed for the sql driver, then sqlite
RUN apk add build-base
RUN apk add sqlite-dev

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go test ./...

RUN go build -o server .

CMD [ "/app/server" ]