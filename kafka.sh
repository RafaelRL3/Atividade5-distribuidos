# !/usr/bin/env bash
set -euo pipefail

BROKERS="localhost:9092"
TOPIC="bench_topic"
MSGS=10000
ROUNDS=30

# cria tópico se não existir (idempotente)
docker exec kafka kafka-topics.sh \
  --bootstrap-server "$BROKERS" \
  --create --if-not-exists \
  --topic "$TOPIC" \
  --partitions 1 \
  --replication-factor 1 2>/dev/null || true

echo "Kafka benchmark — $ROUNDS runs × $MSGS msgs"

for ((i=1; i<=ROUNDS; i++)); do
  echo "---- Iteration #$i ----"

  GROUP="bench-$i-$(date +%s%N)"   # GroupID único por rod.

  # consumidor em background
  go run kafka/consumer_kafka.go -brokers "$BROKERS" -topic "$TOPIC" -n "$MSGS" -group "$GROUP" &
  CPID=$!

  sleep 0.5   # mínima margem para conectar

  # produtor envia lote
  go run kafka/producer_kafka.go -brokers "$BROKERS" -topic "$TOPIC" -n "$MSGS"

  # espera consumidor terminar
  wait "$CPID"
done

