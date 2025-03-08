#Etapa 1: Construcción

FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . . 

RUN go build -o faasAPI.exe .

FROM alpine:latest 

WORKDIR /app

COPY --from=builder /app/faasAPI.exe /app
COPY --from=builder /app/.env /app

EXPOSE 8080

CMD ["./faasAPI.exe"]