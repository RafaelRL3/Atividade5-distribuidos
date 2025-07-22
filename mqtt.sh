#!/usr/bin/env bash

ITERATIONS=30
MESSAGES=200000
BROKER="tcp://localhost:1883"
TOPIC="bench/topic"

for i in $(seq 1 "$ITERATIONS"); do
  echo "--- Iteration #$i  $MESSAGES messages ---"

  # Inicia o consumidor em segundo plano
  go run mqtt/consumer_mqtt.go -n "$MESSAGES" -broker "$BROKER" -topic "$TOPIC" &
  CONSUMER_PID=$!

  sleep 1   # dá tempo para ele se inscrever

  # Publica as mensagens
  go run mqtt/producer_mqtt.go -n "$MESSAGES" -broker "$BROKER" -topic "$TOPIC"

  # Aguarda o consumer imprimir a latência média e sair
  wait "$CONSUMER_PID"
done
