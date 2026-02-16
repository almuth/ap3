#!/bin/bash

# AhadPOS Go API - Docker Deployment Script

set -e

echo "=================================="
echo "AhadPOS Go API - Docker Setup"
echo "=================================="
echo

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "Error: Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Check if ahadpos3.db exists in parent directory
if [ ! -f "../ahadpos3.db" ]; then
    echo "Warning: ahadpos3.db not found in parent directory."
    echo "Please ensure the database file exists or create an empty one."
    echo
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Create necessary directories
echo "Creating directories..."
mkdir -p logs
echo "Done."
echo

# Build and start services
echo "Building and starting AhadPOS API..."
if command -v docker-compose &> /dev/null; then
    docker-compose up -d --build
else
    docker compose up -d --build
fi

echo
echo "=================================="
echo "AhadPOS API is now running!"
echo "=================================="
echo
echo "API URL: http://localhost:8080"
echo "Health Check: http://localhost:8080/health"
echo
echo "Useful commands:"
echo "  View logs: docker logs -f ahadpos-api"
echo "  Stop services: docker-compose down (or: docker compose down)"
echo "  Restart services: docker-compose restart (or: docker compose restart)"
echo
echo "Default credentials:"
echo "  Username: admin"
echo "  Password: admin123"
echo
