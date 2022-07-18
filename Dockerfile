FROM golang:1.18

WORKDIR /usr/src/app
COPY . .
RUN go mod download && \
    go build -v -o /usr/local/bin/app main.go

EXPOSE 8080

CMD ["app"]