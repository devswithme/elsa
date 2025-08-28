package watch

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	internalWatch "github.com/risoftinc/elsa/internal/watch"
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
	fmt.Printf("📁 Watching Go files in: %s\n", internalWatch.GetCurrentDir())
	fmt.Printf("⏱️ Restart delay: %v\n", watchDelay)
	fmt.Printf("🔍 File extensions: %v\n", watchExtensions)
	fmt.Printf("🚫 Excluded dirs: %v\n", watchExcludeDirs)
	fmt.Printf("Press Ctrl+C to stop watching\n\n")

	// Create watch options
	watchOptions := &internalWatch.WatchOptions{
		Extensions:   watchExtensions,
		ExcludeDirs:  watchExcludeDirs,
		OnFileChange: nil, // We'll handle file changes in the event loop
	}

	// Create file watcher
	fileWatcher, err := internalWatch.NewFileWatcher(watchOptions)
	if err != nil {
		log.Fatal("Error creating watcher:", err)
	}
	defer fileWatcher.Close()

	// Add directories to watch
	if err := fileWatcher.AddDirectoriesToWatch(); err != nil {
		log.Fatal("Error adding directories to watch:", err)
	}

	// Create process manager
	processManager := internalWatch.NewProcessManager()

	// Start initial command
	if err := processManager.StartCommand(command); err != nil {
		log.Fatal("Error starting initial command:", err)
	}

	// Start watching
	events, errors := fileWatcher.Watch()

	// Process management
	var debounce <-chan time.Time
	var isRestarting bool

	go func() {
		for {
			select {
			case event := <-events:
				if !isRestarting {
					fmt.Printf("📝 File changed: %s\n", event.Name)
					isRestarting = true
					debounce = time.After(watchDelay)
				}

			case <-debounce:
				fmt.Printf("🔄 Restarting command...\n")
				if err := processManager.RestartCommand(command, watchDelay); err != nil {
					fmt.Printf("❌ Error restarting command: %v\n", err)
				}
				debounce = nil
				isRestarting = false

			case err := <-errors:
				fmt.Printf("❌ Watcher error: %v\n", err)
			}
		}
	}()

	// Handle SIGINT / SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n👋 Stopping watch mode...")
	processManager.StopCommand()
}
