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

Examples:
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
	fmt.Printf("Press Ctrl+C to stop watching\n\n")

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

	// Process management
	var currentProcess *exec.Cmd
	var debounce <-chan time.Time
	var isRestarting bool

	// Start initial command
	startCommand(command, &currentProcess)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if shouldRestart(event) && !isRestarting {
					fmt.Printf("📝 File changed: %s\n", event.Name)
					isRestarting = true
					debounce = time.After(watchDelay)
				}

			case <-debounce:
				fmt.Println("🔄 Restarting...")
				stopCommand(currentProcess)
				// Tunggu sebentar untuk memastikan port release
				time.Sleep(2*time.Second + watchDelay)
				startCommand(command, &currentProcess)
				debounce = nil
				isRestarting = false

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("❌ Watcher error: %v\n", err)
			}
		}
	}()

	// Handle SIGINT / SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n👋 Stopping watch mode...")
	stopCommand(currentProcess)
}

func startCommand(command string, currentProcess **exec.Cmd) {
	fmt.Printf("▶️  Running: %s\n", command)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		fmt.Printf("❌ Error starting command: %v\n", err)
		return
	}

	*currentProcess = cmd

	go func() {
		if err := cmd.Wait(); err != nil {
			// Jangan tampilkan error untuk command yang di-terminate
			if err.Error() != "signal: interrupt" && err.Error() != "exit status 1" {
				fmt.Printf("❌ Command exited with error: %v\n", err)
			}
		} else {
			fmt.Println("✅ Command completed")
		}
	}()
}

func stopCommand(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	pid := cmd.Process.Pid
	fmt.Printf("🛑 Killing process PID: %d\n", pid)

	if runtime.GOOS == "windows" {
		// pakai taskkill di Windows
		_ = exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", pid)).Run()
	} else {
		// pakai Process.Kill() yang cross-platform
		_ = cmd.Process.Kill()
	}

	// Tunggu sebentar untuk memastikan process benar-benar mati
	time.Sleep(1 * time.Second)
}

func addDirectoriesToWatch(watcher *fsnotify.Watcher) error {
	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirName := filepath.Base(path)
			for _, excluded := range watchExcludeDirs {
				if dirName == excluded {
					return filepath.SkipDir
				}
			}
			if err := watcher.Add(path); err != nil {
				fmt.Printf("⚠️  Warning: Could not watch directory %s: %v\n", path, err)
			}
		}
		return nil
	})
}

func shouldRestart(event fsnotify.Event) bool {
	// Hanya restart pada Write dan Create, bukan Remove dan Rename
	if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
		return false
	}
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
