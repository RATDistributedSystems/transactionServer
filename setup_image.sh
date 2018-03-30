#!/bin/bash

# front end
git clone https://github.com/RATDistributedSystems/frontend

CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo --ldflags="-s" -o transactionServer

# Build the image
docker build -t rattransaction .

# Remove remnants
rm -rf transactionServer frontend