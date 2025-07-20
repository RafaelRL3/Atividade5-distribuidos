# Messaging Systems Benchmark

This project benchmarks latency performance across different messaging systems using a consistent business contract.

## Business Contract

All messaging systems implement the same contract:
- **Message Format**: Timestamp as string (nanoseconds since Unix epoch)
- **Message Size**: ~19 bytes (timestamp: 19 digits as string)
- **QoS Level**: At-least-once delivery guarantee
- **Measurement**: End-to-end latency from producer send to consumer receive
- **Output**: Individual latency measurements written to files for post-processing

## Project Structure

```
├── simplified/          # Custom TCP-based queue
│   ├── server.go        # Queue server
│   ├── producer.go      # Message producer
│   ├── consumer.go      # Message consumer
│   └── simplified.sh    # Test runner
├── rabbitmq/           # RabbitMQ implementation
│   ├── producer_rmq.go  # RabbitMQ producer
│   ├── consumer_rmq.go  # RabbitMQ consumer
│   └── rabbitmq.sh     # Test runner
├── mqtt/               # MQTT implementation
│   ├── producer_mqtt.go # MQTT producer
│   ├── consumer_mqtt.go # MQTT consumer
│   └── mqtt.sh         # Test runner
├── kafka/              # Kafka implementation
│   ├── producer_kafka.go # Kafka producer
│   ├── consumer_kafka.go # Kafka consumer
│   └── kafka.sh        # Test runner
├── results/            # Latency measurement files
│   ├── simplified/     # Simplified queue results
│   ├── rabbitmq/       # RabbitMQ results
│   ├── mqtt/           # MQTT results
│   └── kafka/          # Kafka results
├── analyze_results.go  # Post-processing analysis tool
├── run_all_tests.sh    # Master test runner (runs all available systems)
└── README.md           # This file
```

## Prerequisites

### For All Tests
- Go 1.19+
- Git (to clone dependencies)

### For Specific Systems
- **RabbitMQ**: RabbitMQ server running on localhost:5672
- **MQTT**: MQTT broker (e.g., Mosquitto) running on localhost:1883
- **Kafka**: Kafka cluster running on localhost:9092

## Installation

1. Clone the repository and navigate to the project directory

2. Install Go dependencies:
```bash
go mod init messaging-benchmark
go get github.com/eclipse/paho.mqtt.golang
go get github.com/segmentio/kafka-go
go get github.com/streadway/amqp
```

3. Create results directories:
```bash
mkdir -p results/{simplified,rabbitmq,mqtt,kafka}
```

## Running Tests

### Quick Start - Run All Tests
```bash
chmod +x run_all_tests.sh
./run_all_tests.sh
```
This script will automatically detect available services and run all applicable tests.

### Individual System Tests

#### Simplified TCP Queue
No external dependencies required.
```bash
chmod +x simplified.sh
./simplified.sh
```

#### RabbitMQ
Ensure RabbitMQ is running:
```bash
# Start RabbitMQ (example for Ubuntu/Debian)
sudo systemctl start rabbitmq-server

# Run tests
chmod +x rabbitmq.sh
./rabbitmq.sh
```

#### MQTT
Ensure MQTT broker is running:
```bash
# Start Mosquitto (example)
mosquitto -c /etc/mosquitto/mosquitto.conf -d

# Run tests
chmod +x mqtt.sh
./mqtt.sh
```

#### Kafka
Ensure Kafka is running:
```bash
# Start Kafka (example)
bin/kafka-server-start.sh config/server.properties

# Create topic (optional - will auto-create if not exists)
bin/kafka-topics.sh --create --topic bench_topic --bootstrap-server localhost:9092

# Run tests
chmod +x kafka.sh
./kafka.sh
```

## Analyzing Results

After running all tests, analyze the results:
```bash
go run analyze_results.go
```

This will calculate and display average latencies for each system and test iteration.

## Message Details

- **Format**: Unix timestamp in nanoseconds as string
- **Example**: `1642771200123456789`
- **Size**: ~19 bytes (19 digits)
- **Encoding**: UTF-8 string

## Test Configuration

- **Default Messages per Test**: 10,000 (configurable via `-n` flag)
- **Test Iterations**: 30 per system
- **Measurement Unit**: Microseconds
- **Output Format**: One latency value per line in result files

## Customization

All programs accept command-line flags:
- `-n`: Number of messages (default varies by system)
- `-broker`/`-brokers`/`-addr`: Connection endpoints
- `-topic`: Topic/queue name (where applicable)

Example:
```bash
go run simplified/producer.go -n 5000 -addr localhost:9000
```