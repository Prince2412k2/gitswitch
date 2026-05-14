package cmd

import (
	"fmt"
	"path/filepath"

	"gitswitch/internal/config"
	"gitswitch/internal/git"
	"gitswitch/internal/ui"
)

func runClone(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: gitswitch clone <url> [dir]")
	}

	repoURL := args[0]
	targetDir := ""
	if len(args) > 1 {
		targetDir = args[1]
	} else {
		targetDir = git.DirFromURL(repoURL)
	}

	profiles, err := config.LoadAll()
	if err != nil {
		return fmt.Errorf("loading profiles: %w", err)
	}
	if len(profiles) == 0 {
		fmt.Print("\n  no profiles yet — run `gitswitch new` to create one\n\n")
		return nil
	}

	result := ui.RunPicker(profiles, filepath.Base(targetDir))
	if result.Canceled {
		fmt.Print("\n  canceled\n\n")
		return nil
	}

	p := result.Profile

	fmt.Printf("\n  cloning %s …\n\n", repoURL)

	if err := git.CloneWithProfile(p, repoURL, targetDir); err != nil {
		return fmt.Errorf("clone failed: %w", err)
	}

	if err := git.ApplyProfileTo(p, targetDir); err != nil {
		return fmt.Errorf("applying profile: %w", err)
	}

	ui.PrintApplied(p, filepath.Base(targetDir))
	return nil
}
