#!/bin/bash

# Messaging Benchmark - MQTT
# Business Contract: Timestamp messages (19 bytes), QoS 1 (at-least-once delivery)

echo "Starting MQTT benchmark..."
echo "Business Contract: Timestamp messages (~19 bytes), QoS 1 (at-least-once delivery)"

# Check if MQTT broker is running
if ! nc -z localhost 1883; then
    echo "Error: MQTT broker is not running on localhost:1883"
    echo "Please start MQTT broker (e.g., Mosquitto) before running this test"
    exit 1
fi

# Create results directory
mkdir -p results/mqtt

echo "Running 30 test iterations with 10,000 messages each..."

for i in {1..30}
do
   echo "--- Iteration #$i: 10,000 messages ---"

   # Start consumer in background with specific output file
   OUTPUT_FILE="results/mqtt/test_$(date +%s)_${i}.txt"
   go run mqtt/consumer_mqtt.go -n 10000 -output "$OUTPUT_FILE" &
   CONSUMER_PID=$!

   sleep 2  # Give consumer time to connect and subscribe

   # Send messages
   go run mqtt/producer_mqtt.go -n 10000

   # Wait for consumer to finish
   wait $CONSUMER_PID
   echo "Iteration #$i completed"
done

echo "MQTT benchmark completed!"
echo "Results saved in results/mqtt/"
echo "Run 'go run analyze_results.go' to analyze all results"