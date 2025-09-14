package watch

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"go.risoftinc.com/elsa/constants"
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
	fmt.Printf(constants.MsgWatchRunning+"\n", command)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(constants.WindowsShell, constants.WindowsShellArgs, command)
	} else {
		// Use universal shell for Unix systems (Linux, macOS, BSD)
		cmd = exec.Command(constants.UnixShell, constants.UnixShellArgs, command)
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
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf(constants.MsgWatchError+"\n", fmt.Sprintf("Panic in process monitor: %v", r))
			}
		}()

		if err := cmd.Wait(); err != nil {
			// Don't show error for terminated commands
			if err.Error() != "signal: interrupt" &&
				err.Error() != "exit status 1" &&
				err.Error() != "signal: killed" {
				fmt.Printf(constants.MsgWatchExitedWithError+"\n", err)
			}
		} else {
			fmt.Println(constants.MsgWatchCompleted)
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
	fmt.Printf(constants.MsgWatchKillingProcess+"\n", pid)

	if runtime.GOOS == "windows" {
		// Use taskkill on Windows
		_ = exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", pid)).Run()
	} else {
		// For Unix systems (Linux, macOS, BSD)
		// 1. Try graceful kill first
		_ = pm.currentProcess.Process.Signal(syscall.SIGTERM)
		time.Sleep(constants.ProcessKillDelay)

		// 2. Force kill if still alive
		if IsProcessRunning(pid) {
			fmt.Printf(constants.MsgWatchForceKilling + "\n")
			_ = pm.currentProcess.Process.Kill()
		}
	}

	// Wait a bit to ensure process is really dead
	time.Sleep(constants.ProcessWaitDelay)
	pm.currentProcess = nil
}

// RestartCommand stops the current command and starts a new one
func (pm *ProcessManager) RestartCommand(command string, delay time.Duration) error {
	fmt.Println(constants.MsgWatchRestartingProcess)
	pm.StopCommand()

	// Wait a bit to ensure port release
	time.Sleep(constants.RestartWaitDelay + delay)

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
