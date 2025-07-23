#!/usr/bin/env bash
#
# Benchmark Kafka-go – stack sobe e desce a cada iteração
set -euo pipefail

# ===================== CONFIG =====================
TOPIC="bench_topic"
BROKER_EXT="localhost:9094"
N_MSG=200000
ITER=30
# ==================================================

run_quiet() { "$@" >/dev/null 2>&1; }           # executa sem stdout/stderr
wait_for_port() { local hp=$1; until nc -z "${hp%:*}" "${hp#*:}" 2>/dev/null; do sleep 1; done; }

for i in $(seq 1 "$ITER"); do
  run_quiet docker compose up -d
  run_quiet wait_for_port "$BROKER_EXT"

  # cria tópico (ignora se já existir)
  run_quiet docker compose exec kafka \
      kafka-topics.sh --create \
      --topic "$TOPIC" \
      --bootstrap-server kafka:9092 \
      --partitions 1 --replication-factor 1 || true

  TMP=$(mktemp)         # captura média impressa pelo consumidor

  go run kafka/consumer_kafka.go \
      -n "$N_MSG" -broker "$BROKER_EXT" -topic "$TOPIC" >"$TMP" &
  CONS_PID=$!

  sleep 1
  run_quiet go run kafka/producer_kafka.go \
      -n "$N_MSG" -broker "$BROKER_EXT" -topic "$TOPIC"

  wait "$CONS_PID"

  avg_us=$(tail -n1 "$TMP" | tr -d '\r\n')
  rm -f "$TMP"

  printf 'Iter %2d: %s µs\n' "$i" "$avg_us"   # ← unidade incluída

  run_quiet docker compose down --volumes --remove-orphans
done