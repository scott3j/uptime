package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
)

func ensureLogDirectory() string {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Failed to get home directory: %v", err))
	}

	// Create the full path for the logs directory
	logDir := filepath.Join(homeDir, "riseworks", "logs")

	// Create all directories in the path if they don't exist
	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		panic(fmt.Sprintf("Failed to create log directory: %v", err))
	}

	return logDir
}

func main() {
	// Ensure log directory exists and get its path
	logDir := ensureLogDirectory()

	// Create log file with current timestamp in the logs directory
	currentTime := time.Now().Format("2006-01-02")
	logPath := filepath.Join(logDir, fmt.Sprintf("site_monitor_%s.log", currentTime))

	logFile, err := os.OpenFile(
		logPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file: %v", err))
	}
	defer logFile.Close()

	// Create multi-writer for both console and file
	consoleWriter := zerolog.ConsoleWriter{
		Out:        colorable.NewColorableStdout(),
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	logger := zerolog.New(multi).With().Timestamp().Logger()

	// Log startup message with file location
	logger.Info().
		Str("log_file", logPath).
		Msg("Starting site monitor")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create ticker for 30 second intervals
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Log initial check
	checkStatus(client, &logger)

	// Main loop
	for {
		select {
		case <-ticker.C:
			checkStatus(client, &logger)
		case <-sigChan:
			logger.Info().Msg("Shutting down monitor")
			return
		}
	}
}

func checkStatus(client *http.Client, logger *zerolog.Logger) {
	resp, err := client.Get("https://riseworks.io")
	if err != nil {
		logger.Error().
			Err(err).
			Str("url", "https://riseworks.io").
			Msg("Failed to make HTTP request")
		return
	}
	defer resp.Body.Close()

	logger.Info().
		Str("url", "https://riseworks.io").
		Int("status_code", resp.StatusCode).
		Str("status", resp.Status).
		Msg("Site status check")
}
