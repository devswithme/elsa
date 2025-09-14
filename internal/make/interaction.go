package make

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// promptOutputDirectory prompts user for output directory
func (tm *TemplateManager) promptOutputDirectory(templateType string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("📁 Please enter the output directory for %s: ", templateType)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %v", err)
	}

	outputDir := strings.TrimSpace(input)
	if outputDir == "" {
		return "", fmt.Errorf("output directory cannot be empty")
	}

	fmt.Printf("📁 Using output directory: %s\n", outputDir)
	return outputDir, nil
}
