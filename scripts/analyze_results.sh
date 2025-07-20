#!/bin/bash

# Results Analysis Script
# Analyzes benchmark results and generates summary statistics

echo "=== Messaging Systems Benchmark Results Analysis ==="
echo ""

# Function to analyze a result file
analyze_file() {
    local file=$1
    local system=$2
    
    if [ ! -f "$file" ]; then
        echo "❌ $system: Results file not found ($file)"
        return
    fi
    
    if [ ! -s "$file" ]; then
        echo "❌ $system: Results file is empty ($file)"
        return
    fi
    
    echo "📊 $system Analysis:"
    
    # Extract latency values (assuming they're numeric lines)
    latencies=$(grep -E '^[0-9]+$' "$file")
    
    if [ -z "$latencies" ]; then
        echo "   ⚠️  No valid latency data found"
        echo "   📄 File preview:"
        head -10 "$file" | sed 's/^/      /'
        echo ""
        return
    fi
    
    # Calculate statistics using awk
    echo "$latencies" | awk '
    BEGIN {
        min = 999999999
        max = 0
        sum = 0
        count = 0
    }
    {
        if ($1 > 0) {
            sum += $1
            count++
            if ($1 < min) min = $1
            if ($1 > max) max = $1
            values[count] = $1
        }
    }
    END {
        if (count > 0) {
            avg = sum / count
            
            # Calculate median
            n = asort(values)
            if (n % 2 == 1) {
                median = values[int(n/2) + 1]
            } else {
                median = (values[int(n/2)] + values[int(n/2) + 1]) / 2
            }
            
            # Calculate 95th percentile
            p95_index = int(n * 0.95)
            if (p95_index == 0) p95_index = 1
            p95 = values[p95_index]
            
            printf "   📈 Iterations: %d\n", count
            printf "   ⚡ Min Latency: %d μs (%.3f ms)\n", min, min/1000
            printf "   📊 Avg Latency: %d μs (%.3f ms)\n", avg, avg/1000
            printf "   📍 Median Latency: %d μs (%.3f ms)\n", median, median/1000
            printf "   📊 95th Percentile: %d μs (%.3f ms)\n", p95, p95/1000
            printf "   🔥 Max Latency: %d μs (%.3f ms)\n", max, max/1000
            printf "   📏 Latency Range: %d μs\n", max - min
        } else {
            print "   ❌ No valid data points found"
        }
    }'
    echo ""
}

# Analyze each system
analyze_file "results_custom_queue.txt" "Custom Queue Server"
analyze_file "results_rabbitmq.txt" "RabbitMQ"
analyze_file "results_kafka.txt" "Apache Kafka"
analyze_file "results_mqtt.txt" "MQTT"

echo "=== Summary Comparison ==="
echo ""

# Create a comparison table
echo "| System              | Avg Latency (μs) | Avg Latency (ms) | Status |"
echo "|---------------------|------------------|------------------|--------|"

for system_file in "Custom Queue:results_custom_queue.txt" "RabbitMQ:results_rabbitmq.txt" "Kafka:results_kafka.txt" "MQTT:results_mqtt.txt"; do
    IFS=':' read -r system_name file <<< "$system_file"
    
    if [ -f "$file" ] && [ -s "$file" ]; then
        avg=$(grep -E '^[0-9]+$' "$file" | awk '{sum+=$1; count++} END {if(count>0) printf "%.0f", sum/count; else print "N/A"}')
        if [ "$avg" != "N/A" ]; then
            avg_ms=$(echo "scale=3; $avg/1000" | bc -l 2>/dev/null || echo "0")
            printf "| %-19s | %16s | %16s | ✅     |\n" "$system_name" "$avg" "$avg_ms"
        else
            printf "| %-19s | %16s | %16s | ❌     |\n" "$system_name" "N/A" "N/A"
        fi
    else
        printf "| %-19s | %16s | %16s | ⚠️      |\n" "$system_name" "No Data" "No Data"
    fi
done

echo ""
echo "=== Message Size Analysis ==="
echo "All systems use the same message format:"
echo "• Payload: Unix timestamp in nanoseconds (string)"
echo "• Size: ~19-20 bytes (e.g., '1690804800123456789')"
echo "• Protocol overhead varies by system:"
echo "  - Custom TCP: ~4 bytes (PUSH/MSG commands)"
echo "  - RabbitMQ AMQP: ~8-12 bytes (headers)"
echo "  - Kafka: ~20-30 bytes (record headers)"
echo "  - MQTT: ~2-4 bytes (minimal overhead)"
echo ""

echo "=== Recommendations ==="
echo "🚀 Lowest latency systems are typically better for real-time applications"
echo "📊 Consider throughput vs latency tradeoffs for your use case"
echo "🔧 Protocol overhead affects small message performance"
echo "💾 Persistence and durability features add latency but improve reliability"
echo ""

# Check if bc is available for calculations
if ! command -v bc &> /dev/null; then
    echo "💡 Install 'bc' for enhanced mathematical calculations: sudo apt-get install bc"
fi