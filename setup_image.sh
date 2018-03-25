#!/bin/bash

CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo --ldflags="-s" -o transactionserver

# Build the image
docker build -t rattransaction .

# Remove remnants
rm -f transactionserver