package watch

import (
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

			case err := <-errors:
				fmt.Printf(constants.MsgWatchError+"\n", err)
			}
		}
	}()

	// Handle SIGINT / SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n" + constants.MsgWatchStopping)
	processManager.StopCommand()
}
