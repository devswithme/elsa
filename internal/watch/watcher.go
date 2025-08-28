package watch

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
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
		Extensions: []string{".go"},
		ExcludeDirs: []string{
			".git", "vendor", "tmp", "temp", "build", "dist",
			"bin", "pkg", ".vscode", ".idea", "coverage", "testdata",
		},
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
				fmt.Printf("⚠️  Warning: Could not watch directory %s: %v\n", path, err)
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
	events := make(chan fsnotify.Event)
	errors := make(chan error)

	go func() {
		for {
			select {
			case event, ok := <-fw.watcher.Events:
				if !ok {
					close(events)
					return
				}
				if fw.ShouldRestart(event) {
					// Only call callback if it's not nil
					if fw.onFileChange != nil {
						fw.onFileChange(event.Name)
					}
					events <- event
				}

			case err, ok := <-fw.watcher.Errors:
				if !ok {
					close(errors)
					return
				}
				errors <- err
			}
		}
	}()

	return events, errors
}

// Close closes the file watcher
func (fw *FileWatcher) Close() error {
	return fw.watcher.Close()
}

// GetCurrentDir returns the current working directory
func GetCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}
