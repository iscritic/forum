FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o forum ./cmd/

RUN chmod +x ./forum

EXPOSE 4269

CMD ["./forum"]
