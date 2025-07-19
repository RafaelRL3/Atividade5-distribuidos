#!/bin/bash

for i in {1..30}
do
   echo "--- Iteration #$i 10000 messages---"

   # Start consumer in background and give it time to be ready
   go run rabbitmq/consumer_rmq.go -n 10000 &
   CONSUMER_PID=$!

   sleep 1  # Dá tempo para o consumer se conectar e começar a consumir

   # Agora sim envia as mensagens
   go run rabbitmq/producer_rmq.go -n 10000

   # Espera o consumer terminar
   wait $CONSUMER_PID
done
