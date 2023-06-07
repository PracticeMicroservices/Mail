FROM golang:1.18-alpine AS builder

RUN mkdir /app

COPY . /app
COPY templates /templates
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build  -o mail-service ./cmd/api
RUN chmod +x /app/mail-service

#build a tiny docker image
FROM  alpine:latest


RUN mkdir /app

COPY --from=builder /app/mail-service /app
COPY templates /templates

CMD ["/app/mail-service"]