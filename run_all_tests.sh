#!/bin/bash

# Messaging Systems Benchmark - Master Test Runner
# Runs all messaging system benchmarks and analyzes results

echo "=========================================="
echo "  Messaging Systems Benchmark Suite"
echo "=========================================="
echo "Business Contract: Timestamp messages (~19 bytes), at-least-once delivery"
echo "Testing: Simplified TCP Queue, RabbitMQ, MQTT, Kafka"
echo ""

# Create all results directories
mkdir -p results/{simplified,rabbitmq,mqtt,kafka}

# Function to check if a service is running
check_service() {
    local service=$1
    local host=$2
    local port=$3
    
    if nc -z "$host" "$port" 2>/dev/null; then
        echo "✓ $service is running on $host:$port"
        return 0
    else
        echo "✗ $service is NOT running on $host:$port"
        return 1
    fi
}

# Check service availability
echo "Checking service availability..."
SIMPLIFIED_AVAILABLE=true
RABBITMQ_AVAILABLE=$(check_service "RabbitMQ" "localhost" "5672" && echo true || echo false)
MQTT_AVAILABLE=$(check_service "MQTT Broker" "localhost" "1883" && echo true || echo false)
KAFKA_AVAILABLE=$(check_service "Kafka" "localhost" "9092" && echo true || echo false)

echo ""
echo "Services to test:"
echo "- Simplified TCP Queue: Always available (built-in server)"
[ "$RABBITMQ_AVAILABLE" = true ] && echo "- RabbitMQ: Available" || echo "- RabbitMQ: Skipping (not running)"
[ "$MQTT_AVAILABLE" = true ] && echo "- MQTT: Available" || echo "- MQTT: Skipping (not running)"
[ "$KAFKA_AVAILABLE" = true ] && echo "- Kafka: Available" || echo "- Kafka: Skipping (not running)"

echo ""
read -p "Press Enter to start tests, or Ctrl+C to cancel..."

# Run Simplified TCP Queue test
echo ""
echo "=========================================="
echo "Running Simplified TCP Queue benchmark..."
echo "=========================================="
chmod +x simplified.sh
./simplified.sh

# Run RabbitMQ test if available
if [ "$RABBITMQ_AVAILABLE" = true ]; then
    echo ""
    echo "=========================================="
    echo "Running RabbitMQ benchmark..."
    echo "=========================================="
    chmod +x rabbitmq.sh
    ./rabbitmq.sh
fi

# Run MQTT test if available
if [ "$MQTT_AVAILABLE" = true ]; then
    echo ""
    echo "=========================================="
    echo "Running MQTT benchmark..."
    echo "=========================================="
    chmod +x mqtt.sh
    ./mqtt.sh
fi

# Run Kafka test if available
if [ "$KAFKA_AVAILABLE" = true ]; then
    echo ""
    echo "=========================================="
    echo "Running Kafka benchmark..."
    echo "=========================================="
    chmod +x kafka.sh
    ./kafka.sh
fi

# Analyze all results
echo ""
echo "=========================================="
echo "Analyzing Results..."
echo "=========================================="
go run analyze_results.go

echo ""
echo "=========================================="
echo "Benchmark Complete!"
echo "=========================================="
echo "Individual results saved in results/ directories"
echo "Re-run analysis anytime with: go run analyze_results.go"