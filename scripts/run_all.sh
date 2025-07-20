#!/bin/bash

# Run All Benchmarks Script
# Executes all messaging system benchmarks sequentially

echo "=== Messaging Systems Benchmark Suite ==="
echo "Testing message latency across different systems"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    exit 1
fi

# Function to check if a port is open
check_port() {
    local port=$1
    local service=$2
    if nc -z localhost $port 2>/dev/null; then
        echo "✓ $service is running on port $port"
        return 0
    else
        echo "✗ $service is not running on port $port"
        return 1
    fi
}

echo "Checking infrastructure availability..."

# Check services
CUSTOM_SERVER_READY=false
RABBITMQ_READY=false
KAFKA_READY=false
MQTT_READY=false

if check_port 9000 "Custom Queue Server"; then
    CUSTOM_SERVER_READY=true
fi

if check_port 5672 "RabbitMQ"; then
    RABBITMQ_READY=true
fi

if check_port 9092 "Kafka"; then
    KAFKA_READY=true
fi

if check_port 1883 "MQTT Broker"; then
    MQTT_READY=true
fi

echo ""

# Run benchmarks based on available services
if [ "$CUSTOM_SERVER_READY" = true ]; then
    echo "=== Running Custom Queue Server Benchmark ==="
    echo "Messages: 200,000 per iteration × 30 iterations = 6,000,000 total"
    ./scripts/simplified.sh > results_custom_queue.txt 2>&1 &
    CUSTOM_PID=$!
    echo "Custom queue benchmark started (PID: $CUSTOM_PID)"
else
    echo "Skipping Custom Queue Server benchmark (server not running)"
fi

if [ "$RABBITMQ_READY" = true ]; then
    echo ""
    echo "=== Running RabbitMQ Benchmark ==="
    echo "Messages: 10,000 per iteration × 30 iterations = 300,000 total"
    ./scripts/rabbitmq.sh > results_rabbitmq.txt 2>&1 &
    RABBITMQ_PID=$!
    echo "RabbitMQ benchmark started (PID: $RABBITMQ_PID)"
else
    echo "Skipping RabbitMQ benchmark (service not available)"
fi

if [ "$KAFKA_READY" = true ]; then
    echo ""
    echo "=== Running Kafka Benchmark ==="
    echo "Messages: 10,000 per iteration × 30 iterations = 300,000 total"
    ./scripts/kafka.sh > results_kafka.txt 2>&1 &
    KAFKA_PID=$!
    echo "Kafka benchmark started (PID: $KAFKA_PID)"
else
    echo "Skipping Kafka benchmark (service not available)"
fi

if [ "$MQTT_READY" = true ]; then
    echo ""
    echo "=== Running MQTT Benchmark ==="
    echo "Messages: 5,000 per iteration × 30 iterations = 150,000 total"
    ./scripts/mqtt.sh > results_mqtt.txt 2>&1 &
    MQTT_PID=$!
    echo "MQTT benchmark started (PID: $MQTT_PID)"
else
    echo "Skipping MQTT benchmark (service not available)"
fi

echo ""
echo "Benchmarks are running in parallel..."
echo "Results will be saved to results_*.txt files"
echo ""
echo "To monitor progress:"
echo "  tail -f results_custom_queue.txt"
echo "  tail -f results_rabbitmq.txt" 
echo "  tail -f results_kafka.txt"
echo "  tail -f results_mqtt.txt"
echo ""
echo "Use 'ps aux | grep go' to see running processes"
echo "Use './scripts/analyze_results.sh' to analyze results when complete"