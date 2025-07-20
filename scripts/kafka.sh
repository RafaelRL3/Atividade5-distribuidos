#!/bin/bash

# Kafka Benchmark Script
# Runs 30 iterations of 10,000 messages each

for i in {1..30}
do
   echo "--- Iteration #$i: 10,000 messages ---"

   # Start consumer in background and give it time to be ready
   go run kafka/consumer_kafka.go -n 10000 &
   CONSUMER_PID=$!

   sleep 1  # Give time for consumer to connect and start consuming

   # Send messages
   go run kafka/producer_kafka.go -n 10000

   # Wait for consumer to finish
   wait $CONSUMER_PID
done