package root

import (
	"fmt"
	"strings"

	"github.com/risoftinc/elsa/constants"
	"github.com/spf13/cobra"
)

// DisplayHelper provides helper functions for displaying root command information
type DisplayHelper struct{}

// NewDisplayHelper creates a new display helper instance
func NewDisplayHelper() *DisplayHelper {
	return &DisplayHelper{}
}

// ShowRootHelp displays the custom root help template
func (h *DisplayHelper) ShowRootHelp(cmd *cobra.Command, version string) {
	fmt.Println(h.getBanner(version) + constants.UsageText)
	h.printAvailableCommands(cmd)
	fmt.Printf("\nFlags:\n%s\n", constants.HelpText)
}

// ShowSuggestions displays command suggestions when an unknown command is entered
func (h *DisplayHelper) ShowSuggestions(suggestions []string) {
	if len(suggestions) > 0 {
		fmt.Printf("%s\n", constants.SuggestionText)
		for _, suggestion := range suggestions {
			fmt.Printf("  elsa %s\n", suggestion)
		}
		fmt.Println()
	}
}

// GetVersionTemplate returns the formatted version template
func (h *DisplayHelper) GetVersionTemplate(version, goVersion, platform string) string {
	return fmt.Sprintf(constants.VersionTemplate, version, goVersion, platform)
}

// getBanner returns the formatted banner with version
func (h *DisplayHelper) getBanner(version string) string {
	return fmt.Sprintf(constants.BannerTemplate, version)
}

// printAvailableCommands prints all available commands in a formatted way
func (h *DisplayHelper) printAvailableCommands(cmd *cobra.Command) {
	maxLen := h.calculateMaxCommandLength(cmd)

	for _, c := range cmd.Commands() {
		if (!c.IsAvailableCommand() || c.Hidden) && c.Name() != "help" {
			continue
		}
		fmt.Printf("  %-*s %s\n", maxLen+1, c.Name(), c.Short)
	}
}

// calculateMaxCommandLength calculates the maximum length of command names for formatting
func (h *DisplayHelper) calculateMaxCommandLength(cmd *cobra.Command) int {
	maxLen := 0
	for _, c := range cmd.Commands() {
		if (!c.IsAvailableCommand() || c.Hidden) && c.Name() != "help" {
			continue
		}
		if l := len(c.Name()); l > maxLen {
			maxLen = l
		}
	}
	return maxLen
}

// IsFlag checks if the given argument is a flag
func (h *DisplayHelper) IsFlag(arg string) bool {
	return strings.HasPrefix(arg, "-")
}
