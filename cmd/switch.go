package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"gitswitch/internal/config"
	"gitswitch/internal/git"
	"gitswitch/internal/ui"
)

func runSwitch() error {
	if !git.IsRepo() {
		fmt.Fprintln(os.Stderr, "  ✗  not inside a git repository")
		return fmt.Errorf("not a git repo")
	}

	profiles, err := config.LoadAll()
	if err != nil {
		return fmt.Errorf("loading profiles: %w", err)
	}
	if len(profiles) == 0 {
		fmt.Println("\n  no profiles yet — run `gitswitch new` to create one\n")
		return nil
	}

	root, _ := git.RepoRoot()
	repoName := filepath.Base(root)

	curName, curEmail := git.CurrentLocal()
	if curName != "" || curEmail != "" {
		fmt.Printf("\n  current  %s  <%s>\n\n", curName, curEmail)
	}

	result := ui.RunPicker(profiles, repoName)
	if result.Canceled {
		fmt.Println("\n  canceled\n")
		return nil
	}

	if err := git.ApplyProfile(result.Profile); err != nil {
		return fmt.Errorf("applying profile: %w", err)
	}

	ui.PrintApplied(result.Profile, repoName)
	return nil
}
