package new

import (
	"fmt"
	"regexp"
	"strings"
)

// parseTemplateName parses template name and version from template argument
// Examples: "xarch", "xarch@v1.2.3", "xarch@latest", "xarch@main"
func parseTemplateName(templateArg string) (name, version string, err error) {
	if !strings.Contains(templateArg, "@") {
		return templateArg, "latest", nil
	}

	parts := strings.Split(templateArg, "@")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid template format")
	}

	name = parts[0]
	version = parts[1]

	// Handle latest version
	if version == "latest" {
		version = ""
	}

	return name, version, nil
}

// generateModuleName generates a valid module name from project name
func generateModuleName(projectName string) string {
	// Convert project name to valid module name
	// Remove invalid characters and convert to lowercase
	moduleName := strings.ToLower(projectName)

	// Replace invalid characters with hyphens
	moduleName = regexp.MustCompile(`[^a-z0-9._-]`).ReplaceAllString(moduleName, "-")

	// Remove consecutive hyphens
	moduleName = regexp.MustCompile(`-+`).ReplaceAllString(moduleName, "-")

	// Remove leading/trailing hyphens and dots
	moduleName = strings.Trim(moduleName, "-.")

	// Remove any remaining invalid characters
	moduleName = regexp.MustCompile(`[^a-z0-9._-]`).ReplaceAllString(moduleName, "")

	// Ensure it starts with a letter or digit
	if moduleName == "" || !regexp.MustCompile(`^[a-z0-9]`).MatchString(moduleName) {
		moduleName = "project-" + moduleName
	}

	// Ensure it doesn't start or end with dot or hyphen
	moduleName = strings.Trim(moduleName, ".-")

	// Ensure minimum length
	if len(moduleName) < 3 {
		moduleName = moduleName + "-app"
	}

	// Final cleanup - remove any remaining invalid characters
	moduleName = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*[a-z0-9]$|^[a-z0-9]+$`).FindString(moduleName)
	if moduleName == "" {
		moduleName = "my-project"
	}

	return moduleName
}

// validateModuleName validates the module name according to Go standards
func validateModuleName(moduleName string) error {
	if moduleName == "" {
		return fmt.Errorf("module name cannot be empty")
	}

	// Check minimum length
	if len(moduleName) < 3 {
		return fmt.Errorf("module name must be at least 3 characters long")
	}

	// Check maximum length (Go module path limit)
	if len(moduleName) > 255 {
		return fmt.Errorf("module name cannot exceed 255 characters")
	}

	// Go module name validation regex
	// Must contain only letters, digits, dots, hyphens, underscores, and slashes
	moduleRegex := regexp.MustCompile(`^[a-zA-Z0-9._/-]+$`)
	if !moduleRegex.MatchString(moduleName) {
		return fmt.Errorf("invalid module name format. Module name must:\n" +
			"  - Contain only letters, digits, dots, hyphens, underscores, and slashes\n" +
			"  - Examples: github.com/user/repo, example.com/my-module, mycompany.com/api/v1")
	}

	// Check for consecutive dots
	if strings.Contains(moduleName, "..") {
		return fmt.Errorf("module name cannot contain consecutive dots")
	}

	// Check that module name doesn't start or end with dot, hyphen, or slash
	if strings.HasPrefix(moduleName, ".") || strings.HasSuffix(moduleName, ".") {
		return fmt.Errorf("module name cannot start or end with a dot")
	}

	if strings.HasPrefix(moduleName, "-") || strings.HasSuffix(moduleName, "-") {
		return fmt.Errorf("module name cannot start or end with a hyphen")
	}

	if strings.HasPrefix(moduleName, "/") || strings.HasSuffix(moduleName, "/") {
		return fmt.Errorf("module name cannot start or end with a slash")
	}

	// Check for reserved names
	reservedNames := []string{
		"test", "example", "internal", "vendor", "cmd", "pkg", "api", "web", "app",
		"main", "init", "go", "golang", "golang.org", "github.com", "gitlab.com",
		"bitbucket.org", "gopkg.in", "go.uber.org", "go.mongodb.org",
	}

	// Check if module name is a reserved name (case insensitive)
	lowerModuleName := strings.ToLower(moduleName)
	for _, reserved := range reservedNames {
		if lowerModuleName == reserved {
			return fmt.Errorf("module name '%s' is reserved and cannot be used", moduleName)
		}
	}

	// Check for common patterns that should be avoided
	if strings.HasPrefix(moduleName, "go-") {
		return fmt.Errorf("module name should not start with 'go-' prefix")
	}

	if strings.HasSuffix(moduleName, ".go") {
		return fmt.Errorf("module name should not end with '.go' extension")
	}

	// Check for valid domain-like structure (if it contains dots)
	if strings.Contains(moduleName, ".") {
		// Split by dots first, then check each part
		dotParts := strings.Split(moduleName, ".")
		if len(dotParts) < 2 {
			return fmt.Errorf("module name with dots should follow domain-like structure (e.g., example.com/module)")
		}

		// Check each part separated by dots
		for i, part := range dotParts {
			if part == "" {
				return fmt.Errorf("module name cannot contain empty parts between dots")
			}

			// Check that each part contains only valid characters (including slash)
			partRegex := regexp.MustCompile(`^[a-zA-Z0-9_/-]+$`)
			if !partRegex.MatchString(part) {
				return fmt.Errorf("module name parts can only contain letters, digits, hyphens, underscores, and slashes")
			}

			// First and last parts should not start/end with hyphen or slash
			if (i == 0 || i == len(dotParts)-1) && (strings.HasPrefix(part, "-") || strings.HasSuffix(part, "-") || strings.HasPrefix(part, "/") || strings.HasSuffix(part, "/")) {
				return fmt.Errorf("module name parts cannot start or end with hyphen or slash")
			}
		}
	}

	return nil
}
