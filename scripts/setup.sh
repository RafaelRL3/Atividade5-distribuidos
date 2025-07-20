#!/bin/bash

# Environment Setup Script for Messaging Benchmark
# This script sets up the development environment and dependencies

echo "Setting up Messaging Systems Benchmark Environment..."

# Create project structure
echo "Creating project directories..."
mkdir -p simplified rabbitmq kafka mqtt scripts

# Initialize Go module
echo "Initializing Go module..."
go mod init messaging-benchmark

# Install Go dependencies
echo "Installing Go dependencies..."
go get github.com/eclipse/paho.mqtt.golang
go get github.com/segmentio/kafka-go
go get github.com/streadway/amqp

# Make scripts executable
echo "Setting script permissions..."
chmod +x scripts/*.sh 2>/dev/null || true

echo ""
echo "=== Infrastructure Setup Instructions ==="
echo ""
echo "1. RabbitMQ Setup:"
echo "   Docker: docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management"
echo "   Ubuntu: sudo apt-get install rabbitmq-server"
echo "   macOS: brew install rabbitmq"
echo ""
echo "2. Apache Kafka Setup:"
echo "   Download from: https://kafka.apache.org/downloads"
echo "   Or use Docker Compose (recommended)"
echo ""
echo "3. MQTT Broker (Mosquitto) Setup:"
echo "   Docker: docker run -it -p 1883:1883 eclipse-mosquitto"
echo "   Ubuntu: sudo apt-get install mosquitto"
echo "   macOS: brew install mosquitto"
echo ""
echo "Environment setup complete!"
echo "Please ensure the message brokers are running before executing benchmarks."