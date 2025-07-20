#!/bin/bash

# MQTT Benchmark Script
# Runs 30 iterations of 5,000 messages each

for i in {1..30}
do
   echo "--- Iteration #$i: 5,000 messages ---"

   # Start consumer in background and give it time to be ready
   go run mqtt/consumer_mqtt.go -n 5000 &
   CONSUMER_PID=$!

   sleep 1  # Give time for consumer to connect and subscribe

   # Send messages
   go run mqtt/producer_mqtt.go -n 5000

   # Wait for consumer to finish
   wait $CONSUMER_PID
done