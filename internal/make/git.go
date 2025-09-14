package make

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// cloneRepository clones a repository to the specified directory
func (tm *TemplateManager) cloneRepository(url, targetDir, commit string) error {
	// Create target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %v", err)
	}

	// Clone repository
	cmd := exec.Command("git", "clone", url, targetDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}

	// Checkout specific commit if provided
	if commit != "" && commit != "latest" {
		checkoutCmd := exec.Command("git", "checkout", commit)
		checkoutCmd.Dir = targetDir
		if err := checkoutCmd.Run(); err != nil {
			// Try to fetch and checkout if direct checkout fails
			fetchCmd := exec.Command("git", "fetch", "origin", commit)
			fetchCmd.Dir = targetDir
			if fetchErr := fetchCmd.Run(); fetchErr == nil {
				checkoutCmd = exec.Command("git", "checkout", commit)
				checkoutCmd.Dir = targetDir
				if err := checkoutCmd.Run(); err != nil {
					// Try direct checkout without fetch
					checkoutCmd = exec.Command("git", "checkout", commit)
					checkoutCmd.Dir = targetDir
					if err := checkoutCmd.Run(); err != nil {
						return fmt.Errorf("failed to checkout commit %s: %v", commit, err)
					}
				}
			} else {
				return fmt.Errorf("failed to checkout commit %s: %v", commit, err)
			}
		}
	}

	return nil
}

// getLatestCommitFromRemote gets the latest commit hash from remote repository
func (tm *TemplateManager) getLatestCommitFromRemote(gitURL string) string {
	cmd := exec.Command("git", "ls-remote", gitURL, "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				commit := parts[0]
				// Return short commit hash (7 characters)
				if len(commit) > 7 {
					return commit[:7]
				}
				return commit
			}
		}
	}

	return ""
}

// getActualCommitFromClonedRepo gets the actual commit hash from a cloned repository
func (tm *TemplateManager) getActualCommitFromClonedRepo(repoPath string) string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	commit := strings.TrimSpace(string(output))
	// Return short commit hash (7 characters)
	if len(commit) > 7 {
		return commit[:7]
	}
	return commit
}
