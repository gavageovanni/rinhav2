FROM golang:latest

WORKDIR /app

COPY .. .

RUN go get -d -v ./...

RUN go build -o app

EXPOSE 8080

CMD ["./app"]

#FROM golang:1.21.3-alpine as builder
#RUN apk add --no-cache gcc musl-dev
#WORKDIR /app
#COPY go.mod go.sum ./
#RUN go mod download
#COPY . .
#RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server
#
#FROM alpine:latest
#WORKDIR /root/
#COPY --from=builder /app/main .
#EXPOSE 3001
#CMD ["./main"]