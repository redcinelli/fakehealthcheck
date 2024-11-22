# FROM debian:bullseye-slim AS up_downloader

# RUN apt-get update && apt-get install -y wget tar ca-certificates

# WORKDIR /pf
# RUN wget -O pf-host-agent.tgz "https://artifacts.elastic.co/downloads/prodfiler/pf-host-agent-8.16.0-linux-x86_64.tar.gz" && tar xzf pf-host-agent.tgz


# Use the official Golang image as a base
FROM golang:1.20 AS builder

# Set the working directory in the container
WORKDIR /app

# Copy the Go modules manifests and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# Use a lightweight image for the final stage
FROM gcr.io/distroless/static-debian11

# Set the working directory in the container
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app/app .
# COPY --from=up_downloader /pf/pf-host-agent-8.16.0-linux-x86_64/pf-host-agent .

# Expose port 8080 for the application
EXPOSE 8080

# RUN chmod +x pf-host-agent

# Run the application
# CMD ["./pf-host-agent", "-project-id=1", "-secret-token=${ES_token}", "-collection-agent=0e9b454bad67432baf59923546ced6a2.profiling.us-central1.gcp.cloud.es.io:443", "&&", "./app" ]
CMD ["./app"]