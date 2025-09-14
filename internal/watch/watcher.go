package watch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"go.risoftinc.com/elsa/constants"
)

// FileWatcher handles file system watching with configurable options
type FileWatcher struct {
	watcher      *fsnotify.Watcher
	extensions   []string
	excludeDirs  []string
	onFileChange func(string)
}

// WatchOptions configures the file watcher behavior
type WatchOptions struct {
	Extensions   []string
	ExcludeDirs  []string
	OnFileChange func(string)
}

// DefaultWatchOptions returns sensible defaults for Go development
func DefaultWatchOptions() *WatchOptions {
	return &WatchOptions{
		Extensions:   []string{constants.DefaultWatchExtensions},
		ExcludeDirs:  strings.Split(constants.DefaultWatchExcludeDirs, ","),
		OnFileChange: nil, // No default callback to avoid duplication
	}
}

// NewFileWatcher creates a new FileWatcher instance
func NewFileWatcher(options *WatchOptions) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("error creating watcher: %v", err)
	}

	return &FileWatcher{
		watcher:      watcher,
		extensions:   options.Extensions,
		excludeDirs:  options.ExcludeDirs,
		onFileChange: options.OnFileChange,
	}, nil
}

// AddDirectoriesToWatch recursively adds directories to watch, excluding specified dirs
func (fw *FileWatcher) AddDirectoriesToWatch() error {
	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirName := filepath.Base(path)
			for _, excluded := range fw.excludeDirs {
				if dirName == excluded {
					return filepath.SkipDir
				}
			}
			if err := fw.watcher.Add(path); err != nil {
				fmt.Printf(constants.MsgWatchWarning+"\n", path, err)
			}
		}
		return nil
	})
}

// ShouldRestart determines if a file change should trigger a restart
func (fw *FileWatcher) ShouldRestart(event fsnotify.Event) bool {
	// Only restart on Write and Create, not Remove and Rename
	if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
		return false
	}

	ext := filepath.Ext(event.Name)
	for _, watchExt := range fw.extensions {
		if ext == watchExt {
			return true
		}
	}
	return false
}

// Watch starts watching for file changes and returns event channels
func (fw *FileWatcher) Watch() (chan fsnotify.Event, chan error) {
	events := make(chan fsnotify.Event, 100) // Buffered channel to prevent blocking
	errors := make(chan error, 10)           // Buffered channel to prevent blocking

	go func() {
		defer func() {
			if r := recover(); r != nil {
				select {
				case errors <- fmt.Errorf("panic in watcher: %v", r):
				default:
				}
			}
			// Close channels when goroutine exits
			close(events)
			close(errors)
		}()

		for {
			select {
			case event, ok := <-fw.watcher.Events:
				if !ok {
					// Channel closed, exit gracefully
					return
				}
				if fw.ShouldRestart(event) {
					// Only call callback if it's not nil
					if fw.onFileChange != nil {
						fw.onFileChange(event.Name)
					}

					// Non-blocking send to events channel
					select {
					case events <- event:
					default:
						// Channel full, skip this event
					}
				}

			case err, ok := <-fw.watcher.Errors:
				if !ok {
					// Channel closed, exit gracefully
					return
				}

				// Non-blocking send to errors channel
				select {
				case errors <- err:
				default:
					// Channel full, skip this error
				}
			}
		}
	}()

	return events, errors
}

// Close closes the file watcher
func (fw *FileWatcher) Close() error {
	if fw.watcher != nil {
		return fw.watcher.Close()
	}
	return nil
}

// IsClosed checks if the watcher is closed
func (fw *FileWatcher) IsClosed() bool {
	return fw.watcher == nil
}

// GetCurrentDir returns the current working directory
func GetCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}
