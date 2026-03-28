# Server and Service Health Monitor

A lightweight, comprehensive monitoring tool written in Go. It actively monitors system resources, web services, and background processes, sending instant alerts via Enterprise WeChat when anomalies are detected. Perfect for use with simple cron jobs or scheduled tasks.

## Value Proposition

- **All-in-One Monitoring**: Combines system metrics (CPU, Memory, Disk) and application health (HTTP, TCP, PID) into a single, easy-to-deploy binary.
- **Instant Alerting**: Natively integrates with Enterprise WeChat to keep your team informed of issues immediately.
- **Robust & Performant**: Written in Go with concurrent execution, timeouts, and automatic retries (up to 3 attempts) to ensure reliable checks without bogging down your system.
- **Simple Deployment**: No complex agents or server setups required. Just a single binary, a YAML configuration file, and a cron job.

## Features

- **System Resource Monitoring**: Alerts when CPU, Memory, or Disk usage exceeds configurable thresholds. Automatically reports the top CPU-consuming processes when the CPU threshold is breached.
- **HTTP/HTTPS Health Checks**: Validates web endpoints by verifying expected status codes and expected response bodies.
- **TCP Health Checks**: Verifies connectivity to critical backend services via TCP host and port.
- **Process (PID) Checks**: Ensures important background jobs and daemons are alive by checking their PID files.

## Configuration File Structure (`config.yaml`)

```yaml
ewechat:
  corp_id: "your enterprise wechat corp_id"
  corp_secret: "your enterprise wechat corp_secret"
  agent_id: your enterprise wechat agent_id
  receivers: "your enterprise wechat receivers, for more receiver: receiver1|receiver2"

# System monitoring thresholds (Percentage)
cpu_usage_rate: 90.0  # Alerts when CPU usage exceeds 90%
mem_usage_rate: 90.0  # Alerts when Memory usage exceeds 90%
disk_usage_rate: 90.0 # Alerts when Disk usage exceeds 90%

# HTTP Service Checks
check_url:
  - url: "https://xxxx/api/health"
    method: get # default
    expect_status_code: 200 # default
    expect_body: "ok" # returns true if response body contains this string.
  - url: "https://xxxxx"
    method: post
    expect_status_code: 403
    expect_body: "As long as the returned string contains expect_body, it's a success"

# TCP Service Checks
check_tcp:
  - host: "example.com or ip"
    port: 443

# Process Checks (via PID file)
check_pid:
  - pid: "/var/run/my-service.pid"
```

## Instructions

### 1. Setup Configuration

```sh
cp config-example.yaml config.yaml
# Edit config.yaml to match your environment
```

### 2. Compile and Run Manually

```sh
go build -o server-check
./server-check
```

### 3. Automated Monitoring (Recommended)

Use `cron` on Linux/macOS to run the checks at regular intervals.

```sh
# Run every 5 minutes
*/5 * * * * cd /path/to/your/app && ./server-check
```

If you use Windows, you can use the **Task Scheduler** to run the executable periodically.
