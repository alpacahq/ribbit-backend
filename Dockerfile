FROM golang:1.16-alpine
RUN apk update && apk add --virtual build-dependencies build-base gcc wget git openssl

WORKDIR /app

COPY ./ ./
RUN go mod download
RUN cp .env.sample .env
RUN chmod 100 ./initdb.sh
# RUN go build -o mvp ./entry/

EXPOSE 8080

CMD [ "go", "run", "./entry/main.go" ]