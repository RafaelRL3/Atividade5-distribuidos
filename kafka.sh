#!/usr/bin/env bash
set -euo pipefail

MSGS=10
ROUNDS=30
TOPIC=bench_topic

# cria o tópico caso não exista
docker exec -it kafka kafka-topics.sh \
  --create --topic ${TOPIC} \
  --bootstrap-server localhost:9092 \
  --partitions 1 --replication-factor 1 2>/dev/null || true

for i in $(seq 1 $ROUNDS); do
  echo "---- Iteration #$i ($MSGS msgs) ----"

  GROUP="bench-$i"                             # grupo único p/ cada rodada
  go run kafka/consumer_kafka.go -n $MSGS -group $GROUP &  # consumidor
  PID=$!

  sleep 1
  go run kafka/producer_kafka.go -n $MSGS                 # produtor

  wait $PID
done
