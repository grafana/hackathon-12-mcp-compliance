#!/bin/bash

# Build the server
echo "Building FedRAMP Compliance Assistant..."
cd ..
go build -o fedramp-compliance-assistant cmd/server/main.go
cd test

# Run the client test
echo "Running client test..."
go test -v

# Run the client program
echo "Running client program..."
go run client.go 