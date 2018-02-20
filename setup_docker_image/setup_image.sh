#!/bin/bash

cd ..
CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o transactionserver

# Build the image
docker build -t transactionserver .

# Remove remnants
rm -f transactionserver