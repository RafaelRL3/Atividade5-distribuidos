# Messaging Systems Benchmark

This project benchmarks different messaging systems by measuring end-to-end latency across multiple implementations:

- **Custom Queue Server** (simplified implementation)
- **RabbitMQ** (AMQP broker)
- **Apache Kafka** (distributed streaming platform)
- **MQTT** (lightweight pub/sub protocol)

## Business Contract

All messaging systems implement the same business contract for fair comparison:

### Message Format
- **Payload**: Unix timestamp in nanoseconds as string
- **Purpose**: Measure end-to-end latency from producer to consumer
- **Size**: ~19-20 bytes (timestamp as string: "1690804800123456789")

### Producer Contract
- Generates current timestamp (`time.Now().UnixNano()`)
- Converts to string format
- Publishes to designated topic/queue
- Configurable message count via `-n` flag

### Consumer Contract
- Subscribes to designated topic/queue
- Receives timestamp payload
- Calculates latency: `current_time - received_timestamp`
- Collects latency statistics
- Outputs average latency in microseconds
- Exits after receiving specified number of messages (`-n` flag)

### Common Configuration
- Default message count: 1000 messages
- Topic/Queue names follow pattern: `bench_*` or similar
- Connection timeouts and error handling
- Graceful shutdown procedures

## Quick Start

### Option 1: Docker Infrastructure (Recommended)
```bash
# Clone or create project directory
mkdir messaging-benchmark && cd messaging-benchmark

# Setup environment and dependencies
./scripts/setup.sh

# Start all infrastructure services
docker-compose up -d

# Wait for services to be ready (check health)
docker-compose ps

# Run all benchmarks
./scripts/run_all.sh

# Analyze results
./scripts/analyze_results.sh
```

### Option 2: Manual Setup
```bash
# Setup Go environment
./scripts/setup.sh

# Start services manually (see Prerequisites section)
# Then run individual benchmarks or all at once
./scripts/run_all.sh
```

### Go Dependencies
```bash
go mod init messaging-benchmark
go get github.com/eclipse/paho.mqtt.golang
go get github.com/segmentio/kafka-go
go get github.com/streadway/amqp
```

#### Docker Infrastructure (Recommended)
```bash
# Start all services at once
docker-compose up -d

# Check service health
docker-compose ps

# View logs if needed
docker-compose logs rabbitmq
docker-compose logs kafka
docker-compose logs mosquitto

# Stop all services
docker-compose down
```

#### Individual Service Setup

#### RabbitMQ
```bash
# Using Docker
docker run -d --name rabbitmq \
  -p 5672:5672 -p 15672:15672 \
  rabbitmq:3-management

# Or install locally
# Ubuntu/Debian: apt-get install rabbitmq-server
# macOS: brew install rabbitmq
```

#### Apache Kafka
```bash
# Using Docker Compose (recommended)
# Create docker-compose.yml with Kafka + Zookeeper
docker-compose up -d

# Or download from Apache Kafka website
# Requires Java 8+
```

#### MQTT Broker (Mosquitto)
```bash
# Using Docker
docker run -it -p 1883:1883 eclipse-mosquitto

# Or install locally
# Ubuntu/Debian: apt-get install mosquitto
# macOS: brew install mosquitto
```

## Project Structure

```
messaging-benchmark/
├── simplified/
│   ├── server.go      # Custom queue server
│   ├── producer.go    # Custom queue producer
│   └── consumer.go    # Custom queue consumer
├── rabbitmq/
│   ├── producer_rmq.go # RabbitMQ producer
│   └── consumer_rmq.go # RabbitMQ consumer
├── kafka/
│   ├── producer_kafka.go # Kafka producer
│   └── consumer_kafka.go # Kafka consumer
├── mqtt/
│   ├── producer_mqtt.go  # MQTT producer
│   └── consumer_mqtt.go  # MQTT consumer
├── scripts/
│   ├── simplified.sh     # Benchmark script for custom queue
│   ├── rabbitmq.sh      # Benchmark script for RabbitMQ
│   ├── kafka.sh         # Benchmark script for Kafka
│   ├── mqtt.sh          # Benchmark script for MQTT
│   ├── setup.sh         # Environment setup script
│   ├── run_all.sh       # Run all benchmarks
│   └── analyze_results.sh # Results analysis script
├── docker-compose.yml    # Infrastructure setup
├── mosquitto.conf       # MQTT broker configuration
└── README.md
```

## Running the Benchmarks

### All Benchmarks at Once
```bash
# Run all available benchmarks automatically
./scripts/run_all.sh

# This will:
# 1. Check which services are running
# 2. Start appropriate benchmarks in parallel
# 3. Save results to separate files
# 4. Show progress monitoring commands
```

### Individual Benchmark Runs

### 1. Custom Queue Server Benchmark

The custom implementation uses a simple TCP server with FIFO queue.

**Message Flow:**
- Producer sends: `PUSH {timestamp}\n`
- Server responds: `OK\n`
- Consumer sends: `PULL\n`
- Server responds: `MSG {timestamp}\n` or `EMPTY\n`

```bash
# Make script executable
chmod +x scripts/simplified.sh

# Run benchmark (30 iterations of 200,000 messages each)
./scripts/simplified.sh
```

**Configuration:**
- Server listens on port 9000
- Protocol: Custom TCP-based
- Message size: ~19-20 bytes per message
- Total benchmark: 6,000,000 messages

### 2. RabbitMQ Benchmark

Uses AMQP protocol with direct queue messaging.

```bash
# Ensure RabbitMQ is running
sudo systemctl start rabbitmq-server
# or with Docker: docker start rabbitmq

# Make script executable
chmod +x scripts/rabbitmq.sh

# Run benchmark (30 iterations of 10,000 messages each)
./scripts/rabbitmq.sh
```

**Configuration:**
- Queue: `bench_queue`
- Exchange: Default (direct)
- Message size: ~19-20 bytes per message
- Total benchmark: 300,000 messages

### 3. Apache Kafka Benchmark

Uses Kafka's high-throughput streaming capabilities.

```bash
# Ensure Kafka is running
# With Docker: docker-compose up -d kafka
# Check: docker-compose ps kafka

# Create topic (optional - auto-created)
docker exec benchmark-kafka kafka-topics --create --topic bench_topic \
  --bootstrap-server localhost:9092 \
  --partitions 1 --replication-factor 1

# Run benchmark script (30 iterations of 10,000 messages)
./scripts/kafka.sh

# Or run manually
go run kafka/consumer_kafka.go -n 10000 &
sleep 1
go run kafka/producer_kafka.go -n 10000
```

**Configuration:**
- Topic: `bench_topic`
- Brokers: `localhost:9092`
- Message size: ~19-20 bytes per message
- Total benchmark: 300,000 messages

### 4. MQTT Benchmark

Uses lightweight pub/sub messaging for IoT scenarios.

```bash
# Ensure MQTT broker is running
# With Docker: docker-compose up -d mosquitto
# Check: docker-compose ps mosquitto

# Run benchmark script (30 iterations of 5,000 messages)
./scripts/mqtt.sh

# Or run manually
go run mqtt/consumer_mqtt.go -n 5000 &
sleep 1
go run mqtt/producer_mqtt.go -n 5000
```

**Configuration:**
- Topic: `bench/topic`
- Broker: `tcp://localhost:1883`
- QoS Level: 1 (at least once delivery)
- Message size: ~19-20 bytes per message
- Total benchmark: 150,000 messages

## Manual Testing

### Single Message Test
```bash
# Custom Queue
go run simplified/consumer.go -n 1 &
go run simplified/server.go &
sleep 1
go run simplified/producer.go -n 1

# RabbitMQ
go run rabbitmq/consumer_rmq.go -n 1 &
sleep 1
go run rabbitmq/producer_rmq.go -n 1

# Kafka
go run kafka/consumer_kafka.go -n 1 &
sleep 1
go run kafka/producer_kafka.go -n 1

# MQTT
go run mqtt/consumer_mqtt.go -n 1 &
sleep 1
go run mqtt/producer_mqtt.go -n 1
```

## Results Analysis

### Automated Analysis
```bash
# Analyze all benchmark results
./scripts/analyze_results.sh

# This provides:
# - Statistical summary (min, avg, median, 95th percentile, max)
# - Comparison table across all systems
# - Message size analysis
# - Performance recommendations
```

### Manual Results Review
```bash
# View individual result files
cat results_custom_queue.txt    # Custom queue results
cat results_rabbitmq.txt        # RabbitMQ results  
cat results_kafka.txt           # Kafka results
cat results_mqtt.txt            # MQTT results

# Monitor running benchmarks
tail -f results_custom_queue.txt
```

## Performance Considerations

### Message Size Analysis
- **Timestamp payload**: 19-20 bytes (e.g., "1690804800123456789")
- **Protocol overhead**:
  - Custom TCP: ~4 bytes (PUSH/MSG commands)
  - RabbitMQ AMQP: ~8-12 bytes (headers)
  - Kafka: ~20-30 bytes (record headers)
  - MQTT: ~2-4 bytes (minimal overhead)

### Benchmark Differences
- **Custom Queue**: 200,000 messages (high throughput test)
- **RabbitMQ**: 10,000 messages (moderate throughput)
- **Kafka/MQTT**: 1,000-10,000 messages (standard test)

These differences account for the inherent performance characteristics of each system.

## Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure no other services use the required ports
2. **Go modules**: Run `go mod tidy` if dependencies are missing
3. **Broker connectivity**: Check if message brokers are running and accessible
4. **Firewall**: Ensure ports are open (5672, 9092, 1883, 9000)

### Debug Mode
Uncomment detailed statistics in consumer files for verbose output:
```go
// Change this line in consumer files:
fmt.Println(avg.Microseconds())
// To this:
fmt.Printf("received %d msgs\nmin=%v max=%v avg=%v\n", len(lats), min, max, avg)
```

## License

MIT License - Feel free to modify and distribute.