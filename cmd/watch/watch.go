package watch

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var (
	WatchCmd = &cobra.Command{
		Use:   "watch [command]",
		Short: "Watch Go files and auto-restart on changes",
		Long: `Watch Go files in the current directory and automatically restart the specified command when changes are detected.
		
Only Go files (*.go) are monitored to avoid unnecessary restarts from temporary files or uploads.
Automatically excludes common Go development folders like vendor, build, bin, pkg, etc.

Example:
  elsa watch "go run main.go"
  elsa watch "go build && ./elsa"
  elsa watch "go test ./..."`,
		Args: cobra.MinimumNArgs(1),
		Run:  runWatch,
	}

	// Watch options
	watchExtensions  = []string{".go"}
	watchExcludeDirs = []string{".git", "vendor", "tmp", "temp", "build", "dist", "bin", "pkg", ".vscode", ".idea", "coverage", "testdata"}
	watchDelay       = 500 * time.Millisecond
)

func init() {
	WatchCmd.Flags().StringSliceVarP(&watchExtensions, "ext", "e", watchExtensions, "File extensions to watch")
	WatchCmd.Flags().StringSliceVarP(&watchExcludeDirs, "exclude", "x", watchExcludeDirs, "Directories to exclude from watching")
	WatchCmd.Flags().DurationVarP(&watchDelay, "delay", "d", watchDelay, "Delay before restarting (e.g., 500ms, 1s)")
}

func runWatch(cmd *cobra.Command, args []string) {
	command := strings.Join(args, " ")
	fmt.Printf("🚀 Starting watch mode for: %s\n", command)
	fmt.Printf("📁 Watching Go files in: %s\n", getCurrentDir())
	fmt.Printf("⏱️ Restart delay: %v\n", watchDelay)
	fmt.Printf("🔍 File extensions: %v\n", watchExtensions)
	fmt.Printf("🚫 Excluded dirs: %v\n", watchExcludeDirs)
	fmt.Println("Press Ctrl+C to stop watching")
	fmt.Println()

	// Create watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error creating watcher:", err)
	}
	defer watcher.Close()

	// Add directories to watch
	if err := addDirectoriesToWatch(watcher); err != nil {
		log.Fatal("Error adding directories to watch:", err)
	}

	// Channel for restart signals and process management
	restartChan := make(chan bool, 1)
	processChan := make(chan *exec.Cmd, 1)
	stopChan := make(chan bool, 1)

	// Process management
	var currentProcess *exec.Cmd

	// Start the initial command
	go runCommand(command, restartChan, processChan, stopChan)

	// Handle file system events
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Only restart on Go file changes
				if shouldRestart(event) {
					fmt.Printf("📝 File changed: %s\n", event.Name)
					// Stop current process before restarting
					if currentProcess != nil && currentProcess.Process != nil {
						fmt.Println("🛑 Stopping current process for restart...")
						forceKillProcess(currentProcess)
						// Wait a bit for process to fully terminate and port to be released
						time.Sleep(500 * time.Millisecond)
					}
					select {
					case restartChan <- true:
					default:
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("❌ Watcher error: %v\n", err)
			}
		}
	}()

	// Handle process management
	go func() {
		for {
			select {
			case process := <-processChan:
				currentProcess = process
			case <-stopChan:
				return
			}
		}
	}()

	// Handle restart signals
	go func() {
		for range restartChan {
			time.Sleep(watchDelay)
			fmt.Println("🔄 Restarting...")
			go runCommand(command, restartChan, processChan, stopChan)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n👋 Stopping watch mode...")
	stopChan <- true
	watcher.Close()
}

func runCommand(command string, restartChan chan bool, processChan chan *exec.Cmd, stopChan chan bool) {
	fmt.Printf("▶️  Running: %s\n", command)

	// Split command for cross-platform compatibility
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Send process to main goroutine for management
	select {
	case processChan <- cmd:
	default:
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("❌ Error starting command: %v\n", err)
		fmt.Println("⏳ Waiting for file changes...")
		return
	}

	// Check if we should stop
	select {
	case <-stopChan:
		fmt.Println("🛑 Stopping current process...")
		if cmd.Process != nil {
			forceKillProcess(cmd)
		}
		return
	default:
	}

	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		fmt.Printf("❌ Command failed: %v\n", err)
		fmt.Println("⏳ Waiting for file changes...")
	} else {
		fmt.Println("✅ Command completed successfully")
		fmt.Println("⏳ Waiting for file changes...")
	}
}

func addDirectoriesToWatch(watcher *fsnotify.Watcher) error {
	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip excluded directories
		if info.IsDir() {
			dirName := filepath.Base(path)
			for _, excluded := range watchExcludeDirs {
				if dirName == excluded {
					return filepath.SkipDir
				}
			}

			// Add directory to watcher
			if err := watcher.Add(path); err != nil {
				fmt.Printf("⚠️  Warning: Could not watch directory %s: %v\n", path, err)
			}
		}

		return nil
	})
}

func shouldRestart(event fsnotify.Event) bool {
	// Only restart on Go file changes
	if event.Op&fsnotify.Write == 0 && event.Op&fsnotify.Create == 0 {
		return false
	}

	// Check if it's a Go file
	ext := filepath.Ext(event.Name)
	for _, watchExt := range watchExtensions {
		if ext == watchExt {
			return true
		}
	}

	return false
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

// forceKillProcess forcefully terminates a process and its children
func forceKillProcess(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}

	// Kill the main process
	cmd.Process.Kill()

	// On Windows, we need to be more aggressive
	if runtime.GOOS == "windows" {
		// Use taskkill to force kill the process tree
		killCmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", cmd.Process.Pid))
		killCmd.Run()
	}

	// Wait for process to terminate
	cmd.Wait()
}
