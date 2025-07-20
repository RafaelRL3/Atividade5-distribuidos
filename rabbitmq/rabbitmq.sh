#!/bin/bash

# Messaging Benchmark - RabbitMQ
# Business Contract: Timestamp messages (19 bytes), at-least-once delivery

echo "Starting RabbitMQ benchmark..."
echo "Business Contract: Timestamp messages (~19 bytes), at-least-once delivery"

# Check if RabbitMQ is running
if ! nc -z localhost 5672; then
    echo "Error: RabbitMQ server is not running on localhost:5672"
    echo "Please start RabbitMQ server before running this test"
    exit 1
fi

# Create results directory
mkdir -p results/rabbitmq

echo "Running 30 test iterations with 10,000 messages each..."

for i in {1..30}
do
   echo "--- Iteration #$i: 10,000 messages ---"

   # Start consumer in background with specific output file
   OUTPUT_FILE="results/rabbitmq/test_$(date +%s)_${i}.txt"
   go run rabbitmq/consumer_rmq.go -n 10000 -output "$OUTPUT_FILE" &
   CONSUMER_PID=$!

   sleep 2  # Give consumer time to connect and declare queue

   # Send messages
   go run rabbitmq/producer_rmq.go -n 10000

   # Wait for consumer to finish
   wait $CONSUMER_PID
   echo "Iteration #$i completed"
done

echo "RabbitMQ benchmark completed!"
echo "Results saved in results/rabbitmq/"
echo "Run 'go run analyze_results.go' to analyze all results"