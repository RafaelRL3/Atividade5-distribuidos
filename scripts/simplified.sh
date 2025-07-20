#!/bin/bash

# Start server in background
go run simplified/server.go &
SERVER_PID=$!
sleep 1  # Dá tempo para o servidor iniciar

for i in {1..30}
do
   echo "--- Iteration #$i 200000 messages---"

   # Start consumer in background and give it time to be ready
   go run simplified/consumer.go -n 200000 &
   CONSUMER_PID=$!

   sleep 1  # Dá tempo para o consumer se conectar e começar a consumir

   # Agora sim envia as mensagens
   go run simplified/producer.go -n 200000

   # Espera o consumer terminar
   wait $CONSUMER_PID
done

# Termina o servidor
kill $SERVER_PID