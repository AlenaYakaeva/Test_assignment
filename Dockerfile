# build stage
FROM golang:1.25 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o example-app-service cmd/main.go
RUN ls /app

# runtime stage
FROM alpine:3.14

WORKDIR /root
COPY --from=builder /app/example-app-service .

EXPOSE 8080
CMD [ "./example-app-service" ]