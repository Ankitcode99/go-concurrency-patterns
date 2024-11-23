package main

import "time"

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Service   string    `json:"service"`
	Message   string    `json:"message"`
	IP        string    `json:"ip"`
}

// LogAnalytics represents processed analytics for logs
type LogAnalytics struct {
	TotalErrors    int
	UniqueIPs      map[string]bool
	ErrorPatterns  map[string]int
	ServiceMetrics map[string]int
	ProcessedBy    int // Worker ID that processed this batch
}
