package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gitswitch/internal/config"
)

// RepoRoot returns the root of the current git repo, or "" if not in one.
func RepoRoot() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", nil // not a git repo, not an error
	}
	return strings.TrimSpace(string(out)), nil
}

// IsRepo returns true if the cwd is inside a git repo.
func IsRepo() bool {
	root, err := RepoRoot()
	return err == nil && root != ""
}

// CurrentLocal reads the local git config user.name and user.email.
func CurrentLocal() (name, email string) {
	n, _ := exec.Command("git", "config", "--local", "user.name").Output()
	e, _ := exec.Command("git", "config", "--local", "user.email").Output()
	return strings.TrimSpace(string(n)), strings.TrimSpace(string(e))
}

// CurrentSSHCommand reads the local core.sshCommand.
func CurrentSSHCommand() string {
	out, _ := exec.Command("git", "config", "--local", "core.sshCommand").Output()
	return strings.TrimSpace(string(out))
}

// ApplyProfile writes a profile's identity into the repo's local git config.
func ApplyProfile(p config.Profile) error {
	root, err := RepoRoot()
	if err != nil || root == "" {
		return fmt.Errorf("not inside a git repository")
	}

	cmds := [][]string{
		{"git", "config", "--local", "user.name", p.Name},
		{"git", "config", "--local", "user.email", p.Email},
	}

	if p.SSHKey != "" {
		expanded := expandHome(p.SSHKey)
		sshCmd := fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes", expanded)
		cmds = append(cmds, []string{"git", "config", "--local", "core.sshCommand", sshCmd})
	}

	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run %v: %w", args, err)
		}
	}

	return nil
}

// ConfigPath returns the path to .git/config in the current repo.
func ConfigPath() (string, error) {
	root, err := RepoRoot()
	if err != nil || root == "" {
		return "", fmt.Errorf("not inside a git repository")
	}
	return filepath.Join(root, ".git", "config"), nil
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
