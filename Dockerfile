FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV ENVIRONMENT=local

# Build the Go application
RUN go build -o forum ./cmd/

# Expose the port your application will run on
EXPOSE 4269

# Copy the local JSON configuration file
COPY config/local.json /app/config/local.json

# Command to run the Go application
CMD ["./forum"]
