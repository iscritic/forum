FROM golang:1.20-alpine

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o forum ./cmd/

EXPOSE 8080

CMD ["./forum"]