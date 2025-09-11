package new

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"go.risoftinc.com/elsa/constants"
)

// cloneTemplate clones the template repository
func (tm *TemplateManager) cloneTemplate(templateURL, cachedPath, version string) error {
	// Create cache directory
	if err := os.MkdirAll(filepath.Dir(cachedPath), 0755); err != nil {
		return err
	}

	// Remove existing cache if it exists
	if _, err := os.Stat(cachedPath); err == nil {
		if err := os.RemoveAll(cachedPath); err != nil {
			return err
		}
	}

	// Clone repository
	fmt.Printf(constants.NewInfoCloningTemplate+"\n", templateURL)
	cmd := exec.Command("git", "clone", templateURL, cachedPath)
	if err := cmd.Run(); err != nil {
		return err
	}

	// If version is specified and not latest, checkout to that version
	if version != "" && version != "latest" {
		// Try to checkout as tag first, then as branch
		checkoutCmd := exec.Command("git", "checkout", version)
		checkoutCmd.Dir = cachedPath
		if err := checkoutCmd.Run(); err != nil {
			// If tag checkout fails, try branch
			checkoutCmd = exec.Command("git", "checkout", "-b", version, "origin/"+version)
			checkoutCmd.Dir = cachedPath
			if err := checkoutCmd.Run(); err != nil {
				return fmt.Errorf(constants.NewErrorVersionNotFound, version, filepath.Base(templateURL))
			}
		}
	}

	return nil
}

// cleanGitHistory removes git history from the project
func (tm *TemplateManager) cleanGitHistory(projectPath string) error {
	gitPath := filepath.Join(projectPath, ".git")
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		return nil // No .git directory to remove
	}

	return os.RemoveAll(gitPath)
}
