# Smart Idle Shutdown

Smart Idle Shutdown is a console application designed to automatically shut down a computer based on specific load conditions. It provides flexibility through command-line flags to customize the behavior according to your preferences.

## Usage

```bash
./smart_idle_shutdown [flags]
```

## Available Flags

- `-interval`: Check interval in seconds (default: 1 second)
- `-threshold`: Load Threshold in Percent (default: 30.0)
- `-consecutive`: Number of consecutive times the load should be below the threshold (default: 3)

## Example

```bash
-interval 120 -threshold 80.5 -consecutive 5
```

This example sets the interval for checking to 120 seconds, the load threshold to 80.5, and requires the load to be below the threshold for 5 consecutive times before initiating a shutdown. So that's 10 minutes of downtime in a row

## Purpose

The primary purpose of Smart Idle Shutdown is to save energy by automatically shutting down the computer when the system load falls below a specified threshold for a certain duration. This is especially useful for scenarios where the computer is frequently left idle.
