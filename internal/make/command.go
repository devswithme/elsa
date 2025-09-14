package make

import (
	"fmt"
	"strings"
)

// MakeCommand handles the make command execution
type MakeCommand struct {
	templateManager *TemplateManager
}

// NewMakeCommand creates a new make command
func NewMakeCommand() *MakeCommand {
	return &MakeCommand{
		templateManager: NewTemplateManager(),
	}
}

// Execute executes the make command
func (mc *MakeCommand) Execute(args []string) error {
	if len(args) < 2 {
		return mc.showHelp()
	}

	templateType := args[0]
	name := args[1]

	// Validate name format
	if err := mc.validateName(name); err != nil {
		return err
	}

	// Generate file
	return mc.templateManager.GenerateFile(templateType, name)
}

// validateName validates the name format
func (mc *MakeCommand) validateName(name string) error {
	// Check for invalid characters
	if strings.Contains(name, "..") {
		return fmt.Errorf("invalid name: '..' not allowed")
	}

	if strings.HasPrefix(name, "/") || strings.HasSuffix(name, "/") {
		return fmt.Errorf("invalid name: cannot start or end with '/'")
	}

	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("invalid name: name cannot be empty")
	}

	return nil
}

// showHelp shows the help information
func (mc *MakeCommand) showHelp() error {
	// Try to load config to show available types
	config, err := mc.templateManager.LoadProjectConfig(".")
	if err != nil {
		fmt.Println("Usage: elsa make <template-type> <name>")
		fmt.Println("Example: elsa make repository user_repository")
		fmt.Println("Example: elsa make service user_service")
		return nil
	}

	fmt.Println("Usage: elsa make <template-type> <name>")
	fmt.Println()
	fmt.Println("Available template types:")
	for templateType := range config.Make {
		fmt.Printf("  - %s\n", templateType)
	}
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  elsa make repository user_repository")
	fmt.Println("  elsa make service user_service")
	fmt.Println("  elsa make repository health/health_repository")

	return nil
}

// ListAvailableTypes lists all available template types
func (mc *MakeCommand) ListAvailableTypes() error {
	config, err := mc.templateManager.LoadProjectConfig(".")
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	fmt.Println("Available template types:")
	for templateType := range config.Make {
		fmt.Printf("  - %s\n", templateType)
	}

	return nil
}
