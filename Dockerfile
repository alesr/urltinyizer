FROM alpine:latest

RUN apk add --update go

WORKDIR /app

COPY urltinyizer .

EXPOSE 8080

ENTRYPOINT ["./urltinyizer"]
