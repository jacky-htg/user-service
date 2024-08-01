FROM golang:alpine

WORKDIR /app
COPY user-service /app/
COPY ./cmd/cmd-user-service /app/

CMD ["/app/user-service"]
