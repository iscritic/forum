# # FROM golang:1.20-alpine
# FROM ubuntu:latest

# # RUN apk add --no-cache gcc musl-dev
# # RUN apk add --no-cache curl
# RUN apt update && \
#     apt install golang sqlite3 -y

# WORKDIR /app

# # COPY go.mod go.sum ./

# COPY . .

# # RUN go mod download
# # RUN go mod tidy

# ENV ENVIRONMENT=local

# # Build the Go application
# # RUN CGO_ENABLED=1 GOOS=linux go build -o forum ./cmd/
# # RUN go build -o forum ./cmd/

# # Expose the port your application will run on
# EXPOSE 8080

# # Copy the local JSON configuration file
# # COPY config/local.json /app/config/local.json

# # Command to run the Go application
# CMD ["bash", "-c", "go run ./cmd/"]

FROM golang:1.20-alpine

RUN apk add --no-cache gcc musl-dev
RUN apk add --no-cache curl

WORKDIR /app

COPY . .
# COPY go.mod go.sum ./

RUN go mod download

# RUN go build -o forum ./cmd/
RUN CGO_ENABLED=1 GOOS=linux go build -o forum ./cmd/

EXPOSE 8080

CMD ["./forum"]