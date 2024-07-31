FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /go-message-service ./cmd

FROM gcr.io/distroless/base-debian10

COPY --from=builder /go-message-service /go-message-service

EXPOSE 8080

ENTRYPOINT ["/go-message-service"]