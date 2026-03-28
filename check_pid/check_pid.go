package check_pid

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// CheckHealth reads the pid from the specified file and checks if the process is running.
func CheckHealth(ctx context.Context, pidFile string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return false, fmt.Errorf("failed to read pid file: %w", err)
	}

	pidStr := strings.TrimSpace(string(data))
	if pidStr == "" {
		return false, fmt.Errorf("pid file is empty")
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return false, fmt.Errorf("invalid pid in file: %w", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		// on Unix systems FindProcess always succeeds, but returning error just in case for cross-platform
		return false, fmt.Errorf("failed to find process: %w", err)
	}

	// Sending signal 0 checks if the process exists and we have permissions
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false, nil
	}

	return true, nil
}
