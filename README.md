# Uptime Monitor

monitors https://riseworks.io for uptime by checking HTTP status codes every 15 seconds. Results are logged both to the terminal and to log file.

## Features

- Monitors website status every 30 seconds
- Colored terminal output
- Daily log files stored in `~/riseworks/logs/`
- Graceful shutdown with Ctrl+C

## Installation

1. Clone the repository and install dependencies:

```bash
go mod init rise-uptime
go get github.com/rs/zerolog
go get github.com/mattn/go-colorable
```
