package cmd

import (
	"fmt"
	"os"
	"strings"

	"gitswitch/internal/config"
	"gitswitch/internal/git"
	"gitswitch/internal/ui"
)

func runNew() error {
	fmt.Println()

	result := ui.RunForm()
	if result.Aborted {
		fmt.Println("\n  canceled\n")
		return nil
	}

	p := result.Profile

	if err := config.Save(p); err != nil {
		return fmt.Errorf("saving profile: %w", err)
	}

	path, _ := config.ProfilePath(p)

	fmt.Printf("\n  ✓  profile saved\n")
	fmt.Printf("     %s\n\n", path)

	if git.IsRepo() {
		fmt.Print("  apply to current repo? [y/N]  ")
		var ans string
		fmt.Fscan(os.Stdin, &ans)
		if strings.ToLower(strings.TrimSpace(ans)) == "y" {
			if err := git.ApplyProfile(p); err != nil {
				return fmt.Errorf("applying profile: %w", err)
			}
			root, _ := git.RepoRoot()
			ui.PrintApplied(p, root)
		}
	}

	return nil
}
