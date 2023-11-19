#First stage
FROM golang:1.20-alpine as builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o main ./cmd


#Second stage
FROM builder as test
RUN go test ./... -v


#Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 50051

CMD ["./main"]
