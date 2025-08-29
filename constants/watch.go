package constants

import "time"

// Watch file extensions constants
const (
	// DefaultWatchExtensions are the default file extensions to watch
	DefaultWatchExtensions = ".go"
)

// Watch exclude directories constants
const (
	// DefaultWatchExcludeDirs are the default directories to exclude from watching
	DefaultWatchExcludeDirs = ".git,vendor,tmp,temp,build,dist,bin,pkg,.vscode,.idea,coverage,testdata"
)

// Watch timing constants
const (
	// DefaultWatchDelay is the default delay before restarting (500ms)
	DefaultWatchDelay = 500 * time.Millisecond

	// ProcessKillDelay is the delay before force killing a process (500ms)
	ProcessKillDelay = 500 * time.Millisecond

	// ProcessWaitDelay is the delay to wait for process to die (1s)
	ProcessWaitDelay = 1 * time.Second

	// RestartWaitDelay is the delay before restarting a command (2s)
	RestartWaitDelay = 2 * time.Second
)

// Watch command constants
const (
	// WatchCommandUsage is the usage description for watch command
	WatchCommandUsage = "watch [command]"

	// WatchCommandShort is the short description for watch command
	WatchCommandShort = "Watch Go files and auto-restart on changes"

	// WatchCommandLong is the long description for watch command
	WatchCommandLong = `Watch Go files in the current directory and automatically restart the specified command when changes are detected.
		
Only Go files (*.go) are monitored to avoid unnecessary restarts from temporary files or uploads.
Automatically excludes common Go development folders like vendor, build, bin, pkg, etc.

Examples:
  elsa watch "go run main.go"
  elsa watch "go build && ./elsa"
  elsa watch "go test ./..."`

	// WatchFlagExt is the flag name for file extensions
	WatchFlagExt = "ext"

	// WatchFlagExtShort is the short flag name for file extensions
	WatchFlagExtShort = "e"

	// WatchFlagExtUsage is the usage description for extensions flag
	WatchFlagExtUsage = "File extensions to watch"

	// WatchFlagExclude is the flag name for exclude directories
	WatchFlagExclude = "exclude"

	// WatchFlagExcludeShort is the short flag name for exclude directories
	WatchFlagExcludeShort = "x"

	// WatchFlagExcludeUsage is the usage description for exclude flag
	WatchFlagExcludeUsage = "Directories to exclude from watching"

	// WatchFlagDelay is the flag name for restart delay
	WatchFlagDelay = "delay"

	// WatchFlagDelayShort is the short flag name for restart delay
	WatchFlagDelayShort = "d"

	// WatchFlagDelayUsage is the usage description for delay flag
	WatchFlagDelayUsage = "Delay before restarting (e.g., 500ms, 1s)"
)

// Watch message constants
const (
	// MsgWatchStarting is the message when watch mode starts
	MsgWatchStarting = RocketEmoji + " Starting watch mode for: %s"

	// MsgWatchDirectory is the message showing watched directory
	MsgWatchDirectory = FolderEmoji + " Watching Go files in: %s"

	// MsgWatchDelay is the message showing restart delay
	MsgWatchDelay = TimerEmoji + " Restart delay: %v"

	// MsgWatchExtensions is the message showing watched extensions
	MsgWatchExtensions = MagnifyingGlassEmoji + " File extensions: %v"

	// MsgWatchExcludedDirs is the message showing excluded directories
	MsgWatchExcludedDirs = ProhibitedEmoji + " Excluded dirs: %v"

	// MsgWatchPressCtrlC is the message to stop watching
	MsgWatchPressCtrlC = "Press Ctrl+C to stop watching"

	// MsgWatchFileChanged is the message when a file changes
	MsgWatchFileChanged = FileChangeEmoji + " File changed: %s"

	// MsgWatchRestarting is the message when restarting command
	MsgWatchRestarting = RestartEmoji + " Restarting command..."

	// MsgWatchRestartError is the message when restart fails
	MsgWatchRestartError = ErrorEmoji + " Error restarting command: %v"

	// MsgWatchError is the message when watcher has an error
	MsgWatchError = ErrorEmoji + " Watcher error: %v"

	// MsgWatchStopping is the message when stopping watch mode
	MsgWatchStopping = WaveEmoji + " Stopping watch mode..."

	// MsgWatchRunning is the message when command starts running
	MsgWatchRunning = PlayEmoji + " Running: %s"

	// MsgWatchCompleted is the message when command completes
	MsgWatchCompleted = SuccessEmoji + " Command completed"

	// MsgWatchExitedWithError is the message when command exits with error
	MsgWatchExitedWithError = ErrorEmoji + " Command exited with error: %v"

	// MsgWatchKillingProcess is the message when killing a process
	MsgWatchKillingProcess = StopEmoji + " Killing process PID: %d"

	// MsgWatchForceKilling is the message when force killing a process
	MsgWatchForceKilling = WarningEmoji + " Process still running, force killing..."

	// MsgWatchRestartingProcess is the message when restarting process
	MsgWatchRestartingProcess = RestartEmoji + " Restarting..."

	// MsgWatchWarning is the message for warnings
	MsgWatchWarning = WarningEmoji + " Warning: Could not watch directory %s: %v"
)
