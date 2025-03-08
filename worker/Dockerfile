#Etapa 1: Construcción

FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . . 

RUN go build -o worker .

FROM alpine:latest 

WORKDIR /app

RUN apk add --no-cache docker-cli

COPY --from=builder /app/worker /app

EXPOSE 8080

CMD ["./worker"]