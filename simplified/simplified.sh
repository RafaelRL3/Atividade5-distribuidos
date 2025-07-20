#!/bin/bash

# Messaging Benchmark - Simplified TCP Queue
# Business Contract: Timestamp messages (19 bytes), at-least-once delivery

echo "Starting Simplified TCP Queue benchmark..."
echo "Business Contract: Timestamp messages (~19 bytes), at-least-once delivery"

# Create results directory
mkdir -p results/simplified

# Start server in background
echo "Starting queue server..."
go run simplified/server.go &
SERVER_PID=$!
sleep 2  # Give server time to start

echo "Running 30 test iterations with 10,000 messages each..."

for i in {1..30}
do
   echo "--- Iteration #$i: 10,000 messages ---"

   # Start consumer in background with specific output file
   OUTPUT_FILE="results/simplified/test_$(date +%s)_${i}.txt"
   go run simplified/consumer.go -n 10000 -output "$OUTPUT_FILE" &
   CONSUMER_PID=$!

   sleep 1  # Give consumer time to connect

   # Send messages
   go run simplified/producer.go -n 10000

   # Wait for consumer to finish
   wait $CONSUMER_PID
   echo "Iteration #$i completed"
done

# Terminate server
echo "Stopping queue server..."
kill $SERVER_PID

echo "Simplified TCP Queue benchmark completed!"
echo "Results saved in results/simplified/"
echo "Run 'go run analyze_results.go' to analyze all results"