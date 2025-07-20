#!/bin/bash

# Messaging Benchmark - Kafka
# Business Contract: Timestamp messages (19 bytes), at-least-once delivery

echo "Starting Kafka benchmark..."
echo "Business Contract: Timestamp messages (~19 bytes), at-least-once delivery"

# Check if Kafka is running
if ! nc -z localhost 9092; then
    echo "Error: Kafka broker is not running on localhost:9092"
    echo "Please start Kafka before running this test"
    exit 1
fi

# Create results directory
mkdir -p results/kafka

# Create topic if it doesn't exist (requires kafka-topics.sh in PATH)
echo "Ensuring topic 'bench_topic' exists..."
if command -v kafka-topics.sh &> /dev/null; then
    kafka-topics.sh --create --if-not-exists --topic bench_topic \
        --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1 2>/dev/null || true
else
    echo "Warning: kafka-topics.sh not found in PATH. Assuming topic exists or will be auto-created."
fi

echo "Running 30 test iterations with 10,000 messages each..."

for i in {1..30}
do
   echo "--- Iteration #$i: 10,000 messages ---"

   # Use unique consumer group for each test to avoid offset issues
   CONSUMER_GROUP="bench-consumer-group-$(date +%s)-${i}"
   
   # Start consumer in background with specific output file
   OUTPUT_FILE="results/kafka/test_$(date +%s)_${i}.txt"
   go run kafka/consumer_kafka.go -n 10000 -group "$CONSUMER_GROUP" -output "$OUTPUT_FILE" &
   CONSUMER_PID=$!

   sleep 3  # Give consumer time to connect and join group

   # Send messages
   go run kafka/producer_kafka.go -n 10000

   # Wait for consumer to finish
   wait $CONSUMER_PID
   echo "Iteration #$i completed"
done

echo "Kafka benchmark completed!"
echo "Results saved in results/kafka/"
echo "Run 'go run analyze_results.go' to analyze all results"