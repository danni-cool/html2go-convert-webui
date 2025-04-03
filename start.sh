#!/bin/bash

# Make sure the static directory exists
if [ ! -d "static" ]; then
  echo "Creating static directory from public..."
  mkdir -p static
  cp -r public/* static/
fi

# Run the server with automatic port selection
go run cmd/server/main.go "$@"
