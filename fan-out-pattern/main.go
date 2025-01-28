package main

/**
Fan-out patterns allow us to essentially split our single input channel into multiple output channels.
This is a useful pattern to distribute work items into multiple uniform actors.
*/

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

// simulateLogSource generates sample logs (in real life, this could be reading from Kafka, file, etc.)
func simulateLogSource(ctx context.Context) <-chan LogEntry {
	out := make(chan LogEntry)

	go func() {
		defer close(out)

		services := []string{"auth-service", "payment-service", "user-service"}
		levels := []string{"ERROR", "INFO", "WARN"}
		ips := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3", "192.168.1.4"}

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Millisecond):
				entry := LogEntry{
					Timestamp: time.Now(),
					Level:     levels[time.Now().UnixNano()%int64(len(levels))],
					Service:   services[time.Now().UnixNano()%int64(len(services))],
					Message:   fmt.Sprintf("Sample log message %d", time.Now().UnixNano()),
					IP:        ips[time.Now().UnixNano()%int64(len(ips))],
				}
				out <- entry
			}
		}
	}()

	return out
}

// processLogs implements the fan-out pattern for log processing
func processLogs(ctx context.Context, workerID int, input <-chan LogEntry) <-chan LogAnalytics {
	out := make(chan LogAnalytics)

	go func() {
		defer close(out)

		// Process logs in batches
		batchSize := 10
		currentBatch := make([]LogEntry, 0, batchSize)
		errorPattern := regexp.MustCompile(`error|exception|failed|timeout|invalid`)

		for {
			select {
			case <-ctx.Done():
				if len(currentBatch) > 0 {
					out <- analyzeBatch(currentBatch, workerID, errorPattern)
				}
				return
			case log, ok := <-input:
				if !ok {
					if len(currentBatch) > 0 {
						out <- analyzeBatch(currentBatch, workerID, errorPattern)
					}
					return
				}

				currentBatch = append(currentBatch, log)

				if len(currentBatch) >= batchSize {
					out <- analyzeBatch(currentBatch, workerID, errorPattern)
					currentBatch = make([]LogEntry, 0, batchSize)
				}
			}
		}
	}()

	return out
}

// analyzeBatch processes a batch of logs and returns analytics
func analyzeBatch(batch []LogEntry, workerID int, errorPattern *regexp.Regexp) LogAnalytics {
	analytics := LogAnalytics{
		UniqueIPs:      make(map[string]bool),
		ErrorPatterns:  make(map[string]int),
		ServiceMetrics: make(map[string]int),
		ProcessedBy:    workerID,
	}

	for _, log := range batch {
		// Track unique IPs
		analytics.UniqueIPs[log.IP] = true

		// Count by service
		analytics.ServiceMetrics[log.Service]++

		// Track errors
		if log.Level == "ERROR" {
			analytics.TotalErrors++

			// Extract error patterns
			if matches := errorPattern.FindAllString(strings.ToLower(log.Message), -1); matches != nil {
				for _, match := range matches {
					analytics.ErrorPatterns[match]++
				}
			}
		}
	}

	return analytics
}

// mergeResults combines results from all workers
func mergeResults(ctx context.Context, channels []<-chan LogAnalytics) <-chan LogAnalytics {
	var wg sync.WaitGroup
	out := make(chan LogAnalytics)

	// Start an output goroutine for each input channel
	output := func(c <-chan LogAnalytics) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go output(c)
	}

	// Start a goroutine to close out once all output goroutines are done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create input channel for logs
	logSource := simulateLogSource(ctx)

	// Fan out to multiple workers
	numWorkers := 5
	workers := make([]<-chan LogAnalytics, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = processLogs(ctx, i+1, logSource)
	}

	// Merge results from all workers
	results := mergeResults(ctx, workers)

	// Aggregate final results
	totalErrors := 0
	uniqueIPs := make(map[string]bool)
	errorPatterns := make(map[string]int)
	serviceMetrics := make(map[string]int)

	// Process results
	for result := range results {
		fmt.Printf("Result received from worker %d : %v\n", result.ProcessedBy, result)
		totalErrors += result.TotalErrors

		// Merge unique IPs
		for ip := range result.UniqueIPs {
			uniqueIPs[ip] = true
		}

		// Merge error patterns
		for pattern, count := range result.ErrorPatterns {
			errorPatterns[pattern] += count
		}

		// Merge service metrics
		for service, count := range result.ServiceMetrics {
			serviceMetrics[service] += count
		}

		// Print worker activity
		fmt.Printf("Worker %d processed batch: %d errors\n", result.ProcessedBy, result.TotalErrors)
	}

	// Print final analysis
	fmt.Printf("\nFinal Analysis:\n")
	fmt.Printf("Total Errors: %d\n", totalErrors)
	fmt.Printf("Unique IPs: %d\n", len(uniqueIPs))
	fmt.Printf("\nError Patterns:\n")
	for pattern, count := range errorPatterns {
		fmt.Printf("- %s: %d\n", pattern, count)
	}
	fmt.Printf("\nService Metrics:\n")
	for service, count := range serviceMetrics {
		fmt.Printf("- %s: %d logs\n", service, count)
	}
}
