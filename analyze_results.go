package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type TestResult struct {
	System    string
	TestFile  string
	Count     int
	AvgMicros float64
	MinMicros int64
	MaxMicros int64
	P50Micros int64
	P95Micros int64
	P99Micros int64
}

func main() {
	systems := []string{"simplified", "rabbitmq", "mqtt", "kafka"}
	var allResults []TestResult

	for _, system := range systems {
		resultsDir := filepath.Join("results", system)
		if _, err := os.Stat(resultsDir); os.IsNotExist(err) {
			fmt.Printf("No results directory found for %s\n", system)
			continue
		}

		files, err := filepath.Glob(filepath.Join(resultsDir, "test_*.txt"))
		if err != nil {
			log.Printf("Error reading %s results: %v", system, err)
			continue
		}

		if len(files) == 0 {
			fmt.Printf("No test files found for %s\n", system)
			continue
		}

		fmt.Printf("\n=== %s Results ===\n", strings.ToUpper(system))
		systemResults := []TestResult{}

		for _, file := range files {
			result, err := analyzeFile(system, file)
			if err != nil {
				log.Printf("Error analyzing %s: %v", file, err)
				continue
			}
			systemResults = append(systemResults, result)
		}

		if len(systemResults) > 0 {
			// Sort by filename for consistent output
			sort.Slice(systemResults, func(i, j int) bool {
				return systemResults[i].TestFile < systemResults[j].TestFile
			})

			// Print individual test results
			for i, result := range systemResults {
				fmt.Printf("Test %d: avg=%.1fμs, min=%dμs, max=%dμs, p50=%dμs, p95=%dμs, p99=%dμs (%d msgs)\n",
					i+1, result.AvgMicros, result.MinMicros, result.MaxMicros,
					result.P50Micros, result.P95Micros, result.P99Micros, result.Count)
			}

			// Calculate overall statistics
			var totalLatency float64
			var totalCount int
			for _, result := range systemResults {
				totalLatency += result.AvgMicros * float64(result.Count)
				totalCount += result.Count
			}
			overallAvg := totalLatency / float64(totalCount)

			fmt.Printf("\nOverall %s: %.1fμs average across %d tests (%d total messages)\n",
				system, overallAvg, len(systemResults), totalCount)

			allResults = append(allResults, systemResults...)
		}
	}

	// Print comparison summary
	if len(allResults) > 0 {
		printComparisonSummary(allResults)
	}
}

func analyzeFile(system, filename string) (TestResult, error) {
	file, err := os.Open(filename)
	if err != nil {
		return TestResult{}, err
	}
	defer file.Close()

	var latencies []int64
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		latency, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			log.Printf("Invalid latency value in %s: %s", filename, line)
			continue
		}
		latencies = append(latencies, latency)
	}

	if err := scanner.Err(); err != nil {
		return TestResult{}, err
	}

	if len(latencies) == 0 {
		return TestResult{}, fmt.Errorf("no valid latency measurements found")
	}

	// Sort for percentile calculations
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	// Calculate statistics
	var sum int64
	min := latencies[0]
	max := latencies[len(latencies)-1]

	for _, lat := range latencies {
		sum += lat
	}

	avg := float64(sum) / float64(len(latencies))
	p50 := percentile(latencies, 0.50)
	p95 := percentile(latencies, 0.95)
	p99 := percentile(latencies, 0.99)

	return TestResult{
		System:    system,
		TestFile:  filepath.Base(filename),
		Count:     len(latencies),
		AvgMicros: avg,
		MinMicros: min,
		MaxMicros: max,
		P50Micros: p50,
		P95Micros: p95,
		P99Micros: p99,
	}, nil
}

func percentile(sortedData []int64, p float64) int64 {
	if len(sortedData) == 0 {
		return 0
	}
	if p <= 0 {
		return sortedData[0]
	}
	if p >= 1 {
		return sortedData[len(sortedData)-1]
	}

	index := p * float64(len(sortedData)-1)
	lower := int(index)
	upper := lower + 1

	if upper >= len(sortedData) {
		return sortedData[lower]
	}

	// Linear interpolation
	weight := index - float64(lower)
	return int64(float64(sortedData[lower])*(1-weight) + float64(sortedData[upper])*weight)
}

func printComparisonSummary(allResults []TestResult) {
	fmt.Printf("\n=== SYSTEM COMPARISON ===\n")

	systemStats := make(map[string][]float64)
	for _, result := range allResults {
		systemStats[result.System] = append(systemStats[result.System], result.AvgMicros)
	}

	type SystemSummary struct {
		Name      string
		AvgMicros float64
		TestCount int
	}

	var summaries []SystemSummary
	for system, latencies := range systemStats {
		var sum float64
		for _, lat := range latencies {
			sum += lat
		}
		avg := sum / float64(len(latencies))
		summaries = append(summaries, SystemSummary{
			Name:      system,
			AvgMicros: avg,
			TestCount: len(latencies),
		})
	}

	// Sort by average latency (lowest first)
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].AvgMicros < summaries[j].AvgMicros
	})

	fmt.Printf("Rank | System     | Avg Latency | Tests\n")
	fmt.Printf("-----|------------|-------------|------\n")
	for i, summary := range summaries {
		fmt.Printf("%-4d | %-10s | %8.1fμs | %d\n",
			i+1, summary.Name, summary.AvgMicros, summary.TestCount)
	}

	if len(summaries) > 1 {
		fastest := summaries[0]
		fmt.Printf("\n%s is the fastest system tested.\n", strings.ToUpper(fastest.Name))

		for i := 1; i < len(summaries); i++ {
			ratio := summaries[i].AvgMicros / fastest.AvgMicros
			fmt.Printf("%s is %.1fx slower than %s\n",
				strings.ToUpper(summaries[i].Name), ratio, strings.ToUpper(fastest.Name))
		}
	}
}
