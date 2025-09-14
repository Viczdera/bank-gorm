#BUILD STAGE
FROM golang:1.24-alpine3.21 AS builder
WORKDIR /app
COPY . .
#build
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | gunzip | tar xvf -

#RUN STAGE
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
#from the app
COPY app.env .
COPY wait-for.sh .
COPY start.sh .
COPY db/migrations ./migrations

#publish port 
EXPOSE 9090

# CMD ["/app/main"]
# ENTRYPOINT [ "/app/start.sh" ]
