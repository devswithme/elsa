package watch

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// ProcessManager handles starting, stopping, and monitoring processes
type ProcessManager struct {
	currentProcess *exec.Cmd
}

// NewProcessManager creates a new ProcessManager instance
func NewProcessManager() *ProcessManager {
	return &ProcessManager{}
}

// StartCommand starts a new command and stores the process reference
func (pm *ProcessManager) StartCommand(command string) error {
	fmt.Printf("▶️  Running: %s\n", command)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		// Use universal shell for Unix systems (Linux, macOS, BSD)
		cmd = exec.Command("sh", "-c", command)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	pm.currentProcess = cmd

	// Monitor process in background
	go func() {
		if err := cmd.Wait(); err != nil {
			// Don't show error for terminated commands
			if err.Error() != "signal: interrupt" && err.Error() != "exit status 1" {
				fmt.Printf("❌ Command exited with error: %v\n", err)
			}
		} else {
			fmt.Println("✅ Command completed")
		}
	}()

	return nil
}

// StopCommand stops the current running process
func (pm *ProcessManager) StopCommand() {
	if pm.currentProcess == nil || pm.currentProcess.Process == nil {
		return
	}

	pid := pm.currentProcess.Process.Pid
	fmt.Printf("🛑 Killing process PID: %d\n", pid)

	if runtime.GOOS == "windows" {
		// Use taskkill on Windows
		_ = exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", pid)).Run()
	} else {
		// For Unix systems (Linux, macOS, BSD)
		// 1. Try graceful kill first
		_ = pm.currentProcess.Process.Signal(syscall.SIGTERM)
		time.Sleep(500 * time.Millisecond)

		// 2. Force kill if still alive
		if IsProcessRunning(pid) {
			fmt.Printf("⚠️ Process still running, force killing...\n")
			_ = pm.currentProcess.Process.Kill()
		}
	}

	// Wait a bit to ensure process is really dead
	time.Sleep(1 * time.Second)
	pm.currentProcess = nil
}

// RestartCommand stops the current command and starts a new one
func (pm *ProcessManager) RestartCommand(command string, delay time.Duration) error {
	fmt.Println("🔄 Restarting...")
	pm.StopCommand()

	// Wait a bit to ensure port release
	time.Sleep(2*time.Second + delay)

	return pm.StartCommand(command)
}

// IsProcessRunning checks if a process is still running (cross-platform)
func IsProcessRunning(pid int) bool {
	if runtime.GOOS == "windows" {
		// Windows: use tasklist
		cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid))
		output, _ := cmd.Output()
		return strings.Contains(string(output), fmt.Sprintf("%d", pid))
	} else {
		// Unix: use kill -0 to check process
		cmd := exec.Command("kill", "-0", fmt.Sprintf("%d", pid))
		return cmd.Run() == nil
	}
}

// GetCurrentProcess returns the current running process
func (pm *ProcessManager) GetCurrentProcess() *exec.Cmd {
	return pm.currentProcess
}
