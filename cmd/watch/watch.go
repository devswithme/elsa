package watch

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"go.risoftinc.com/elsa/constants"
	internalWatch "go.risoftinc.com/elsa/internal/watch"
)

var (
	WatchCmd = &cobra.Command{
		Use:   constants.WatchCommandUsage,
		Short: constants.WatchCommandShort,
		Long:  constants.WatchCommandLong,
		Args:  cobra.MinimumNArgs(1),
		Run:   runWatch,
	}

	// Watch options
	watchExtensions  = []string{constants.DefaultWatchExtensions}
	watchExcludeDirs = strings.Split(constants.DefaultWatchExcludeDirs, ",")
	watchDelay       = constants.DefaultWatchDelay
)

func init() {
	WatchCmd.Flags().StringSliceVarP(&watchExtensions, constants.WatchFlagExt, constants.WatchFlagExtShort, watchExtensions, constants.WatchFlagExtUsage)
	WatchCmd.Flags().StringSliceVarP(&watchExcludeDirs, constants.WatchFlagExclude, constants.WatchFlagExcludeShort, watchExcludeDirs, constants.WatchFlagExcludeUsage)
	WatchCmd.Flags().DurationVarP(&watchDelay, constants.WatchFlagDelay, constants.WatchFlagDelayShort, watchDelay, constants.WatchFlagDelayUsage)
}

func runWatch(cmd *cobra.Command, args []string) {
	command := strings.Join(args, " ")
	fmt.Printf(constants.MsgWatchStarting+"\n", command)
	fmt.Printf(constants.MsgWatchDirectory+"\n", internalWatch.GetCurrentDir())
	fmt.Printf(constants.MsgWatchDelay+"\n", watchDelay)
	fmt.Printf(constants.MsgWatchExtensions+"\n", watchExtensions)
	fmt.Printf(constants.MsgWatchExcludedDirs+"\n", watchExcludeDirs)
	fmt.Printf(constants.MsgWatchPressCtrlC + "\n\n")

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
	defer func() {
		if err := fileWatcher.Close(); err != nil {
			fmt.Printf(constants.MsgWatchError+"\n", fmt.Sprintf("Error closing watcher: %v", err))
		}
	}()

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

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Process management
	var debounce <-chan time.Time
	var isRestarting bool

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf(constants.MsgWatchError+"\n", fmt.Sprintf("Panic in event loop: %v", r))
			}
		}()

		for {
			select {
			case <-ctx.Done():
				fmt.Println(constants.InfoEmoji + " Context cancelled, stopping event loop")
				return

			case event, ok := <-events:
				if !ok {
					fmt.Println(constants.InfoEmoji + " Events channel closed, stopping event loop")
					return
				}
				if !isRestarting {
					fmt.Printf(constants.MsgWatchFileChanged+"\n", event.Name)
					isRestarting = true
					debounce = time.After(watchDelay)
				}

			case <-debounce:
				fmt.Printf(constants.MsgWatchRestarting + "\n")
				if err := processManager.RestartCommand(command, watchDelay); err != nil {
					fmt.Printf(constants.MsgWatchRestartError+"\n", err)
				}
				debounce = nil
				isRestarting = false

			case err, ok := <-errors:
				if !ok {
					fmt.Println(constants.InfoEmoji + " Errors channel closed, stopping event loop")
					return
				}
				fmt.Printf(constants.MsgWatchError+"\n", err)

				// If it's a critical error, try to restart the watcher
				if strings.Contains(err.Error(), "watcher") || strings.Contains(err.Error(), "panic") {
					fmt.Println(constants.WarningEmoji + " Critical watcher error detected, attempting to restart...")
					// Note: In a production app, you might want to implement watcher restart here
				}
			}
		}
	}()

	// Handle SIGINT / SIGTERM with immediate cancellation
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for first signal
	<-sigChan

	fmt.Println("\n" + constants.MsgWatchStopping)

	// Cancel context immediately to stop goroutines
	cancel()

	// Stop process and close watcher
	processManager.StopCommand()
	fileWatcher.Close()

	// Give goroutines a moment to clean up
	time.Sleep(100 * time.Millisecond)

	// If we get another signal, force exit
	select {
	case <-sigChan:
		fmt.Println(constants.StopEmoji + " Force exit requested...")
		os.Exit(1)
	case <-time.After(2 * time.Second):
		// Normal exit after cleanup
	}
}
