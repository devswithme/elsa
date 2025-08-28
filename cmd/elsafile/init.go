package elsafile

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Elsafile in the current directory",
	Long: `Initialize creates a new Elsafile in the current directory.
The Elsafile is similar to a Makefile and allows you to define custom commands
for your project. Commands defined in Elsafile can be run using 'elsa run command'.

Example Elsafile content:
  build:
    go build -o bin/app .
  
  test:
    go test ./...
  
  clean:
    rm -rf bin/
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := createElsafile(); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating Elsafile: %v\n", err)
			os.Exit(1)
		}
	},
}

func createElsafile() error {
	// Check if Elsafile already exists
	if _, err := os.Stat("Elsafile"); err == nil {
		return fmt.Errorf("Elsafile already exists in current directory")
	}

	// Create default Elsafile content
	content := `# Elsa - Engineer's Little Smart Assistant
# This file defines custom commands for your project
# Commands can be run using: elsa run command_name

# Build the project
build:
	go build -o bin/app .

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run the application
run:
	go run .

# Format code
fmt:
	go fmt ./...
	go vet ./...

# Generate documentation
docs:
	godoc -http=:6060
`

	// Write Elsafile
	if err := os.WriteFile("Elsafile", []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write Elsafile: %v", err)
	}

	fmt.Println("✅ Created Elsafile successfully!")
	fmt.Println("📝 You can now define custom commands in the Elsafile")
	fmt.Println("🚀 Run commands using: elsa run command_name")
	fmt.Println("")
	fmt.Println("Example:")
	fmt.Println("  elsa run build    # Run the build command")
	fmt.Println("  elsa run test     # Run the test command")

	return nil
}
