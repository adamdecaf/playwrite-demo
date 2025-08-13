# Stage 1: Build the Go binary
FROM golang:1.25 AS go-builder

WORKDIR /app
COPY main.go .
COPY go.mod .
COPY go.sum .
RUN go build -o playwright-client main.go

# Stage 2: Set up Playwright server and copy Go binary
FROM mcr.microsoft.com/playwright:v1.54.2-jammy

WORKDIR /app

# Install Node.js dependencies
COPY package.json ./
RUN npm install

# Copy server.js and the Go binary
COPY server.js .
COPY --from=go-builder /app/playwright-client .

# Expose the WebSocket endpoint
EXPOSE 3000

# Start the Playwright server and run the Go client
CMD ["sh", "-c", "node server.js & sleep 5 && ./playwright-client"]
