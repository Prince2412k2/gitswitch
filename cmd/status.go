package cmd

import (
	"fmt"
	"os"

	"gitswitch/internal/config"
	"gitswitch/internal/git"
	"gitswitch/internal/ui"
)

func runStatus() error {
	if !git.IsRepo() {
		fmt.Fprintln(os.Stderr, "  ✗  not inside a git repository")
		return fmt.Errorf("not a git repo")
	}

	name, email := git.CurrentLocal()
	sshCmd := git.CurrentSSHCommand()

	if name == "" && email == "" {
		fmt.Print("\n  no local git identity set — run: gitswitch\n\n")
		return nil
	}

	profiles, _ := config.LoadAll()
	var matched *config.Profile
	for i, p := range profiles {
		if p.Name == name && p.Email == email {
			matched = &profiles[i]
			break
		}
	}

	root, _ := git.RepoRoot()
	ui.PrintStatus(name, email, sshCmd, matched, root)
	return nil
}
