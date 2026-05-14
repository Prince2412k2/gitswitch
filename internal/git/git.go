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
	return ApplyProfileTo(p, root)
}

// ApplyProfileTo writes a profile's identity into the git config of an
// arbitrary repo directory (useful after cloning).
func ApplyProfileTo(p config.Profile, dir string) error {
	cmds := [][]string{
		{"git", "-C", dir, "config", "--local", "user.name", p.Name},
		{"git", "-C", dir, "config", "--local", "user.email", p.Email},
	}

	if p.SSHKey != "" {
		expanded := expandHome(p.SSHKey)
		sshCmd := fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes", expanded)
		cmds = append(cmds, []string{"git", "-C", dir, "config", "--local", "core.sshCommand", sshCmd})
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

// CloneWithProfile clones a repo using the SSH key from the given profile.
// Git's output streams directly to the terminal so the user sees progress.
func CloneWithProfile(p config.Profile, url, dir string) error {
	args := []string{"clone", url}
	if dir != "" {
		args = append(args, dir)
	}

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if p.SSHKey != "" {
		expanded := expandHome(p.SSHKey)
		sshVal := fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes", expanded)
		cmd.Env = append(os.Environ(), "GIT_SSH_COMMAND="+sshVal)
	}

	return cmd.Run()
}

// DirFromURL derives a local directory name from a clone URL, the same way
// git itself does: take the last path segment and strip a trailing ".git".
func DirFromURL(rawURL string) string {
	// works for https://github.com/org/repo.git and git@github.com:org/repo.git
	s := strings.TrimRight(rawURL, "/")
	// for SCP-style URLs, replace the colon-prefixed path separator
	if idx := strings.LastIndex(s, ":"); idx > strings.LastIndex(s, "/") {
		s = s[idx+1:]
	}
	base := filepath.Base(s)
	return strings.TrimSuffix(base, ".git")
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
